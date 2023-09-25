package auth

import (
	"GeekProject/newGeekProject/day2/webook/internal/service/sms"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

// SmsService 装饰器用于来短信验证服务的权限控制
type SmsService struct {
	svc sms.ServiceSmsInterface
	key string
}

func (s SmsService) Send(ctx context.Context, biz string, args []sms.NameArg, numbers ...string) error {
	//在这儿使用权限控制，用biz里面的token
	var claims Claims
	//进行token验证,如果这儿能解析成功，说明就是对应的业务方
	//ParseWithClaims 一定要使用指针， 因为ParseWithClaims会去修改claims里面的值
	token, err := jwt.ParseWithClaims(biz, claims, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("token 不合法")
	}
	return s.svc.Send(ctx, biz, args, numbers...)
}

type Claims struct {
	jwt.RegisteredClaims
	Tpl string
}
