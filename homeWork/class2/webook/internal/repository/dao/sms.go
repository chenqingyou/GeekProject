package dao

import (
	mysms "GeekProject/homeWork/class2/webook/internal/service/sms"
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateConditions = errors.New("条件冲突")
)

type SMSDaoInterface interface {
	InsertSMS(cxt context.Context, SmsDB SmsDB) error
	FindByTpl(cxt context.Context, phone string) (SmsDB, error)
}

func NewSMSDao(db *gorm.DB) SMSDaoInterface {
	return &SMSDao{db: db}
}

func (sd *SMSDao) InsertSMS(cxt context.Context, SmsDB SmsDB) error {
	//写入数据库
	nowTime := time.Now().UnixMilli()
	SmsDB.CreatTime = nowTime
	err := sd.db.WithContext(cxt).Create(&SmsDB).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		//判断是否是唯一
		if mysqlErr.Number == uniqueConflictsErrNo {
			return ErrUserDuplicateConditions
		}
	}
	return err
}

func (sd *SMSDao) FindByTpl(cxt context.Context, Tpl string) (SmsDB, error) {
	var req SmsDB
	err := sd.db.WithContext(cxt).First(&req, "tpl = ?", Tpl).Error
	return req, err
}

type SMSDao struct {
	db *gorm.DB
}

// SmsDB 直接对应数据库中的表结构
type SmsDB struct {
	Id        int64           `gorm:"primaryKey,autoIncrement"`
	Tpl       string          `gorm:"unique"`
	NameArg   []mysms.NameArg `gorm:"json"`
	Number    []string
	CreatTime int64
}
