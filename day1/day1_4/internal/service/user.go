package service

import (
	"GeekProject/day1/day1_4/internal/domain"
	"GeekProject/day1/day1_4/internal/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicateEmailRepSvc = repository.ErrUserDuplicateEmailRep
	ErrUserNotFoundSvc          = repository.ErrUserNotFoundRep
)

type UserService struct {
	userRep *repository.UserRep
}

func NewUserService(repU *repository.UserRep) *UserService {
	return &UserService{
		userRep: repU,
	}
}

func (us *UserService) CreateUser(ctx *gin.Context, domainU domain.UserDomain) error {
	//对密码进行加密
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(domainU.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = us.userRep.CreateUser(ctx, domain.UserDomain{
		Email:    domainU.Email,
		Password: string(bcryptPassword),
	})
	if err == ErrUserDuplicateEmailRepSvc {
		return ErrUserDuplicateEmailRepSvc
	}
	if err != nil {
		return err
	}
	return err
}

func (us *UserService) LoginUser(ctx *gin.Context, domainU domain.UserDomain) (domain.UserDomain, error) {
	//根据email去查询，这个用户是否存在
	bcryptPassword, err := us.userRep.FindByEmail(ctx, domainU.Email)
	if err == ErrUserNotFoundSvc {
		return domain.UserDomain{}, ErrUserNotFoundSvc
	}
	if err != nil {
		return domain.UserDomain{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(bcryptPassword.Password), []byte(domainU.Password))
	if err != nil {
		return domain.UserDomain{}, err
	}
	return domain.UserDomain{
		ID: bcryptPassword.ID,
	}, err
}

func (us *UserService) EditUser(ctx *gin.Context, domainU domain.UserDomain) error {
	//根据email去查询，这个用户是否存在
	return us.userRep.Edit(ctx, domainU)
}
