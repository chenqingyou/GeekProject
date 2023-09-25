package ratelimit_win

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed slide_window.lua
var SlideWindow string

type RedisSlidingWindowLimiter struct {
	cmd redis.Cmdable
	//窗口大小
	interval time.Duration
	// 阈值
	rate int
	//interval 内允许rate个请求
}

func NewRedisSlidingWindowLimiter(cmd redis.Cmdable, interval time.Duration, rate int) LimitInterface {
	return &RedisSlidingWindowLimiter{
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}

func (r *RedisSlidingWindowLimiter) Limited(ctx context.Context, key string) (bool, error) {
	return r.cmd.Eval(ctx, SlideWindow, []string{key},
		r.interval.Milliseconds(), r.rate, time.Now().UnixMilli()).Bool()
}
