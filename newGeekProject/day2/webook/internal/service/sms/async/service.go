package async

import (
	"GeekProject/newGeekProject/day2/webook/internal/service/sms"
	"context"
)

type SMSService struct {
	svc sms.ServiceSmsInterface
}

func (s *SMSService) Send(ctx context.Context, biz string, args []sms.NameArg, numbers ...string) error {
	//1、正常情况
	//2、服务崩了(怎么判断服务崩溃)
	//3、服务崩了，把短信转到数据库存储
	//4、数据库发送，重试次数
	//5、怎么切换到正常的服务，怎么判断

	//先正常路径
	return nil
}
