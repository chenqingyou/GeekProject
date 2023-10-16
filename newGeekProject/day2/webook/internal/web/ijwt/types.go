package ijwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type HandlerJWTInterface interface {
	ExtractToken(ctx *gin.Context) string
	ClearToken(ctx *gin.Context) error
	CheckSession(ctx *gin.Context, ssid string) error
	SetRefreshToken(ctx *gin.Context, uId int64, ssid string) error
	SetJWTToken(ctx *gin.Context, uId int64, ssid string) error
	SetLoginToken(ctx *gin.Context, uId int64) error
}

// 长短token设置
type RefreshClaims struct {
	Uid  int64
	Ssid string
	jwt.RegisteredClaims
}

// TokenClaims 实现jwt的接口
type TokenClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	Ssid      string
	UserAgent string
}

type jwtHandler struct {
}
