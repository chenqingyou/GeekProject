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

// EditUser 创建用户
func (ur *UserRepository) EditUser(cxt context.Context, domianU domain.UserDomain) error {
	return ur.daoUserDB.EditUser(cxt, dao.UserDB{
		Id:              domianU.Id,
		Email:           domianU.Email,
		Password:        domianU.Password,
		Nickname:        domianU.Nickname,
		Birthday:        domianU.Birthday,
		PersonalProfile: domianU.PersonalProfile,
	})
}

func (ur *UserRepository) FindById(ctx context.Context, id int64) (domain.UserDomain, error) {
	ud, err := ur.daoUserDB.FindById(ctx, id)
	if err != nil {
		return domain.UserDomain{}, err
	}
	return domain.UserDomain{
		Email:           ud.Email,
		Nickname:        ud.Nickname,
		Birthday:        ud.Birthday,
		PersonalProfile: ud.PersonalProfile,
	}, nil
}
