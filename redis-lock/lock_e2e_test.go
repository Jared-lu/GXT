//go:build e2e

package redis_lock

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_e2e_TyLock(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	testcases := []struct {
		name     string
		key      string
		before   func(t *testing.T)
		after    func(t *testing.T)
		wantErr  error
		wantLock *Lock
	}{
		{
			name: "别人持有锁",
			before: func(t *testing.T) {
				// 先设置好锁
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				res, err := rdb.Set(ctx, "key1", "value1", time.Minute).Result()
				assert.NoError(t, err)
				assert.Equal(t, "OK", res)
			},
			after: func(t *testing.T) {
				// 验证一下redis是不是真的有别人设置的key，并且要把测试的key删掉
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				res, err := rdb.GetDel(ctx, "key1").Result()
				assert.NoError(t, err)
				assert.Equal(t, "value1", res)
			},
			key:      "key1",
			wantErr:  ErrFailedToPreemptLock,
			wantLock: nil,
		},
		{
			name: "加锁成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				res, err := rdb.GetDel(ctx, "key2").Result()
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
			},
			key:     "key2",
			wantErr: nil,
			wantLock: &Lock{
				key:    "key2",
				client: rdb,
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			client := NewClient(rdb)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			lock, err := client.TryLock(ctx, tc.key, time.Minute)
			cancel()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantLock.key, lock.key)
			assert.NotEmpty(t, lock.value)
			assert.NotNil(t, lock.client)
			defer tc.after(t)
		})
	}
}

func Test_e2e_Unlock(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	testCases := []struct {
		name    string
		before  func(t *testing.T)
		after   func(t *testing.T)
		lock    *Lock
		wantErr error
	}{
		{
			name: "不是自己的锁",
			before: func(t *testing.T) {
				// 不是自己的锁，也就是value不对
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				res, err := rdb.Set(ctx, "key1", "not my key", time.Minute).Result()
				assert.NoError(t, err)
				assert.Equal(t, "OK", res)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				res, err := rdb.GetDel(ctx, "key1").Result()
				assert.NoError(t, err)
				assert.Equal(t, "not my key", res)
			},
			lock: &Lock{
				key:    "key1",
				value:  "value1",
				client: rdb,
			},
			wantErr: ErrLockNotHeld,
		},
		{
			name: "释放锁成功",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				res, err := rdb.Set(ctx, "key2", "value2", time.Minute).Result()
				assert.NoError(t, err)
				assert.Equal(t, "OK", res)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				// 释放了锁，key就不存在了
				res, err := rdb.Exists(ctx, "key2").Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), res)
			},
			lock: &Lock{
				key:    "key2",
				value:  "value2",
				client: rdb,
			},
			wantErr: nil,
		},
		{
			name: "没有持有锁",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {

			},
			lock: &Lock{
				key:    "non-exist-key",
				value:  "value",
				client: rdb,
			},
			wantErr: ErrLockNotHeld,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			err := tc.lock.Unlock(ctx)
			assert.Equal(t, tc.wantErr, err)
			tc.after(t)
		})
	}
}
