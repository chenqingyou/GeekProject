package failover

import (
	"GeekProject/newGeekProject/day2/webook/internal/service/sms"
	"context"
	"errors"
	"sync/atomic"
)

type SMSServiceFailOver struct {
	svcs []sms.ServiceSmsInterface
	idx  uint64
}

func NewSMSServiceFailOver(svcs []sms.ServiceSmsInterface) sms.ServiceSmsInterface {
	return &SMSServiceFailOver{svcs: svcs}
}

func (f *SMSServiceFailOver) Send(ctx context.Context, tpl string,
	args []sms.NameArg, numbers ...string) error {
	/*
		缺点： • 每次都从头开始轮询，绝大多数请求会在 svcs[0] 就成功，负载不均衡。
		• 如果 svcs 有几十个，轮询都很慢。*/
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tpl, args, numbers...)
		if err == nil {
			return err
		}
		//正常输入日志和健康
	}
	return errors.New("服务商全部都失败了")
}

func (f *SMSServiceFailOver) SendV1(ctx context.Context, tpl string,
	args []sms.NameArg, numbers ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < idx+length; i++ {
		svc := f.svcs[int(i%length)]
		err := svc.Send(ctx, tpl, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			return err
		}

	}
	return errors.New("服务商全部都失败了")
}
