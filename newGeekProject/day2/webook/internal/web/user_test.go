package web

import (
	"GeekProject/newGeekProject/day2/webook/internal/domain"
	"GeekProject/newGeekProject/day2/webook/internal/service"
	svcmocks "GeekProject/newGeekProject/day2/webook/internal/service/mocks"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

/*mockgen -source=./webook/internal/service/user.go -package=svcmocks - destination=./webook/internal/service/mocks/user.mock.go
 */
func TestUserHandler_SignUp(t *testing.T) {
	testCase := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserServiceInterface
		body     string
		wantCode int
		wantBody string
	}{
		{
			name: "正例-注册成功",
			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
				userSvc := svcmocks.NewMockUserServiceInterface(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)
				//codeSvc := svcmocks.NewMockCodeServiceInterface(ctrl)
				return userSvc
			},
			body: `
{
    "email":"112@qq.com",
    "password":"hello#word123",
    "confirmPassword":"hello#word123"
}`,
			wantCode: 200,
			wantBody: "Registered successfully",
		},
		{
			name: "负例-参数不对bind失败",
			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
				userSvc := svcmocks.NewMockUserServiceInterface(ctrl)
				//负例还没有调用到signup方法
				return userSvc
			},
			body: `
{
    "email123":"112@qq.com",
    "password123":"hello#word123",
    "confirmPassword":"hello#word123",
}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "负例-邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserServiceInterface {
				userSvc := svcmocks.NewMockUserServiceInterface(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(service.ErrUserDuplicateEmail)
				//codeSvc := svcmocks.NewMockCodeServiceInterface(ctrl)
				return userSvc
			},
			body: `
{
    "email":"112@qq.com",
    "password":"hello#word123",
    "confirmPassword":"hello#word123"
}`,
			wantCode: 200,
			wantBody: "Mailbox conflict",
		},
	}
	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			server := gin.Default()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			handler := NewUserHandler(tt.mock(ctrl), nil, nil)
			handler.RegisterRoutesCt(server)
			req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(tt.body)))
			require.NoError(t, err)
			//数据是json格式
			req.Header.Set("Content-Type", "application/json")
			resq := httptest.NewRecorder()
			fmt.Println(resq)
			//这儿会遇到初始化第三方的问题
			server.ServeHTTP(resq, req) //这是http路由，请求进去GIN框架的入口
			assert.Equal(t, tt.wantCode, resq.Code)
			assert.Equal(t, tt.wantBody, resq.Body.String())
		})
	}
}

func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userSvc := svcmocks.NewMockUserServiceInterface(ctrl)
	userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("mock err"))
	err := userSvc.SignUp(context.Background(), domain.UserDomain{
		Id:              0,
		Email:           "",
		Password:        "",
		Nickname:        "",
		Phone:           "",
		Birthday:        "",
		PersonalProfile: "",
		Ctime:           time.Time{},
	})
	if err != nil {
		t.Log(err)
	}
}
