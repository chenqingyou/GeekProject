package web

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type jwtHandler struct {
}

// TokenClaims 实现jwt的接口
type TokenClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}

func (j jwtHandler) SetJWTToken(ctx *gin.Context, uId int64) error {
	//使用 ecdsa.GenerateKey 生成了一个 ECDSA 密钥对，并将私钥用于 JWT 的签名
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES512, TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)), //过期时间
			NotBefore: nil,
		},
		Uid:       uId,
		UserAgent: ctx.Request.UserAgent(),
	})
	signedString, err := token.SignedString(privateKey)
	if err != nil {
		return err
	}
	fmt.Printf("signedString [%v]\n", signedString)
	ctx.Header("x-jwt-token", signedString)
	return nil
}
