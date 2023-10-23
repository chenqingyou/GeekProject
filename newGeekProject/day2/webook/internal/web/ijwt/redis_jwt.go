package ijwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strings"
	"time"
)

var (
	//access_token key
	AtKey = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0")
	//refresh_token key
	RtKey = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0")
)

type RedisJwtHandler struct {
	cmd redis.Cmdable
}

func NewRedisJwtHandler(cmd redis.Cmdable) HandlerJWTInterface {
	return &RedisJwtHandler{
		cmd: cmd,
	}
}

func (j *RedisJwtHandler) SetLoginToken(ctx *gin.Context, uId int64) error {
	//只有这个接口，拿出来的是refresh_token，其他的都是短token
	ssid := uuid.New().String()
	err := j.SetJWTToken(ctx, uId, ssid)
	if err != nil {
		return err
	}
	err = j.SetRefreshToken(ctx, uId, ssid)
	if err != nil {
		return err
	}
	return nil
}

func (j *RedisJwtHandler) ClearToken(ctx *gin.Context) error {
	//需要将长短token都设置成一个非法的值
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")
	claims := ctx.MustGet("claims").(*TokenClaims)
	err := j.cmd.Set(ctx, fmt.Sprintf("users:ssid:%s", claims.Ssid), "",
		time.Hour*24*7).Err()
	return err
}

func (j *RedisJwtHandler) CheckSession(ctx *gin.Context, ssid string) error {
	val, err := j.cmd.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	switch err {
	case redis.Nil:
		return nil
	case nil:
		if val == 0 {
			return nil
		}
		return errors.New("session 已经失效了")
	default:
		return err
	}
}

func (j *RedisJwtHandler) SetRefreshToken(ctx *gin.Context, uId int64, ssid string) error {
	//使用 ecdsa.GenerateKey 生成了一个 ECDSA 密钥对，并将私钥用于 JWT 的签名
	//privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	//if err != nil {
	//	return err
	//}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), //过期时间
			NotBefore: nil,
		},
		Uid:  uId,
		Ssid: ssid,
	})
	signedString, err := token.SignedString(RtKey)
	if err != nil {
		return err
	}
	fmt.Printf("signedString [%v]\n", signedString)
	ctx.Header("x-refresh-token", signedString)
	return nil
}

func (j *RedisJwtHandler) ExtractToken(ctx *gin.Context) string {
	tokenStr := ctx.GetHeader("Authorization")
	if tokenStr == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return ""
	}
	segs := strings.Split(tokenStr, " ")
	return segs[1]

}

func (j *RedisJwtHandler) SetJWTToken(ctx *gin.Context, uId int64, ssid string) error {
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
		Ssid:      ssid,
		UserAgent: ctx.Request.UserAgent(),
	})
	signedString, err := token.SignedString(AtKey)
	if err != nil {
		return err
	}
	fmt.Printf("signedString [%v]\n", signedString)
	ctx.Header("x-jwt-token", signedString)
	return nil
}
