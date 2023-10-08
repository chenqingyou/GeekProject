package ratelimit_win

import (
	"context"
	_ "embed"
	"time"
)

type LocalSlidingWindowLimiter struct {
	//窗口大小
	interval time.Duration
	// 阈值
	rate             int
	localRateLimiter *RateLimiter
	//interval 内允许rate个请求
}

func NewLocalSlidingWindowLimiter(interval time.Duration, rate int, localRateLimiter *RateLimiter) LimitInterface {
	return &LocalSlidingWindowLimiter{
		interval:         interval,
		rate:             rate,
		localRateLimiter: localRateLimiter,
	}
}

func (r *LocalSlidingWindowLimiter) Limited(ctx context.Context, key string) (bool, error) {
	return r.localRateLimiter.Access(key, time.Now().UnixNano()/1e6), nil
}
