package redis_lock

import "time"

type RetryStrategy interface {
	// Next
	// 第一个返回值：重试的间隔
	// 第二个返回值：要不要继续重试
	Next() (time.Duration, bool)
}

// FixedIntervalRetryStrategy 固定间隔重试
type FixedIntervalRetryStrategy struct {
	Interval time.Duration
	MaxCnt   int
	cnt      int
}

func (f *FixedIntervalRetryStrategy) Next() (time.Duration, bool) {
	if f.cnt >= f.MaxCnt {
		return 0, false
	}
	f.cnt++
	return f.Interval, true
}
