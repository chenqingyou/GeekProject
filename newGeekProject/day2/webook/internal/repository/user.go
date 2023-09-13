package repository

import (
	"GeekProject/newGeekProject/day2/webook/internal/domain"
	"GeekProject/newGeekProject/day2/webook/internal/repository/cache"
	"GeekProject/newGeekProject/day2/webook/internal/repository/dao"
	"context"
	"database/sql"
	"time"
)

type UserRepository struct {
	daoUserDB dao.UserDaoInterface
	cache     cache.UserCache
}

type UserRepositoryInterface interface {
	CreateUser(cxt context.Context, domianU domain.UserDomain) error
	FindByEmail(cxt context.Context, email string) (domain.UserDomain, error)
	FindByPhone(cxt context.Context, phone string) (domain.UserDomain, error)
	EditUser(cxt context.Context, domianU domain.UserDomain) error
	FindById(ctx context.Context, id int64) (domain.UserDomain, error)
}

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicate
	ErrUserNotFound       = dao.ErrUserNotFound
)

// NewUserRepository 不想直接暴露结构体
func NewUserRepository(db dao.UserDaoInterface, cache cache.UserCache) UserRepositoryInterface {
	return &UserRepository{
		daoUserDB: db,
		cache:     cache,
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
	return ur.entityToDomain(byEmail), err
}

func (ur *UserRepository) FindByPhone(cxt context.Context, phone string) (domain.UserDomain, error) {
	by, err := ur.daoUserDB.FindByPhone(cxt, phone)
	if err != nil {
		return domain.UserDomain{}, err
	}
	return ur.entityToDomain(by), err
}

// EditUser 创建用户
func (ur *UserRepository) EditUser(cxt context.Context, domianU domain.UserDomain) error {
	return ur.daoUserDB.EditUser(cxt, ur.domainToEntity(domianU))
}

func (ur *UserRepository) FindById(ctx context.Context, id int64) (domain.UserDomain, error) {
	//先找cache
	getU, err := ur.cache.Get(ctx, id)
	if err == nil {
		//缓存里面有数据
		return getU, err
	}
	//if err == cache.ErrUserNotNoFound {
	//	//缓存里面没有，去数据库加载
	//}
	//if err != nil {
	//
	//}
	// 这里怎么办？ err = io.EOF
	// 要不要去数据库加载？
	// 看起来我不应该加载？
	// 看起来我好像也要加载？
	ud, err := ur.daoUserDB.FindById(ctx, id)
	if err != nil {
		return domain.UserDomain{}, err
	}
	getU = ur.entityToDomain(ud)
	//缓存会导致数据不一致
	//go func() {
	err = ur.cache.Set(ctx, getU)
	if err != nil {
		//打日志做监控
		//return domain.UserDomain{}, err
	}
	//	}()
	return getU, nil
	// 选加载 —— 做好兜底，万一 Redis 真的崩了，你要保护住你的数据库
	// 我数据库限流呀！
	// 选不加载 —— 用户体验差一点
	// 缓存里面有数据
	// 缓存里面没有数据
	// 缓存出错了，你也不知道有没有数据

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
