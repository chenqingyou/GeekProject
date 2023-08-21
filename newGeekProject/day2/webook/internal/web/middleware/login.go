package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}
func (lmb *LoginMiddlewareBuilder) DepositPaths(path string) *LoginMiddlewareBuilder {
	lmb.paths = append(lmb.paths, path)
	return lmb
}

func (lmb *LoginMiddlewareBuilder) BuildSess() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range lmb.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		//使用session
		sess := sessions.Default(ctx)
		id := sess.Get("userID")
		if id == nil {
			//没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		sess.Set("userID", id)
		sess.Options(sessions.Options{
			MaxAge: 1 * 60,
		})
		updateTime := sess.Get("update_time")
		nowTime := time.Now().UnixMilli()
		//第一次登录
		if updateTime == nil {
			sess.Set("update_time", nowTime)
			sess.Save()
			return
		}
		updateTimeVail, ok := updateTime.(int64)
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		if nowTime-updateTimeVail > 10*1000 {
			sess.Set("update_time", nowTime)
			sess.Save()
			return
		}
	}
}
