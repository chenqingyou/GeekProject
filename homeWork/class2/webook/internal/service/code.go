package service

import (
	"GeekProject/homeWork/class2/webook/internal/repository"
	"GeekProject/homeWork/class2/webook/internal/service/sms"
	"context"
	"fmt"
	"math/rand"
)

const codeTplId = "1937434"

type CodeServiceInterface interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
	GenerateVerificationCode() string
}

type CodeService struct {
	codeRep repository.CodeRepositoryInterface
	smsSvr  sms.ServiceSmsInterface
}

var (
	ErrSetCodeFrequently    = repository.ErrSetCodeFrequently
	ErrVerifyCodeFrequently = repository.ErrVerifyCodeFrequently
)

func NewCodeService(codeRep repository.CodeRepositoryInterface, smsSvr sms.ServiceSmsInterface) CodeServiceInterface {
	return &CodeService{
		codeRep: codeRep,
		smsSvr:  smsSvr,
	}
}

// Send 发送验证码，我需要什么参数？
func (svc *CodeService) Send(ctx context.Context, biz string, phone string) error {
	//两个步骤，生成一个验证码
	code := svc.GenerateVerificationCode()
	//塞进去
	err := svc.codeRep.SetCode(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	err = svc.smsSvr.Send(ctx, codeTplId, []sms.NameArg{
		{
			Val: code,
		},
	}, phone)
	//if err != nil {
	// 这个地方怎么办？
	// 这意味着， 有这个验证码，但是不好意思，
	// 我能不能删掉这个验证码？
	// 你这个 err 可能是超时的 err，你都不知道，发出了没
	// 在这里重试
	// 要重试的话，初始化的时候，传入一个自己就会重试的 smsSvc
	//}
	return err
}

func (svc *CodeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return svc.codeRep.Verify(ctx, biz, phone, inputCode)
}

func (svc *CodeService) GenerateVerificationCode() string {
	numCode := rand.Intn(100000)
	//不足6位，用0补足
	return fmt.Sprintf("%06d", numCode)
}
