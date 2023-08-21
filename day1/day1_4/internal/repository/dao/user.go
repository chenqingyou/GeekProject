package dao

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound //框架里自动使用
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (ud *UserDao) InsertUser(ctx *gin.Context, userDB UserDB) error {
	newTime := time.Now().UnixMilli()
	userDB.CreatTime = newTime
	err := ud.db.WithContext(ctx).Create(&userDB).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		//判断是否是唯一
		if mysqlErr.Number == uniqueConflictsErrNo {
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (ud *UserDao) FindByEmail(ctx *gin.Context, email string) (UserDB, error) {
	var userEmail UserDB
	err := ud.db.WithContext(ctx).First(&userEmail, "email = ?", email).Error
	return userEmail, err
}

func (ud *UserDao) FindByID(ctx *gin.Context, id int64) (UserDB, error) {
	var userDetail UserDB
	err := ud.db.WithContext(ctx).First("id = ?", id).Error
	return userDetail, err
}

func (ud *UserDao) EditUser(ctx *gin.Context, userDB UserDB) error {
	// 创建一个结构体用于更新
	updatedUser := UserDB{
		Email:           userDB.Email,
		Password:        userDB.Password,
		Nickname:        userDB.Nickname,
		Birthday:        userDB.Birthday,
		PersonalProfile: userDB.PersonalProfile,
		UpdateTime:      time.Now().UnixMilli(),
	}
	return ud.db.WithContext(ctx).Where("id = ?", userDB.Id).Updates(updatedUser).Error
}

// UserDB 直接对应数据库中的表结构
type UserDB struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 设置为唯一索引
	Email           string `gorm:"unique"`
	Password        string
	Nickname        string `gorm:"size:10"`
	Birthday        string
	PersonalProfile string `gorm:"size:300"`
	//Phone *string
	Phone      sql.NullString `gorm:"unique"`
	CreatTime  int64
	UpdateTime int64
	DeleteTime int64
}
