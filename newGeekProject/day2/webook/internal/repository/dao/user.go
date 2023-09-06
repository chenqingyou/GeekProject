package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

type UserDao struct {
	db *gorm.DB
}

var (
	ErrUserDuplicate = errors.New("邮箱or手机号码冲突")
	ErrUserNotFound  = gorm.ErrRecordNotFound //框架里自动使用
)

type UserDaoInterface interface {
	FindById(cxt context.Context, Id int64) (UserDB, error)
	FindByEmail(cxt context.Context, email string) (UserDB, error)
	FindByPhone(cxt context.Context, phone string) (UserDB, error)
	InsertUser(cxt context.Context, userDB UserDB) error
	EditUser(cxt context.Context, userDB UserDB) error
}

func NewUserDao(db *gorm.DB) UserDaoInterface {
	return &UserDao{db: db}
}

func (ud *UserDao) FindById(cxt context.Context, Id int64) (UserDB, error) {
	var userDBEmail UserDB
	err := ud.db.WithContext(cxt).First(&userDBEmail, "`id` = ?", Id).Error
	//err := ud.db.WithContext(cxt).Where("email = ?", userDB.Email).First(&userDBEmail).Error
	return userDBEmail, err
}
func (ud *UserDao) FindByEmail(cxt context.Context, email string) (UserDB, error) {
	var userDBEmail UserDB
	err := ud.db.WithContext(cxt).First(&userDBEmail, "email = ?", email).Error
	//err := ud.db.WithContext(cxt).Where("email = ?", userDB.Email).First(&userDBEmail).Error
	return userDBEmail, err
}

func (ud *UserDao) FindByPhone(cxt context.Context, phone string) (UserDB, error) {
	var userDBEmail UserDB
	err := ud.db.WithContext(cxt).First(&userDBEmail, "phone = ?", phone).Error
	//err := ud.db.WithContext(cxt).Where("email = ?", userDB.Email).First(&userDBEmail).Error
	return userDBEmail, err
}

func (ud *UserDao) InsertUser(cxt context.Context, userDB UserDB) error {
	nowTime := time.Now().UnixMilli()
	userDB.CreatTime = nowTime
	userDB.UpdateTime = nowTime
	//写入数据库
	err := ud.db.WithContext(cxt).Create(&userDB).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		//判断是否是唯一
		if mysqlErr.Number == uniqueConflictsErrNo {
			return ErrUserDuplicate
		}
	}
	return err
}

func (ud *UserDao) EditUser(cxt context.Context, userDB UserDB) error {
	nowTime := time.Now().UnixMilli()
	//写入数据库
	err := ud.db.WithContext(cxt).Where("id = ?", userDB.Id).Updates(&UserDB{
		Email:           userDB.Email,
		Password:        userDB.Password,
		Nickname:        userDB.Nickname,
		Birthday:        userDB.Birthday,
		PersonalProfile: userDB.PersonalProfile,
		UpdateTime:      nowTime,
	}).Error
	return err
}

// UserDB 直接对应数据库中的表结构
type UserDB struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 设置为唯一索引
	Email           sql.NullString `gorm:"unique"`
	Password        string
	Nickname        string
	Birthday        string
	PersonalProfile string
	//Phone *string
	Phone      sql.NullString `gorm:"unique"`
	CreatTime  int64
	UpdateTime int64
	DeleteTime int64
}
