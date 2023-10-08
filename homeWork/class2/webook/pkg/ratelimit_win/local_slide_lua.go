package ratelimit_win

import (
	"sync"
)

// RateLimiter 结构体代表了一个滑动窗口限流器
type RateLimiter struct {
	mu        sync.Mutex         // 用于确保多个goroutine并发时的线程安全
	data      map[string][]int64 // 用于存储每个键对应的时间戳数组
	window    int64              // 定义滑动窗口的时长，以秒为单位
	threshold int                // 在滑动窗口时长内允许的最大请求次数
}

// NewRateLimiter 初始化并返回一个RateLimiter的实例
func NewRateLimiter(window int64, threshold int) *RateLimiter {
	return &RateLimiter{
		data:      make(map[string][]int64),
		window:    window,
		threshold: threshold,
	}
}

// Access 判断给定的键是否超过了限流器设置的请求限制
func (rl *RateLimiter) Access(key string, now int64) bool {
	rl.mu.Lock()                // 加锁以保证线程安全
	defer rl.mu.Unlock()        // 解锁
	min := now - rl.window*1000 // 计算滑动窗口的开始时间点
	// 过滤掉窗口之外的时间戳，只保留滑动窗口时长内的时间戳
	var newTimestamps []int64
	if timestamps, ok := rl.data[key]; ok {
		for _, ts := range timestamps {
			if ts > min {
				newTimestamps = append(newTimestamps, ts)
			}
		}
	}
	rl.data[key] = newTimestamps // 更新键对应的时间戳数组
	// 如果当前键在滑动窗口时长内的请求次数超过了阈值，返回true表示需要限流
	if len(rl.data[key]) >= rl.threshold {
		return true
	}
	// 如果未达到请求阈值，将当前的时间戳添加到键的时间戳数组中，并返回false表示不限流
	rl.data[key] = append(rl.data[key], now)
	return false
}
