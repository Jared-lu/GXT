package redis_lock

import (
	"context"
	"errors"
	"fmt"
	redismock "github.com/Jared-lu/GXT/redis-lock/mock/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestClient_Lock(t *testing.T) {
	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) redis.Cmdable
		key        string
		expiration time.Duration
		timeout    time.Duration
		retry      RetryStrategy
		wantErr    error
		wantLock   *Lock
	}{
		{
			name: "eval error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismock.NewMockCmdable(ctrl)

				res := redis.NewCmd(context.Background())
				res.SetErr(context.DeadlineExceeded)
				cmd.EXPECT().Eval(gomock.Any(),
					luaLock, []string{"lock-key1"}, gomock.Any()).
					Return(res)

				return cmd

			},
			key:        "lock-key1",
			expiration: time.Minute,
			timeout:    time.Second * 5,
			retry:      nil,
			wantErr:    context.DeadlineExceeded,
			wantLock:   nil,
		},
		{
			name: "lock success",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismock.NewMockCmdable(ctrl)

				res := redis.NewCmd(context.Background())
				res.SetVal(int64(1))
				cmd.EXPECT().Eval(gomock.Any(),
					luaLock, []string{"lock-key2"}, gomock.Any()).
					Return(res)

				return cmd

			},
			key:        "lock-key2",
			expiration: time.Minute,
			timeout:    time.Second * 5,
			retry:      nil,
			wantErr:    nil,
			wantLock: &Lock{
				key:        "lock-key2",
				expiration: time.Minute,
			},
		},
		{
			name: "other hold lock, retry and failed",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismock.NewMockCmdable(ctrl)

				res := redis.NewCmd(context.Background())
				res.SetVal(int64(2))
				cmd.EXPECT().Eval(gomock.Any(),
					luaLock, []string{"lock-key3"}, gomock.Any()).Times(4).Return(res)

				return cmd

			},
			key:        "lock-key3",
			expiration: time.Minute,
			timeout:    time.Second * 5,
			retry: &FixedIntervalRetryStrategy{
				MaxCnt:   3,
				Interval: time.Second,
			},
			wantErr:  fmt.Errorf("超出重试限制, %w", ErrFailedToPreemptLock),
			wantLock: nil,
		},
		{
			name: "retry and success",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismock.NewMockCmdable(ctrl)

				res := redis.NewCmd(context.Background())
				res.SetVal(int64(2))
				cmd.EXPECT().Eval(gomock.Any(),
					luaLock, []string{"lock-key4"}, gomock.Any()).Times(2).Return(res)
				// 重试后成功
				res2 := redis.NewCmd(context.Background())
				res2.SetVal(int64(1))
				cmd.EXPECT().Eval(gomock.Any(),
					luaLock, []string{"lock-key4"}, gomock.Any()).Return(res2)

				return cmd

			},
			key:        "lock-key4",
			expiration: time.Minute,
			timeout:    time.Second * 5,
			retry: &FixedIntervalRetryStrategy{
				MaxCnt:   3,
				Interval: time.Second,
			},
			wantErr: nil,
			wantLock: &Lock{
				key:        "lock-key4",
				expiration: time.Minute,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			client := NewClient(tc.mock(ctrl))

			lock, err := client.Lock(context.Background(), tc.key, tc.expiration, tc.timeout, tc.retry)

			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantLock.key, lock.key)
			assert.Equal(t, tc.wantLock.expiration, lock.expiration)
			assert.NotEmpty(t, lock.value)
			assert.NotNil(t, lock.client)
		})
	}
}

func TestClient_TryLock(t *testing.T) {
	testCases := []struct {
		name     string
		mockCmd  func(ctrl *gomock.Controller) redis.Cmdable
		key      string
		wantErr  error
		wantLock *Lock
	}{
		{
			name: "setnx error",
			mockCmd: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismock.NewMockCmdable(ctrl)
				// 这个用来模拟 setnx 命令的返回值
				res := redis.NewBoolResult(false, context.DeadlineExceeded)
				cmd.EXPECT().SetNX(context.Background(),
					"key1", gomock.Any(), time.Second*10).
					Return(res)
				return cmd
			},
			key:     "key1",
			wantErr: context.DeadlineExceeded,
		},
		{
			name: "preempt lock failed",
			mockCmd: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismock.NewMockCmdable(ctrl)
				// 这个用来模拟 setnx 命令的返回值
				res := redis.NewBoolResult(false, nil)
				cmd.EXPECT().SetNX(context.Background(),
					"key1", gomock.Any(), time.Second*10).
					Return(res)
				return cmd
			},
			key:     "key1",
			wantErr: ErrFailedToPreemptLock,
		},
		{
			name: "success",
			mockCmd: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismock.NewMockCmdable(ctrl)
				// 这个用来模拟 setnx 命令的返回值
				res := redis.NewBoolResult(true, nil)
				cmd.EXPECT().SetNX(context.Background(),
					"key1", gomock.Any(), time.Second*10).
					Return(res)
				return cmd
			},
			key:     "key1",
			wantErr: nil,
			wantLock: &Lock{
				key:        "key1",
				expiration: time.Second * 10,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			client := NewClient(tc.mockCmd(ctrl))

			lock, err := client.TryLock(context.Background(), tc.key, time.Second*10)

			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantLock.key, lock.key)
			assert.Equal(t, tc.wantLock.expiration, lock.expiration)
			assert.NotEmpty(t, lock.value)
		})
	}
}

