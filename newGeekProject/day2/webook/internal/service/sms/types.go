package sms

import "context"

type ServiceSmsInterface interface {
	Send(ctx context.Context, tpl string, args []NameArg, numbers ...string) error
}

type NameArg struct {
	Val  string
	Name string
}
