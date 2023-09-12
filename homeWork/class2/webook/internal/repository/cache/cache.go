package cache

import (
	"GeekProject/homeWork/class2/webook/internal/domain"
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

var (
	ErrSetCodeFrequently    = errors.New("发送验证码太频繁")
	ErrVerifyCodeFrequently = errors.New("验证次数太多")
	ErrUnknowForCode        = errors.New("未知错误")
)

type UserLocalCache struct {
	cache      sync.Map
	expiration time.Duration
}

func NewUserLocalCache() CodeCache {
	return &UserLocalCache{
		expiration: 1 * time.Minute,
	}
}

type CodeCache interface {
	SetCode(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
	Get(ctx context.Context, id int64) (domain.UserDomain, error)
	Set(ctx context.Context, doU domain.UserDomain) error
}

type UserRedisCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewUserRedisCache(client redis.Cmdable) CodeCache {
	return &UserRedisCache{
		client:     client,
		expiration: 1 * time.Minute,
	}
}
