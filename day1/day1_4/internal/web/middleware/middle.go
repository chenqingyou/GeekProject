package middleware

import (
	"GeekProject/day1/day1_4/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type Middleware struct {
	Paths []string
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

// PathAdd 不需要sess的路由
func (ml *Middleware) PathAdd(path string) *Middleware {
	ml.Paths = append(ml.Paths, path)
	return ml
}

// CrossDomain 跨域问题解决
func (ml *Middleware) CrossDomain(server *gin.Engine) {
	server.Use(cors.New(cors.Config{ //Use作用于全部路由
		//AllowOrigins: []string{"http://localhost:3000"},
		//AllowMethods:     []string{"PUT", "PATCH", "POST"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
		//是否允许你带cookie之类的东西
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				//你的开发环境
				return true
			}
			return strings.Contains(origin, "yunming.com")
		},
		MaxAge: 12 * time.Hour,
	}))
}

// Sess  问题
func (ml *Middleware) Sess(server *gin.Engine) {
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("mysession", store))
	server.Use(func(cxt *gin.Context) {
		for _, path := range ml.Paths {
			if cxt.Request.URL.Path == path {
				cxt.Next()
				return
			}
		}
		// 获取当前请求的路径
		// 在这里可以根据需要添加更复杂的逻辑来判断是否需要Session验证
		// 这里只是一个示例，只有路径为"/private"的接口需要Session验证
		session := sessions.Default(cxt)
		userId := session.Get(web.UserIdKey)
		// 如果Session中没有合法的用户名，则返回未授权
		if userId == nil {
			cxt.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			cxt.Abort()
			return
		}
		cxt.Next()
	})
}
