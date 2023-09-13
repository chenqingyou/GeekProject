package service

import (
	"GeekProject/newGeekProject/day2/webook/internal/domain"
	"GeekProject/newGeekProject/day2/webook/internal/repository"
	reposvcmocks "GeekProject/newGeekProject/day2/webook/internal/repository/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func TestEncrypt(t *testing.T) {
	passwd := "hello#world"
	password, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	t.Log(string(password))
	bcrypt.CompareHashAndPassword(password, []byte(passwd))
	assert.NoError(t, err)
}

func TestUserService_Login(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepositoryInterface
		//输入
		user domain.UserDomain
		//输出
		want    domain.UserDomain
		wantErr error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepositoryInterface {
				userRepo := reposvcmocks.NewMockUserRepositoryInterface(ctrl)
				userRepo.EXPECT().FindByEmail(gomock.Any(), "123123@qq.com").Return(domain.UserDomain{
					Email:    "123123@qq.com",
					Password: "$2a$10$s8Rq1Y1Fm4jDRJ3PPxnd1.94k5x1HQg7vhJeu53qBWlNk0BUBorV2",
					Ctime:    now,
				}, nil)
				return userRepo
			},
			user: domain.UserDomain{
				Email:    "123123@qq.com",
				Password: "hello#world",
			},
			want: domain.UserDomain{
				Email:    "123123@qq.com",
				Password: "$2a$10$s8Rq1Y1Fm4jDRJ3PPxnd1.94k5x1HQg7vhJeu53qBWlNk0BUBorV2",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepositoryInterface {
				userRepo := reposvcmocks.NewMockUserRepositoryInterface(ctrl)
				userRepo.EXPECT().FindByEmail(gomock.Any(), "123123@qq.com").Return(domain.UserDomain{}, repository.ErrUserNotFound)
				return userRepo
			},
			user: domain.UserDomain{
				Email:    "123123@qq.com",
				Password: "hello#world",
			},
			want:    domain.UserDomain{},
			wantErr: ErrInvalidUserOrPassword,
		},
		{
			name: "DB错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepositoryInterface {
				userRepo := reposvcmocks.NewMockUserRepositoryInterface(ctrl)
				userRepo.EXPECT().FindByEmail(gomock.Any(), "123123@qq.com").Return(domain.UserDomain{}, errors.New("mock db 错误"))
				return userRepo
			},
			user: domain.UserDomain{
				Email:    "123123@qq.com",
				Password: "hello#world",
			},
			want:    domain.UserDomain{},
			wantErr: errors.New("mock db 错误"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			us := NewUserService(tt.mock(ctrl))
			got, err := us.Login(context.Background(), tt.user)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
