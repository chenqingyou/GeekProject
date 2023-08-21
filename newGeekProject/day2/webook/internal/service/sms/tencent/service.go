package tencent

import (
	mySms "GeekProject/newGeekProject/day2/webook/internal/service/sms"
	"context"
	"fmt"
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	//sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type SmsService struct {
	appId    *string
	signName *string
	client   *sms.Client
}

func NewSmsService(client *sms.Client, appId, signName string) *SmsService {
	return &SmsService{
		client:   client,
		appId:    ekit.ToPtr[string](appId),
		signName: ekit.ToPtr[string](signName),
	}
}

func (s SmsService) Send(ctx context.Context, tplID string, args []mySms.NameArg, numbers ...string) error {
	smsReq := sms.NewSendSmsRequest()
	smsReq.SignName = s.signName
	smsReq.SmsSdkAppId = s.signName
	smsReq.TemplateId = ekit.ToPtr(tplID)
	smsReq.PhoneNumberSet = s.ToSliceFunc(numbers)
	smsReq.TemplateParamSet = slice.Map[mySms.NameArg, *string](args, func(idx int, src mySms.NameArg) *string {
		return &src.Val
	})
	sendSms, err := s.client.SendSms(smsReq)
	if err != nil {
		return err
	}
	for _, status := range sendSms.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) == "Ok" {
			return fmt.Errorf("Failed to send SMS messages[%v],[%v]\n", *status.Message, *status.Code)
		}
	}
	return err
}

func (s SmsService) ToSliceFunc(args []string) []*string {
	return slice.Map[string, *string](args, func(idx int, src string) *string {
		return &src
	})
}
