package ioc

import (
	"GeekProject/newGeekProject/day2/webook/internal/service/sms"
	"GeekProject/newGeekProject/day2/webook/internal/service/sms/memory"
)

func InitSMSService() sms.ServiceSmsInterface {
	return memory.NewService()
}
