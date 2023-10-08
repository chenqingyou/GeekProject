package ratelimit_win

import (
	"fmt"
	"testing"
	"time"
)

func TestRateLimiter_Access(t *testing.T) {
	rl := NewRateLimiter(2, 5)
	for i := 0; i < 100; i++ {
		now := time.Now().UnixNano() / 1e6
		// 检查对于键"some_key"是否达到了请求限制
		if rl.Access("some_key", now) {
			fmt.Printf("Request %d: Limited\n", i+1) // 被限流
		} else {
			fmt.Printf("Request %d: Allowed\n", i+1) // 未被限流
		}
		time.Sleep(50 * time.Millisecond) // 每10秒发起一个请求
	}
}
