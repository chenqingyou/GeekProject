package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	ErrSetCodeFrequently    = errors.New("发送验证码太频繁")
	ErrVerifyCodeFrequently = errors.New("验证次数太多")
	ErrUnknowForCode        = errors.New("未知错误")
)

//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type RedisCodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		client: client,
	}
}

func (cc *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := cc.client.Eval(ctx, luaSetCode, []string{cc.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		return ErrSetCodeFrequently
	default:
		return errors.New("系统错误")
	}
}

func (cc *RedisCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func (cc *RedisCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	res, err := cc.client.Eval(ctx, luaVerifyCode, []string{cc.key(biz, phone)}, inputCode).Int()
	//res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		return false, ErrVerifyCodeFrequently
	case -2:
		return false, nil
	default:
		return false, ErrUnknowForCode
	}
}
