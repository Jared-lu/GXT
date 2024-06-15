package redis_lock

import (
	"context"
	redismock "github.com/Jared-lu/GXT/redis-lock/mock/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

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
				key: "key1",
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
