package web

import (
	"GeekProject/newGeekProject/day2/webook/internal/domain"
	"GeekProject/newGeekProject/day2/webook/internal/service/oauth2/wechat"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OAuth2WechatHandler struct {
	wechatS wechat.ServiceWechatInterface
}

func NewOAuth2WechatHandler(wechatS wechat.ServiceWechatInterface) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		wechatS: wechatS,
	}
}

func (wh *OAuth2WechatHandler) RegisterRouter(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authUrl", wh.AuthURL)
	g.Any("/callback", wh.Callback) //回调
}

func (wh *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	url, err := wh.wechatS.AuthUrl(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 5,
			Msg:  "构造扫码登录错误",
			Data: nil,
		})
		return
	}
	fmt.Printf("url:[%v]\n", url)

}

func (wh *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")
	info, err := wh.wechatS.VerifyCode(ctx, code, state)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 5,
			Msg:  "系统错误",
			Data: nil,
		})
		return
	}
	fmt.Printf("info[%v]", info)
	ctx.String(http.StatusOK, "ok")
}
