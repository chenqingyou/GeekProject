package retryable

import (
	"GeekProject/newGeekProject/day2/webook/internal/service/sms"
	"context"
	"errors"
)

// 短信服务的重试
type SMSService struct {
	svc         sms.ServiceSmsInterface
	retryMaxCnt int
}

func (s *SMSService) Send(ctx context.Context, biz string, args []sms.NameArg, numbers ...string) error {
	err := s.svc.Send(ctx, biz, args, numbers...)
	cnt := 1
	for err != nil && cnt < s.retryMaxCnt {
		err = s.svc.Send(ctx, biz, args, numbers...)
		if err == nil {
			return nil
		}
		cnt++
	}
	return errors.New("重试次数用完了")
}
