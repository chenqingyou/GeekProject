package cache

import (
	"GeekProject/newGeekProject/day2/webook/internal/repository/cache/redismocks"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRedisCodeCache_Set(t *testing.T) {
	testCase := []struct {
		name             string
		mock             func(ctrl *gomock.Controller) redis.Cmdable
		biz, code, phone string
		wantErr          error
	}{
		{
			name: "验证码设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetErr(nil)
				res.SetVal(int64(0))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{fmt.Sprintf("phone_code:%s:%s", "login", "1212312312")}, []any{"091231"}).Return(res)
				return cmd
			},
			biz:     "login",
			phone:   "1212312312",
			code:    "091231",
			wantErr: nil,
		},
		{
			name: "redis错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetErr(errors.New("mock redis err"))
				res.SetVal(int64(0))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{fmt.Sprintf("phone_code:%s:%s", "login", "1212312312")}, []any{"091231"}).Return(res)
				return cmd
			},
			biz:     "login",
			phone:   "1212312312",
			code:    "091231",
			wantErr: errors.New("mock redis err"),
		},
		{
			name: "code发送频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetErr(ErrSetCodeFrequently)
				res.SetVal(int64(-1))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{fmt.Sprintf("phone_code:%s:%s", "login", "1212312312")}, []any{"091231"}).Return(res)
				return cmd
			},
			biz:     "login",
			phone:   "1212312312",
			code:    "091231",
			wantErr: ErrSetCodeFrequently,
		},
		// TODO: Add test cases.
	}
	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			rcc := NewCodeCache(tt.mock(ctrl))
			err := rcc.Set(context.Background(), tt.biz, tt.phone, tt.code)
			assert.Equal(t, err, tt.wantErr)
		})
	}
}
