package intergration

import (
	"GeekProject/newGeekProject/day2/webook/internal/domain"
	"GeekProject/newGeekProject/day2/webook/ioc"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

//集成测试。发送验证码

func TestUserHandler_e2e_SendLoginSMSCode(t *testing.T) {
	server := InitWebServer()
	rdb := ioc.InitRedis()
	testCase := []struct {
		name string
		//准备数据，已经验证数据
		before func(t *testing.T)
		//验证数据
		after    func(t *testing.T)
		body     string
		wantCode int
		wantBody domain.Result
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {
				//不需要，redis里面什么数据都没有
			},
			after: func(t *testing.T) {
				//数据清理
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
				val, err := rdb.GetDel(ctx, "phone_code:login:12312312").Result()
				cancel()
				assert.NoError(t, err)
				//验证验证码位数
				assert.True(t, len(val) == 6)
			},
			body: `{
	"phone":"12312312"
}`,
			wantCode: 200,
			wantBody: domain.Result{
				Code: 0,
				Msg:  "验证码发送成功",
			},
		},
		{
			name: "发送频繁",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
				_, err := rdb.Set(ctx, "phone_code:login:12312312", "123456", time.Minute*10).Result()
				cancel()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				//数据清理
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
				val, err := rdb.GetDel(ctx, "phone_code:login:12312312").Result()
				cancel()
				assert.NoError(t, err)
				//验证验证码位数
				assert.Equal(t, "123456", val)
			},
			body: `{
	"phone":"12312312"
}`,
			wantCode: 200,
			wantBody: domain.Result{
				Code: 5,
				Msg:  "验证码发送太频繁",
			},
		},
		{
			name: "系统错误", //只有一个条件，就是设置的验证码没有过期时间
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
				_, err := rdb.Set(ctx, "phone_code:login:12312312", "123456", 0).Result()
				cancel()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				//数据清理
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
				val, err := rdb.GetDel(ctx, "phone_code:login:12312312").Result()
				cancel()
				assert.NoError(t, err)
				//验证验证码位数
				assert.Equal(t, "123456", val)
			},
			body: `{
	"phone":"12312312"
}`,
			wantCode: 200,
			wantBody: domain.Result{
				Code: 5,
				Msg:  "验证码发送失败",
			},
		},
		{
			name:   "手机号为空", //只有一个条件，就是设置的验证码没有过期时间
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
			},
			body: `{
	"phone":""
}`,
			wantCode: 200,
			wantBody: domain.Result{
				Code: 5,
				Msg:  "输入有误",
			},
		},
		{
			name:   "输入数据有误", //只有一个条件，就是设置的验证码没有过期时间
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
			},
			body: `{
	"phone":
}`,
			wantCode: 400,
		},
	}
	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			tt.before(t)
			req, err := http.NewRequest(http.MethodPost, "/users/loginSms/code", bytes.NewBuffer([]byte(tt.body)))
			require.NoError(t, err)
			//数据是json格式
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			// 这就是 HTTP 请求进去 GIN 框架的入口。
			// 当你这样调用的时候，GIN 就会处理这个请求
			// 响应写回到 resp 里
			fmt.Println(resp)
			//这儿会遇到初始化第三方的问题
			server.ServeHTTP(resp, req) //这是http路由，请求进去GIN框架的入口
			assert.Equal(t, tt.wantCode, resp.Code)
			if resp.Code != 200 {
				return
			}
			var webRes domain.Result
			//用json.Unmarshal会全部读取一遍Body，NewDecoder不需要读取直接用
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tt.wantBody, webRes)
			tt.after(t)
		})
	}
}
