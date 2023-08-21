package repository

import (
	"GeekProject/day1/day1_4/internal/domain"
	"GeekProject/day1/day1_4/internal/repository/dao"
	"github.com/gin-gonic/gin"
)

var (
	ErrUserDuplicateEmailRep = dao.ErrUserDuplicateEmail
	ErrUserNotFoundRep       = dao.ErrUserNotFound
)

type UserRep struct {
	userDao *dao.UserDao
}

func NewUserRep(db *dao.UserDao) *UserRep {
	return &UserRep{
		userDao: db,
	}
}

func (ur *UserRep) CreateUser(ctx *gin.Context, domainU domain.UserDomain) error {
	err := ur.userDao.InsertUser(ctx, dao.UserDB{
		Email:    domainU.Email,
		Password: domainU.Password,
	})
	if err == ErrUserDuplicateEmailRep {
		return ErrUserDuplicateEmailRep
	}
	if err != nil {
		return err
	}
	return err
}

func (ur *UserRep) FindByEmail(ctx *gin.Context, email string) (domain.UserDomain, error) {
	ud, err := ur.userDao.FindByEmail(ctx, email)
	if err == ErrUserNotFoundRep {
		return domain.UserDomain{}, ErrUserNotFoundRep
	}
	if err != nil {
		return domain.UserDomain{}, ErrUserNotFoundRep
	}
	return domain.UserDomain{
		Password: ud.Password,
		ID:       ud.Id,
	}, err
}

func (ur *UserRep) FindById(ctx *gin.Context, userDomain domain.UserDomain) error {
	return ur.userDao.EditUser(ctx, dao.UserDB{
		Email:           userDomain.Email,
		Password:        userDomain.Password,
		Nickname:        userDomain.Nickname,
		Birthday:        userDomain.Birthday,
		PersonalProfile: userDomain.PersonalProfile,
	})
}

func (ur *UserRep) Edit(ctx *gin.Context, userDomain domain.UserDomain) error {
	return ur.userDao.EditUser(ctx, dao.UserDB{
		Id:              userDomain.ID,
		Email:           userDomain.Email,
		Password:        userDomain.Password,
		Nickname:        userDomain.Nickname,
		Birthday:        userDomain.Birthday,
		PersonalProfile: userDomain.PersonalProfile,
	})
}
