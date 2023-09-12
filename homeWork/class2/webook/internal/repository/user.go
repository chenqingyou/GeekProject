package repository

import (
	"GeekProject/homeWork/class2/webook/internal/domain"
	"GeekProject/homeWork/class2/webook/internal/repository/cache"
	"GeekProject/homeWork/class2/webook/internal/repository/dao"
	"context"
	"database/sql"
	"time"
)

type UserRepository struct {
	daoUserDB      dao.UserDaoInterface
	userLocalCache cache.CodeCache
}

type UserRepositoryInterface interface {
	CreateUser(cxt context.Context, domianU domain.UserDomain) error
	FindByEmail(cxt context.Context, email string) (domain.UserDomain, error)
	EditUser(cxt context.Context, domianU domain.UserDomain) error
	FindById(ctx context.Context, id int64) (domain.UserDomain, error)
	FindByPhone(cxt context.Context, phone string) (domain.UserDomain, error)
}

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

// NewUserRepository 不想直接暴露结构体
func NewUserRepository(db dao.UserDaoInterface, localCache cache.CodeCache) UserRepositoryInterface {
	return &UserRepository{
		daoUserDB:      db,
		userLocalCache: localCache,
	}
}

// CreateUser 创建用户
func (ur *UserRepository) CreateUser(cxt context.Context, domianU domain.UserDomain) error {
	return ur.daoUserDB.InsertUser(cxt, ur.domainToEntity(domianU))
}

func (ur *UserRepository) FindByEmail(cxt context.Context, email string) (domain.UserDomain, error) {
	byEmail, err := ur.daoUserDB.FindByEmail(cxt, email)
	if err != nil {
		return domain.UserDomain{}, err
	}
	return domain.UserDomain{
		Id:       byEmail.Id,
		Email:    byEmail.Email.String,
		Password: byEmail.Password,
	}, err
}

// EditUser 创建用户
func (ur *UserRepository) EditUser(cxt context.Context, domianU domain.UserDomain) error {
	//缓存数据到本地
	err := ur.userLocalCache.Set(cxt, domianU)
	if err != nil {
		return err
	}
	return ur.daoUserDB.EditUser(cxt, ur.domainToEntity(domianU))
}

func (ur *UserRepository) FindById(ctx context.Context, id int64) (domain.UserDomain, error) {
	//先去缓存里面找 找不到再去数据库里面查
	getId, err := ur.userLocalCache.Get(ctx, id)
	if err == nil {
		return getId, nil
	}
	ud, err := ur.daoUserDB.FindById(ctx, id)
	if err != nil {
		return domain.UserDomain{}, err
	}
	return domain.UserDomain{
		Email:           ud.Email.String,
		Nickname:        ud.Nickname,
		Birthday:        ud.Birthday,
		PersonalProfile: ud.PersonalProfile,
	}, nil
}

func (ur *UserRepository) FindByPhone(cxt context.Context, phone string) (domain.UserDomain, error) {
	by, err := ur.daoUserDB.FindByPhone(cxt, phone)
	if err != nil {
		return domain.UserDomain{}, err
	}
	return ur.entityToDomain(by), err
}

func (ur *UserRepository) entityToDomain(u dao.UserDB) domain.UserDomain {
	return domain.UserDomain{
		Id:              u.Id,
		Email:           u.Email.String,
		Phone:           u.Phone.String,
		Password:        u.Password,
		Nickname:        u.Nickname,
		Birthday:        u.Birthday,
		PersonalProfile: u.PersonalProfile,
		Ctime:           time.UnixMilli(u.CreatTime),
	}
}

func (ur *UserRepository) domainToEntity(u domain.UserDomain) dao.UserDB {
	return dao.UserDB{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password:        u.Password,
		Nickname:        u.Nickname,
		Birthday:        u.Birthday,
		PersonalProfile: u.PersonalProfile,
		CreatTime:       0,
	}
}
