package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
)

func (lc *UserLocalCache) SetCode(ctx context.Context, biz, phone, code string) error {
	res := checkAndSetCode(lc.key(biz, phone), code)
	switch res {
	case 0:
		return nil
	case -1:
		return ErrSetCodeFrequently
	default:
		return errors.New("系统错误")
	}
}

func (lc *UserLocalCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func (lc *UserLocalCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	res := checkCode(lc.key(biz, phone), inputCode)
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
