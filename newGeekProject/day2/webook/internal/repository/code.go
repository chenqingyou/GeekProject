package repository

import (
	"GeekProject/newGeekProject/day2/webook/internal/repository/cache"
	"context"
)

var (
	ErrSetCodeFrequently    = cache.ErrSetCodeFrequently
	ErrVerifyCodeFrequently = cache.ErrVerifyCodeFrequently
)

type CodeRepositoryInterface interface {
	Store(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

type CodeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(cache cache.CodeCache) CodeRepositoryInterface {
	return &CodeRepository{cache: cache}
}

func (cr *CodeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return cr.cache.Set(ctx, biz, phone, code)
}

func (cr *CodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return cr.cache.Verify(ctx, biz, phone, code)
}
