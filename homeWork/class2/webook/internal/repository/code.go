package repository

import (
	"sync"
	"time"
)

var (
	// 存储验证码和过期时间
	codeMap = make(map[string]string)
	// 存储验证码尝试次数
	codeCountMap = make(map[string]int)
	// 用于保证map操作的原子性
	mu sync.Mutex
)

func setCode(key, val string) {
	mu.Lock()
	defer mu.Unlock()
	codeMap[key] = val
	codeCountMap[key] = 3
	// 设置600秒后过期
	go func() {
		time.Sleep(60 * time.Second)
		mu.Lock()
		delete(codeMap, key)
		delete(codeCountMap, key)
		mu.Unlock()
	}()
}

func checkAndSetCode(key, val string) int {
	mu.Lock()
	defer mu.Unlock()
	ttl, ok := codeMap[key]
	if !ok {
		// key不存在，直接设置
		setCode(key, val)
		return 0
	}
	// key存在但没有过期时间（不应该发生）
	// 返回-2表示系统错误
	if ttl == "" {
		return -2
	}
	// key存在且未过期，但发送太频繁（在10分钟内）
	// 返回-1表示发送太频繁
	return -1
}
