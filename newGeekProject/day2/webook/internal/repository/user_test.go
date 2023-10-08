package repository

import (
	"GeekProject/newGeekProject/day2/webook/internal/domain"
	"GeekProject/newGeekProject/day2/webook/internal/repository/cache"
	cachemocks "GeekProject/newGeekProject/day2/webook/internal/repository/cache/mocks"
	"GeekProject/newGeekProject/day2/webook/internal/repository/dao"
	daomocks "GeekProject/newGeekProject/day2/webook/internal/repository/dao/mocks"
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestUserRepository_FindById(t *testing.T) {
	now := time.Now()
	//去掉毫秒以外的地方
	now = time.UnixMilli(now.UnixMilli())
	tests := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (dao.UserDaoInterface, cache.UserCache)
		id      int64
		want    domain.UserDomain
		wantErr error
	}{
		{
			name: "缓存没有命中，查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDaoInterface, cache.UserCache) {
				uc := cachemocks.NewMockUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.UserDomain{}, cache.ErrUserNotNoFound)
				ud := daomocks.NewMockUserDaoInterface(ctrl)
				ud.EXPECT().FindById(gomock.Any(), int64(123)).Return(dao.UserDB{
					Id: 123,
					Email: sql.NullString{
						String: "123@qq.com",
						Valid:  true,
					},
					Password: "123123123",
					Phone: sql.NullString{
						String: "123@qq.com",
						Valid:  true,
					},
					CreatTime:  now.UnixMilli(),
					UpdateTime: now.UnixMilli(),
				}, nil)
				uc.EXPECT().Set(gomock.Any(), domain.UserDomain{
					Id:       123,
					Email:    "123@qq.com",
					Password: "123123123",
					Phone:    "123@qq.com",
					Ctime:    now,
				}).Return(nil)
				return ud, uc
			},
			id: 123,
			want: domain.UserDomain{
				Id:       123,
				Email:    "123@qq.com",
				Password: "123123123",
				Phone:    "123@qq.com",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "缓存命中，查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDaoInterface, cache.UserCache) {
				uc := cachemocks.NewMockUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.UserDomain{
					Id:       123,
					Email:    "123@qq.com",
					Password: "123123123",
					Phone:    "123@qq.com",
					Ctime:    now,
				}, nil)
				ud := daomocks.NewMockUserDaoInterface(ctrl)
				return ud, uc
			},
			id: 123,
			want: domain.UserDomain{
				Id:       123,
				Email:    "123@qq.com",
				Password: "123123123",
				Phone:    "123@qq.com",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "缓存没有命中，查询失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDaoInterface, cache.UserCache) {
				uc := cachemocks.NewMockUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.UserDomain{}, cache.ErrUserNotNoFound)
				ud := daomocks.NewMockUserDaoInterface(ctrl)
				ud.EXPECT().FindById(gomock.Any(), int64(123)).
					Return(dao.UserDB{}, errors.New("mock db 错误"))
				return ud, uc
			},
			id:      123,
			want:    domain.UserDomain{},
			wantErr: errors.New("mock db 错误"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ur := NewUserRepository(tt.mock(ctrl))
			got, err := ur.FindById(context.Background(), tt.id)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
			time.Sleep(time.Second)
		})
	}
}
