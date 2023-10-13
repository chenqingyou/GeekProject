package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

type jwtHandler struct {
	//access_token key
	atKey []byte
	//refresh_token key
	rtKey []byte
}

func NewJwtHandler() jwtHandler {
	return jwtHandler{
		atKey: []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"),
		rtKey: []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"),
	}
}

// 长短token设置
type RefreshClaims struct {
	Uid int64
	jwt.RegisteredClaims
}

// TokenClaims 实现jwt的接口
type TokenClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}

func (j jwtHandler) SetJWTToken(ctx *gin.Context, uId int64) error {
	//使用 ecdsa.GenerateKey 生成了一个 ECDSA 密钥对，并将私钥用于 JWT 的签名
	//privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	//if err != nil {
	//	return err
	//}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)), //过期时间
			NotBefore: nil,
		},
		Uid:       uId,
		UserAgent: ctx.Request.UserAgent(),
	})
	signedString, err := token.SignedString(j.rtKey)
	if err != nil {
		return err
	}
	fmt.Printf("signedString [%v]\n", signedString)
	ctx.Header("x-jwt-token", signedString)
	return nil
}

func (j jwtHandler) setRefreshToken(ctx *gin.Context, uId int64) error {
	//使用 ecdsa.GenerateKey 生成了一个 ECDSA 密钥对，并将私钥用于 JWT 的签名
	//privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	//if err != nil {
	//	return err
	//}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)), //过期时间
			NotBefore: nil,
		},
		Uid: uId,
	})
	signedString, err := token.SignedString(j.rtKey)
	if err != nil {
		return err
	}
	fmt.Printf("signedString [%v]\n", signedString)
	ctx.Header("x-refresh-token", signedString)
	return nil
}

func ExtractToken(ctx *gin.Context) string {
	tokenStr := ctx.GetHeader("Authorization")
	if tokenStr == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return ""
	}
	segs := strings.Split(tokenStr, " ")
	return segs[1]

}
