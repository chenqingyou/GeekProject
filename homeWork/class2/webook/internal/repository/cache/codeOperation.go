package cache

import (
	"sync"
	"time"
)

var (
	// 存储验证码尝试次数
	codeCount = sync.Map{}
	// 用于保证map操作的原子性
	mu        sync.Mutex
	cacheCode = sync.Map{}
)

type CacheItem struct {
	Code       string
	Cnt        int
	Expiration time.Time
}

func setCodeUnlocked(key, val string) {
	item := CacheItem{
		Code:       val,
		Cnt:        3,
		Expiration: time.Now().Add(5 * time.Minute),
	}
	cacheCode.Store(key, item)
	// 设置600秒后过期
	go func() {
		time.Sleep(60 * time.Second)
		mu.Lock()
		cacheCode.Delete(key)
		mu.Unlock()
	}()
}

func checkAndSetCode(key, val string) int {
	mu.Lock()
	defer mu.Unlock()
	ttl, ok := cacheCode.Load(key)
	if !ok {
		// key不存在，直接设置
		setCodeUnlocked(key, val)
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

func checkCode(key, expectedCode string) int {
	mu.Lock()
	defer mu.Unlock()
	value, ok := cacheCode.Load(key)
	if !ok {
		return -1
	}
	item := value.(CacheItem)
	if item.Expiration.Before(time.Now()) {
		// 已过期
		cacheCode.Delete(key)
		return -1
	}
	if item.Cnt <= 0 {
		// 无效的尝试或已使用
		return -1
	} else if expectedCode == item.Code {
		// 代码匹配
		item.Cnt = -1
		return 0
	} else {
		// 代码不匹配
		item.Cnt--
		return -2
	}
}
