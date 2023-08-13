package repository

import (
	"GeekProject/homeWork/class2/webook/internal/domain"
	"GeekProject/homeWork/class2/webook/internal/repository/dao"
	"context"
)

type UserRepository struct {
	daoUserDB *dao.UserDao
}

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

// NewUserRepository 不想直接暴露结构体
func NewUserRepository(db *dao.UserDao) *UserRepository {
	return &UserRepository{
		daoUserDB: db,
	}
}

// CreateUser 创建用户
func (ur *UserRepository) CreateUser(cxt context.Context, domianU domain.UserDomain) error {
	return ur.daoUserDB.InsertUser(cxt, dao.UserDB{
		Email:    domianU.Email,
		Password: domianU.Password,
	})
}

func (ur *UserRepository) FindByEmail(cxt context.Context, email string) (domain.UserDomain, error) {
	byEmail, err := ur.daoUserDB.FindByEmail(cxt, email)
	if err != nil {
		return domain.UserDomain{}, err
	}
	return domain.UserDomain{
		Id:       byEmail.Id,
		Email:    byEmail.Email,
		Password: byEmail.Password,
	}, err
}
