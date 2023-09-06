package memory

import (
	mySms "GeekProject/newGeekProject/day2/webook/internal/service/sms"
	"context"
	"fmt"
)

type ServiceSmsInterface struct {
}

func NewService() *ServiceSmsInterface {
	return &ServiceSmsInterface{}
}

func (s *ServiceSmsInterface) Send(ctx context.Context, tplID string, args []mySms.NameArg, numbers ...string) error {
	fmt.Println(args)
	return nil
}