func TestClient_Unlock(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		key     string
		value   string
		wantErr error
	}{
		{
			name: "eval error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismock.NewMockCmdable(ctrl)

				res := redis.NewCmd(context.Background())
				res.SetErr(context.DeadlineExceeded)
				cmd.EXPECT().Eval(context.Background(),
					luaUnlock, []string{"key1"}, []any{"value1"}).
					Return(res)

				return cmd

			},
			key:     "key1",
			value:   "value1",
			wantErr: context.DeadlineExceeded,
		},
		{
			name: "unlock failed",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismock.NewMockCmdable(ctrl)

				res := redis.NewCmd(context.Background())
				res.SetVal(int64(0))
				cmd.EXPECT().Eval(context.Background(),
					luaUnlock, []string{"key1"}, []any{"value1"}).
					Return(res)

				return cmd

			},
			key:     "key1",
			value:   "value1",
			wantErr: ErrLockNotHeld,
		},
		{
			name: "unlock success",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismock.NewMockCmdable(ctrl)

				res := redis.NewCmd(context.Background())
				res.SetVal(int64(1))
				cmd.EXPECT().Eval(context.Background(),
					luaUnlock, []string{"key1"}, []any{"value1"}).
					Return(res)

				return cmd

			},
			key:     "key1",
			value:   "value1",
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			lock := &Lock{
				key:    tc.key,
				value:  tc.value,
				client: tc.mock(ctrl),
			}
			err := lock.Unlock(context.Background())
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestClient_Refresh(t *testing.T) {
	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) redis.Cmdable
		key        string
		value      string
		expiration time.Duration
		wantErr    error
	}{
		{
			name: "eval error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismock.NewMockCmdable(ctrl)

				res := redis.NewCmd(context.Background())
				res.SetErr(context.DeadlineExceeded)
				cmd.EXPECT().Eval(context.Background(),
					luaRefresh, []string{"key1"}, []any{"value1", float64(60)}).
					Return(res)

				return cmd

			},
			key:        "key1",
			value:      "value1",
			expiration: time.Minute,
			wantErr:    context.DeadlineExceeded,
		},
		{
			name: "refresh failed",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismock.NewMockCmdable(ctrl)

				res := redis.NewCmd(context.Background())
				res.SetVal(int64(0))
				cmd.EXPECT().Eval(context.Background(),
					luaRefresh, []string{"key1"}, []any{"value1", float64(60)}).
					Return(res)

				return cmd

			},
			key:        "key1",
			value:      "value1",
			expiration: time.Minute,
			wantErr:    ErrLockNotHeld,
		},
		{
			name: "refresh success",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismock.NewMockCmdable(ctrl)

				res := redis.NewCmd(context.Background())
				res.SetVal(int64(1))
				cmd.EXPECT().Eval(context.Background(),
					luaRefresh, []string{"key1"}, []any{"value1", float64(60)}).
					Return(res)

				return cmd

			},
			key:        "key1",
			value:      "value1",
			expiration: time.Minute,
			wantErr:    nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			lock := &Lock{
				key:        tc.key,
				value:      tc.value,
				client:     tc.mock(ctrl),
				expiration: tc.expiration,
			}
			err := lock.Refresh(context.Background())
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func ExampleLock_Refresh() {
	// 假设加锁成功，拿到了lock
	var l *Lock
	// 业务执行完毕的信号
	stopChan := make(chan struct{})
	// 续约失败的信号
	errChan := make(chan error)
	done := false
	go func() {
		timeoutChan := make(chan struct{}, 1)
		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()
		for !done {
			select {
			case <-stopChan:
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				_ = l.Unlock(ctx)
				cancel()
				done = true
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
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
					errChan <- err
					_ = l.Unlock(ctx)
					done = true
				}
			case <-timeoutChan:
				// 尝试一次重试
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				err := l.Refresh(ctx)
				cancel()
				if err == nil {
					// 重试成功
					continue
				}
				errChan <- err
				_ = l.Unlock(ctx)
				done = true
			}
		}
	}()
	// 你的业务
	// 业务方在执行时，要在中间步骤检测errChan有没有信号
	// 如果续约失败了，业务方需要中断正在处理的业务

	// 如果业务在循环中执行
FOR:
	for i := 0; i < 100; i++ {
		select {
		case <-errChan:
			// 中断业务
			break FOR
		default:
			// 正常的业务处理
		}
	}

	// 如果你的业务是在多个步骤中执行
	select {
	case <-errChan:
	// 中断业务
	default:
		// 步骤一
	}
	select {
	case <-errChan:
	// 中断业务
	default:
		// 步骤二
	}

	select {
	case <-errChan:
	// 中断业务
	default:
		// 步骤n
	}

	// 使用context进行传递中断信号
	ctx, cancel := context.WithCancel(context.Background())
	// 防止正常结束了业务，但是没人发送ctx取消信号而导致goroutine泄露
	defer cancel()
	go func() {
		for {
			select {
			case <-errChan:
				cancel()
				continue
			case <-ctx.Done():
				// 这一个是为了防止goroutine泄露
				return
			}
		}
	}()
	// 使用ctx往下传递，下游在执行前需要先判断一下ctx有没有被cancel()
	// Next(ctx)

	// 业务退出就要退出续约循环
	// 不管是被中断了还是正常结束，都需要发送这个信号
	stopChan <- struct{}{}
	close(stopChan)
	close(errChan)

	fmt.Println("Hello")
	// Output:
	// Hello
}

func ExampleLock_AutoRefresh() {
	// 抢到了锁
	var l *Lock
	go func() {
		_ = l.AutoRefresh(time.Millisecond, time.Second*3)
		// 处理错误，需要中断业务

	}()
	// 执行业务

	fmt.Println("Hello")
	// Output:
	// Hello
}
