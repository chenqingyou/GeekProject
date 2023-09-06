package ioc

import (
	"GeekProject/newGeekProject/day2/webook/internal/web"
	"GeekProject/newGeekProject/day2/webook/internal/web/middleware"
	"GeekProject/newGeekProject/day2/webook/pkg/ginx/middlewares/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitGin(mdls []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutesCt(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		//使用限流的中间件
		ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
		/*
				• AllowCrendentials：是否允许带上用户认证 信息（比如 cookie）。
				• AllowHeader：业务请求中可以带上的头。
				• AllowOriginFunc：哪些来源是允许的。
			• 跨域问题是因为发请求的域名+端口和接收请求的域名 + 端口对不上。比如说这里的 localhost:3000 发到 localhost:8080 上。
			• 解决跨域问题的关键是在 preflight 请求里面告诉浏览器自己愿意接收请求。
			• Gin 提供了解决跨域问题的 middleware，可以直接使用。
			• middleware 是一种机制，可以用来解决一些所有业务都关心的问题，使用 Use 方法来注册 middleware。
		*/
		cors.New(cors.Config{ //Use作用于全部路由
			//AllowOrigins: []string{"http://localhost:3000"},
			//AllowMethods:     []string{"PUT", "PATCH", "POST"},
			AllowHeaders: []string{"Authorization", "Content-Type"},
			//token获取
			ExposeHeaders: []string{"x-jwt-token"},
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
		}),
		//存放sess
		sessions.Sessions("mysession", memstore.NewStore([]byte("WZmKWNA1rxZf9TCoRBNsNDIlKdHb6DrzwK2NFF9n7a8ueRfinsAWFqVskMalYtgo"),
			[]byte("lUzWwJAb6zaC1C5lRELHwDRNiYRwIC3nhL80dzBffEy7EsRGnzuOSa8BkqooCZ6W"))),
		middleware.NewLoginJWTMiddlewareBuilder().DepositPaths(
			"/users/signup").DepositPaths("/users/loginJwt").DepositPaths("/users/loginSms/code").DepositPaths("/users/loginSms").BuildSess(),
	}
}

func InitServer() *gin.Engine {
	server := gin.Default()
	//
	//redisClient := redis.NewClient(&redis.Options{
	//	Addr: config.Config.Redis.Addr,
	//})
	//server.Use()
	//注册接口
	//
	//使用cookie
	//store := cookie.NewStore([]byte("secret"))
	//使用内存存放
	//使用redis存放--多实例适合
	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
	//	[]byte("nCbJB9lqD8O1WSVnCaSZweak"),
	//	[]byte("nCbJB9lqD8O1WSVnCaSZweak"))
	//if err != nil {
	//	panic(err)
	//}
	server.Use()
	//size 16 代表最大空闲链接
	//使用自己定义的store
	//mystore := &sqlx_store.Store{}
	//server.Use(sessions.Sessions("mysession", mystore))

	//server.Use(middleware.NewLoginMiddlewareBuilder().DepositPaths(
	//	"/users/signup").DepositPaths("/users/login").BuildSess())

	//使用jwt模式登录
	server.Use()

	return server
}
