package web

import (
	"GeekProject/newGeekProject/day2/webook/internal/domain"
	"GeekProject/newGeekProject/day2/webook/internal/service"
	"GeekProject/newGeekProject/day2/webook/internal/service/oauth2/wechat"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
	"time"
)

type OAuth2WechatHandler struct {
	wechatS wechat.ServiceWechatInterface
	jwtHandler
	userService service.UserServiceInterface
	stateKey    []byte
}

func NewOAuth2WechatHandler(wechatS wechat.ServiceWechatInterface, userService service.UserServiceInterface) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		wechatS:     wechatS,
		userService: userService,
		stateKey:    []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"),
	}
}

func (wh *OAuth2WechatHandler) RegisterRouter(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authUrl", wh.AuthURL)
	g.Any("/callback", wh.Callback) //回调
}

func (wh *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	state := uuid.New()
	url, err := wh.wechatS.AuthUrl(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 5,
			Msg:  "构造扫码登录错误",
			Data: nil,
		})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		State: state,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	})
	signedString, err := token.SignedString(wh.stateKey)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 5,
			Msg:  "系统错误",
			Data: nil,
		})
		return
	}
	ctx.SetCookie("jwt-state", signedString, 600, "/oauth2/wechat", "",
		false, true)

	ctx.JSON(http.StatusOK, domain.Result{
		Code: 0,
		Msg:  "成功",
		Data: nil,
	})
	fmt.Printf("url:[%v][%v]\n", url, signedString)

}

func (wh *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	err := wh.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 5,
			Msg:  "登录失败",
			Data: nil,
		})
		return
	}
	info, err := wh.wechatS.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 5,
			Msg:  "系统错误",
			Data: nil,
		})
		return
	}
	//从userService里面那uid
	createWechat, err := wh.userService.FindByCreateWechat(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 5,
			Msg:  "系统错误",
			Data: nil,
		})
		return
	}
	err = wh.SetJWTToken(ctx, createWechat.Id)
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

func (wh *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	cookie, err := ctx.Cookie("jwt-state")
	if err != nil {
		return fmt.Errorf("拿不到state的cookie[%v]", err)
	}
	var sc StateClaims
	claims, err := jwt.ParseWithClaims(cookie, &sc, func(token *jwt.Token) (interface{}, error) {
		return wh.stateKey, nil
	})
	if err != nil || claims.Valid {
		return fmt.Errorf("token过期[%v]", err)
	}
	if sc.State != state {
		return errors.New("state不相等")
	}
	return err
}

type StateClaims struct {
	State string
	jwt.RegisteredClaims
}
