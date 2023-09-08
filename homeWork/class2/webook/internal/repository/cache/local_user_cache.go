package cache

import (
	"GeekProject/homeWork/class2/webook/internal/domain"
	"context"
	"fmt"
	"time"
)

var ErrUserNotNoFound = fmt.Errorf("user not found")

type UserCacheItem struct {
	Value      domain.UserDomain
	Expiration time.Time
}

func (lc *UserLocalCache) Get(ctx context.Context, id int64) (domain.UserDomain, error) {
	key := lc.Key(id)
	value, ok := lc.cache.Load(key)
	if !ok {
		return domain.UserDomain{}, ErrUserNotNoFound
	}
	item := value.(*UserCacheItem)
	if time.Now().After(item.Expiration) {
		lc.cache.Delete(key)
		return domain.UserDomain{}, ErrUserNotNoFound
	}
	return item.Value, nil
}

func (lc *UserLocalCache) Set(ctx context.Context, doU domain.UserDomain) error {
	expiration := time.Now().Add(lc.expiration)
	item := &UserCacheItem{Value: doU, Expiration: expiration}
	lc.cache.Store(lc.Key(doU.Id), item)
	return nil
}

func (lc *UserLocalCache) Key(id int64) string {
	return fmt.Sprintf("user:info:%v", id)
}
