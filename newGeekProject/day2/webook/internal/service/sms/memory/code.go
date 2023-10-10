package memory

import (
	mySms "GeekProject/newGeekProject/day2/webook/internal/service/sms"
	"context"
	"fmt"
)

type ServiceMemorySmsInterface struct {
}

func NewService() *ServiceMemorySmsInterface {
	return &ServiceMemorySmsInterface{}
}

func (s *ServiceMemorySmsInterface) Send(ctx context.Context, tplID string, args []mySms.NameArg, numbers ...string) error {
	fmt.Println(args)
	return nil
}
