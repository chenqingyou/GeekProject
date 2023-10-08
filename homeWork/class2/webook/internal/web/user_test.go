package web

import (
	"GeekProject/homeWork/class2/webook/internal/domain"
	"GeekProject/homeWork/class2/webook/internal/service"
	svcmocks "GeekProject/homeWork/class2/webook/internal/service/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_LonginJwt(t *testing.T) {
	testsCase := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserServiceInterface
		body     string
		wantCode int
		wantBody domain.Result
	}{
		{
			name: "正例-登录成功",
			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
				userSvc := svcmocks.NewMockUserServiceInterface(ctrl)
				userSvc.EXPECT().Login(gomock.Any(), gomock.Any()).Return(domain.UserDomain{
					Id:    1,
					Email: "112@qq.com",
				}, nil)
				return userSvc
			},
			body: `
{
    "email":"112@qq.com",
    "password":"hello#word123"
}`,
			wantCode: 200,
			wantBody: domain.Result{
				Code: 0,
				Msg:  "Login successful",
				Data: nil,
			},
		},
		{
			name: "负例-参数不对，bind 失败",
			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
				userSvc := svcmocks.NewMockUserServiceInterface(ctrl)
				//负例还没有调用到signup方法
				return userSvc
			},
			body: `
{
    "ema2323il":"112@qq.com",
    "password":"hello#word123",
}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "负例-账号密码不匹配",
			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
				userSvc := svcmocks.NewMockUserServiceInterface(ctrl)
				userSvc.EXPECT().Login(gomock.Any(), gomock.Any()).Return(domain.UserDomain{}, service.ErrInvalidUserOrPassword)
				return userSvc
			},
			body: `
{
    "email":"112@qq.com",
    "password":"hello#word123"
}`,
			wantCode: 200,
			wantBody: domain.Result{
				Code: 4,
				Msg:  "The account or password is incorrect",
				Data: nil,
			},
		},
		{
			name: "负例-系统错误",
			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
				userSvc := svcmocks.NewMockUserServiceInterface(ctrl)
				userSvc.EXPECT().Login(gomock.Any(), gomock.Any()).Return(domain.UserDomain{}, errors.New("系统异常"))
				return userSvc
			},
			body: `
{
    "email":"112@qq.com",
    "password":"hello#word123"
}`,
			wantCode: http.StatusInternalServerError,
			wantBody: domain.Result{
				Code: 5,
				Msg:  "System error",
				Data: nil,
			},
		},
		{
			name: "负例-jwtToken设置错误",
			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
				userSvc := svcmocks.NewMockUserServiceInterface(ctrl)
				userSvc.EXPECT().Login(gomock.Any(), gomock.Any()).Return(domain.UserDomain{
					Id:    1,
					Email: "112@qq.com",
				}, nil)
				return userSvc
			},
			body: `
{
    "email":"112@qq.com",
    "password":"hello#word123"
}`,
			wantCode: 200,
			wantBody: domain.Result{
				Code: 5,
				Msg:  "System err",
				Data: nil,
			},
		},
	}
	for _, tt := range testsCase {
		t.Run(tt.name, func(t *testing.T) {
			server := gin.Default()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			handler := NewUserHandler(tt.mock(ctrl), nil)
			handler.RegisterRoutesCt(server)
			req, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer([]byte(tt.body)))
			require.NoError(t, err)
			//数据是json格式
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req) //这是http路由，请求进去GIN框架的入口
			assert.Equal(t, tt.wantCode, resp.Code)
			var bodyStr domain.Result
			err = json.NewDecoder(resp.Body).Decode(&bodyStr)
			assert.Equal(t, tt.wantBody, bodyStr)

		})
	}
}
