package ratelimit

import (
	"GeekProject/homeWork/class2/webook/internal/service/sms"
	"GeekProject/homeWork/class2/webook/pkg/ratelimit_win"
	"context"
	"fmt"
)

//使用本地缓存实现限流器，并且加入短信容错机制

var errLimit = fmt.Errorf("触发了限流")

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
		return fmt.Errorf("判断是否限流出现问题%err", err)
	}
	//这里使用本地缓存，不会返回错误
	if limited {
		return errLimit
	}
	err = s.svc.Send(ctx, tpl, args, numbers...)
	return err
}

/*
设计一个新的容错机制，同步转异步的容错机制。当满足以下两个条件中的任何一个时，将请求转储到数据库，后续再另外启动一个 goroutine异步发送出去。
触发了限流。
判定服务商已经崩溃。
要求:
如何判定服务商已经崩溃，不允许使用课程上的判断机制，你需要设计一个新的判断机制，并且解释这种判定机制的决策理由。
控制异步重试次数，转储到数据库之后，可以重试 N 次，重试间隔你可以自由决策。
不允许写死任何参数，即用户必须可以控制控制参数。
保持面向接口和依赖注入风格。
写明这种容错机制适合什么场景，并且有什么优缺点
针对提出的缺点，写出后续的改进方案。
提供单元测试。
*/
