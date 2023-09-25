package sms

import "context"

type ServiceSmsInterface interface {
	Send(ctx context.Context, biz string, args []NameArg, numbers ...string) error
}

type NameArg struct {
	Val  string
	Name string
}
