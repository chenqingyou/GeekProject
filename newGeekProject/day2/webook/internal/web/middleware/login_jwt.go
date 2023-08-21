package middleware

import (
	"GeekProject/newGeekProject/day2/webook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}
func (lmb *LoginJWTMiddlewareBuilder) DepositPaths(path string) *LoginJWTMiddlewareBuilder {
	lmb.paths = append(lmb.paths, path)
	return lmb
}

func (lmb *LoginJWTMiddlewareBuilder) BuildSess() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range lmb.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		//现在用jwt来校验
		tokenStr := ctx.GetHeader("Authorization")
		if tokenStr == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//segs := strings.Split(tokenStr, " ")
		//if len(segs) != 2 {
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		//tokenStr := segs[1]
		claims := &web.TokenClaims{}
		//ParseWithClaims 一定要使用指针， 因为ParseWithClaims会去修改claims里面的值
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !token.Valid || claims.Uid == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//登录安全校验
		if claims.UserAgent != ctx.Request.UserAgent() {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//保证token的时效性
		now := time.Now()
		// 每十秒钟刷新一次
		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
			if err != nil {
				// 记录日志
				log.Println("jwt 续约失败", err)
			}
			ctx.Header("x-jwt-token", tokenStr)
		}
		//可以将设置的值传入cxt里面
		ctx.Set("claims", claims)
	}
}
