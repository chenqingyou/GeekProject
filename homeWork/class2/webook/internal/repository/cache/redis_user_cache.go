package cache

import (
	"GeekProject/homeWork/class2/webook/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var ErrUserNotNoFoundRedis = redis.Nil

func (uc *UserRedisCache) Get(ctx context.Context, id int64) (domain.UserDomain, error) {
	key := uc.Key(id)
	bytesU, err := uc.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.UserDomain{}, ErrUserNotNoFound
	}
	var doU domain.UserDomain
	err = json.Unmarshal(bytesU, &doU)
	return doU, err
}

func (uc *UserRedisCache) Set(ctx context.Context, doU domain.UserDomain) error {
	marshalDoU, err := json.Marshal(doU)
	if err != nil {
		return err
	}
	return uc.client.Set(ctx, uc.Key(doU.Id), marshalDoU, uc.expiration).Err()
}

func (uc *UserRedisCache) Key(id int64) string {
	return fmt.Sprintf("user:info:%v", id)
}
