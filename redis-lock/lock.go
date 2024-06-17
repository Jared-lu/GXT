package redis_lock

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	ErrFailedToPreemptLock = errors.New("failed to preempt lock")
	ErrLockNotHeld         = errors.New("lock not held")
)

//go:embed lua/unlock.lua
var luaUnlock string

//go:embed lua/refresh.lua
var luaRefresh string

//go:embed lua/lock.lua
var luaLock string

type Client struct {
	client redis.Cmdable
}

func NewClient(client redis.Cmdable) *Client {
	return &Client{client: client}
}

func (c *Client) Lock(ctx context.Context, key string,
	expiration time.Duration, timeout time.Duration, retry RetryStrategy) (*Lock, error) {
	var ticker *time.Ticker
	val := uuid.New().String()
	for {
		ctx2, cancel := context.WithTimeout(ctx, timeout)
		res, err := c.client.Eval(ctx2, luaLock, []string{key}, val, expiration.Seconds()).Int64()
		cancel()
		if err != nil || errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}

		if res == 1 {
			return &Lock{
				client:     c.client,
				key:        key,
				value:      val,
				expiration: expiration,
				unlockChan: make(chan struct{}, 1),
			}, nil
		}

		// 这里要执行重试
		// 我怎么知道还能不能重试，怎么重试，重试间隔是多久？
		// 调用方传入RetryStrategy来决定重试策略
		interval, ok := retry.Next()
		if !ok {
			return nil, fmt.Errorf("超出重试限制, %w", ErrFailedToPreemptLock)
		}
		if ticker == nil {
			ticker = time.NewTicker(interval)
		} else {
			ticker.Reset(interval)
		}
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (c *Client) TryLock(ctx context.Context,
	key string, expiration time.Duration) (*Lock, error) {
	val := uuid.New().String()
	// 设置特定键值对成功，就代表加锁成功
	ok, err := c.client.SetNX(ctx, key, val, expiration).Result()
	if err != nil {
		// 如果是超时，会进来这里
		return nil, err
	}
	if !ok {
		// 别人抢到了锁
		return nil, ErrFailedToPreemptLock
	}

	return &Lock{
		client:     c.client,
		key:        key,
		value:      val,
		expiration: expiration,
		unlockChan: make(chan struct{}, 1),
	}, nil
}

type Lock struct {
	client redis.Cmdable
	// key + value 才是锁的唯一标识
	key        string
	value      string
	expiration time.Duration
	unlockChan chan struct{}
}

// AutoRefresh 不建议使用这个API，用户最好是自己手动控制Refresh
// interval 多久续约一次
// timeout 调用续约的超时时间
func (l *Lock) AutoRefresh(interval time.Duration, timeout time.Duration) error {
	timeoutChan := make(chan struct{}, 1)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-l.unlockChan:
			// 用户主动释放锁
			return nil
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			err := l.Refresh(ctx)
			cancel()
			switch {
			// 处理error
			case errors.Is(err, context.DeadlineExceeded):
				// 自己给自己发信号，channel一定要带缓存的
				timeoutChan <- struct{}{}
				continue
			case err == nil:
			default:
				return err
			}
		case <-timeoutChan:
			// 尝试一次重试
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			err := l.Refresh(ctx)
			cancel()
			if err == nil {
				// 重试成功
				continue
			}
			return err
		}
	}
}

// Refresh 用户手动续约
func (l *Lock) Refresh(ctx context.Context) error {
	res, err := l.client.Eval(ctx, luaRefresh, []string{l.key}, l.value, l.expiration.Seconds()).Int64()
	if err != nil {
		return err
	}
	if res != 1 {
		// 不是自己的锁
		return ErrLockNotHeld
	}
	return nil
}

func (l *Lock) Unlock(ctx context.Context) error {
	res, err := l.client.Eval(ctx, luaUnlock, []string{l.key}, l.value).Int64()
	defer func() {
		close(l.unlockChan)
	}()
	if err != nil {
		return err
	}
	if res != 1 {
		// 不是自己的锁
		return ErrLockNotHeld
	}
	return nil
}

// Unlock 定义在Lock结构体上，用户使用起来会更接近面向对象的实现：
// lock,_ := c.TryLock()
// lock.Unlock()
//func (l *Lock) Unlock(ctx context.Context) error {
//	// 要先判断一下这把锁是不是我的锁
//	val, err := l.client.Get(ctx, l.key).Result()
//	if err != nil {
//		return err
//	}
//	if val != l.value {
//		return errors.New("不是自己的锁")
//	}
//
//	// 这中间会有并发问题，所以不能这么写
//
//	// 把键值对删掉
//	cnt, err := l.client.Del(ctx, l.key).Result()
//	if err != nil {
//		// 可能超时，可能返回时网络断了
//		// 不确定有没有删成功
//		return err
//	}
//	if cnt != 1 {
//		// 代表你加的锁过期了
//		return errors.New("解锁失败, 锁不存在")
//	}
//
//}
