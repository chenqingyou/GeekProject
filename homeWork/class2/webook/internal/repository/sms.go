package repository

import (
	"GeekProject/homeWork/class2/webook/internal/domain"
	"GeekProject/homeWork/class2/webook/internal/repository/dao"
	"context"
)

type SMSRepositoryInterface interface {
	InsertSMS(ctx context.Context, smsDomain domain.SmsDomain) error
	SendSMS(ctx context.Context, tpl string) (domain.SmsDomain, error)
}

type SMSRepository struct {
	smsDao dao.SMSDao
}

func NewSMSRepository(smsDao dao.SMSDao) SMSRepositoryInterface {
	return &SMSRepository{smsDao: smsDao}
}

func (sr *SMSRepository) InsertSMS(ctx context.Context, smsDomain domain.SmsDomain) error {
	err := sr.smsDao.InsertSMS(ctx, dao.SmsDB{
		Tpl:     smsDomain.Tpl,
		NameArg: smsDomain.NameArg,
		Number:  smsDomain.Numbers,
	})
	if err != nil {
		return err
	}
	return nil
}

func (sr *SMSRepository) SendSMS(ctx context.Context, tpl string) (domain.SmsDomain, error) {
	byPhoneReq, err := sr.smsDao.FindByTpl(ctx, tpl)
	if err != nil {
		return domain.SmsDomain{}, err
	}
	return domain.SmsDomain{
		Tpl:     byPhoneReq.Tpl,
		NameArg: byPhoneReq.NameArg,
		Numbers: byPhoneReq.Number,
	}, nil
}
