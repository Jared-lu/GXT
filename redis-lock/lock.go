package redis_lock

import (
	"context"
	_ "embed"
	"errors"
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

type Client struct {
	client redis.Cmdable
}

func NewClient(client redis.Cmdable) *Client {
	return &Client{client: client}
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
		client: c.client,
		key:    key,
		value:  val,
	}, nil
}

type Lock struct {
	client redis.Cmdable
	// key + value 才是锁的唯一标识
	key   string
	value string
}

func (l *Lock) Unlock(ctx context.Context) error {
	res, err := l.client.Eval(ctx, luaUnlock, []string{l.key}, l.value).Int64()
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
