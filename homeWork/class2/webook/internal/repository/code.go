package repository

import (
	"GeekProject/homeWork/class2/webook/internal/repository/cache"
	"context"
)

var (
	ErrSetCodeFrequently    = cache.ErrSetCodeFrequently
	ErrVerifyCodeFrequently = cache.ErrVerifyCodeFrequently
	ErrUnknowForCode        = cache.ErrUnknowForCode
)

type CodeRepositoryInterface interface {
	SetCode(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

type CodeRepository struct {
	codeCache cache.CodeCache
}

func NewCodeRepository(codeCache cache.CodeCache) CodeRepositoryInterface {
	return &CodeRepository{codeCache: codeCache}
}

func (cr *CodeRepository) SetCode(ctx context.Context, biz, phone, code string) error {
	err := cr.codeCache.SetCode(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	return nil
}

func (cr *CodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return cr.codeCache.Verify(ctx, biz, phone, code)
}
