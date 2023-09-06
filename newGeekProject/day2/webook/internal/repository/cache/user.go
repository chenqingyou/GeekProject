package cache

import (
	"GeekProject/newGeekProject/day2/webook/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrUserNotNoFound = redis.Nil

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.UserDomain, error)
	Set(ctx context.Context, doU domain.UserDomain) error
}

type UserRedisCache struct {
	//单机和集群都可以
	client     redis.Cmdable
	expiration time.Duration
}

// NewUserCache A用到B ，B一定是接口 ==>保证面向接口
// A用到B ，B一定是A的字段==>规避包变量
// A用到B ，A绝对不初始化，一定是外面注入
func NewUserCache(client redis.Cmdable) UserCache {
	return &UserRedisCache{
		client:     client,
		expiration: 15 * time.Minute,
	}
}

// Get 只有error为nil，就认为缓存里面有数据
// 如果没有数据，就返回一个特定的error
func (uc *UserRedisCache) Get(ctx context.Context, id int64) (domain.UserDomain, error) {
	key := uc.Key(id)
	//如果数据不存在就返回err
	bytesU, err := uc.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.UserDomain{}, err
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
