package tencent

import (
	mySms "GeekProject/newGeekProject/day2/webook/internal/service/sms"
	"context"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"testing"
)

func TestSmsService_Send(t *testing.T) {
	type fields struct {
		appId    *string
		signName *string
		client   *sms.Client
	}
	type args struct {
		ctx     context.Context
		tplID   string
		args    []mySms.NameArg
		numbers []string
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SmsTencentService{
				appId:    tt.fields.appId,
				signName: tt.fields.signName,
				client:   tt.fields.client,
			}
			if err := s.Send(tt.args.ctx, tt.args.tplID, tt.args.args, tt.args.numbers...); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
