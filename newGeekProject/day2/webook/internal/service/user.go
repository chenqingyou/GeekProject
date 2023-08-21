package service

import (
	"GeekProject/newGeekProject/day2/webook/internal/domain"
	"GeekProject/newGeekProject/day2/webook/internal/repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicateEmail //这样的用法可以保证耦合度不高，只和下层耦合
	ErrInvalidUserOrPassword = errors.New("账号/密码错误")
)

// UserService 服务端，调用repository
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService 不想直接暴露结构体
func NewUserService(rep *repository.UserRepository) *UserService {
	return &UserService{
		repo: rep,
	}
}

func (us *UserService) SignUp(ctx context.Context, domainU domain.UserDomain) error {
	//密码加密
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(domainU.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	//存放密码
	domainU.Password = string(hashPassword)
	return us.repo.CreateUser(ctx, domainU)
}

func (us *UserService) Login(ctx context.Context, domainU domain.UserDomain) (domain.UserDomain, error) {
	//先查找这个用户是否存在，然后在比较密码对不对
	userMessage, err := us.repo.FindByEmail(ctx, domainU.Email)
	//定义错误码，如果是邮箱不存在
	if err == repository.ErrUserNotFound {
		return domain.UserDomain{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.UserDomain{}, err
	}
	//比较密码
	err = bcrypt.CompareHashAndPassword([]byte(userMessage.Password), []byte(domainU.Password))
	if err != nil {
		return domain.UserDomain{}, ErrInvalidUserOrPassword
	}
	return userMessage, err
}

func (us *UserService) Profile(ctx context.Context, id int64) (domain.UserDomain, error) {
	return us.repo.FindById(ctx, id)
}

func (us *UserService) FindOrCreate(ctx context.Context, phone string) (domain.UserDomain, error) {
	u, err := us.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		// 绝大部分请求进来这里
		// nil 会进来这里
		// 不为 ErrUserNotFound 的也会进来这里
		return u, err
	}
	//代表这个用户没有注册
	u = domain.UserDomain{
		Phone: phone,
	}
	//注册一次
	err = us.repo.CreateUser(ctx, u)
	if err != nil && err != repository.ErrUserNotFound {
		return u, err
	}
	// 因为这里会遇到主从延迟的问题
	return us.repo.FindByPhone(ctx, phone) //在查找一次返回id
}
