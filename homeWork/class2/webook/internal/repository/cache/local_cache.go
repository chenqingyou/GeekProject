package cache

import (
	"GeekProject/homeWork/class2/webook/internal/domain"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrUserNotNoFound = errors.New("不存在")

type UserLocalCache struct {
	cache      *sync.Map
	expiration time.Duration
}

func NewUserLocalCache() *UserLocalCache {
	return &UserLocalCache{
		cache:      &sync.Map{},
		expiration: 15 * time.Minute,
	}
}

func (uc *UserLocalCache) Get(ctx context.Context, id int64) (domain.UserDomain, error) {
	key := uc.Key(id)
	cached, ok := uc.cache.Load(key)
	if !ok {
		return domain.UserDomain{}, ErrUserNotNoFound
	}
	doU, ok := cached.(domain.UserDomain)
	if !ok {
		return domain.UserDomain{}, fmt.Errorf("invalid cached value type")
	}
	return doU, nil
}

func (uc *UserLocalCache) Set(ctx context.Context, doU domain.UserDomain) error {
	key := uc.Key(doU.Id)
	uc.cache.Store(key, doU)
	return nil
}

func (uc *UserLocalCache) Key(id int64) string {
	return fmt.Sprintf("user:info:%v", id)
}
