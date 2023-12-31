package ratelimit

import (
	"GeekProject/newGeekProject/day2/webook/internal/service/sms"
	"GeekProject/newGeekProject/day2/webook/pkg/ratelimit_win"
	"context"
	"fmt"
)

var errLimit = fmt.Errorf("短信服务触发了限流")

type LimitSMSService struct {
	svc     sms.ServiceSmsInterface      //被装饰的对象
	limiter ratelimit_win.LimitInterface //装饰器模式
}

func NewLimitSMSService(svc sms.ServiceSmsInterface, limiter ratelimit_win.LimitInterface) sms.ServiceSmsInterface {
	return &LimitSMSService{
		svc:     svc,
		limiter: limiter,
	}
}

func (s *LimitSMSService) Send(ctx context.Context, tpl string, args []sms.NameArg, numbers ...string) error {
	//可以加一些代码，新特性
	limited, err := s.limiter.Limited(ctx, "sms:tencent")
	if err != nil {
		//系统错误，一般是redis崩溃了
		//可以限流：保守策略，你的下游很坑 性能太差
		//可以不限制，容错策略,业务可用性要求很高
		return fmt.Errorf("短信服务判断是否限流出现问题%err", err)
	}
	if limited {
		return errLimit
	}
	err = s.svc.Send(ctx, tpl, args, numbers...)
	return err
}
