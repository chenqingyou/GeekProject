package main

import (
	"GeekProject/homeWork/class2/webook/internal/repository"
	"GeekProject/homeWork/class2/webook/internal/repository/cache"
	"GeekProject/homeWork/class2/webook/internal/repository/dao"
	"GeekProject/homeWork/class2/webook/internal/service"
	"GeekProject/homeWork/class2/webook/internal/service/sms/memory"
	"GeekProject/homeWork/class2/webook/internal/web"
	"GeekProject/homeWork/class2/webook/internal/web/middleware"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	if err := initViperConfig(); err != nil {
		panic(err)
	}

	serverDB, err := initDB()
	if err != nil {
		fmt.Printf("Open DB init err [%v]\n", err)
		return
	}

	server := initServer()
	user := initUser(serverDB)
	//根据使用习惯
	//方式1：传入分组
	//user.RegisterRoutesV1(server.Group("/users"))
	//方式2:定义好分组
	user.RegisterRoutesCt(server)
	server.Run(":8080")
}

func initServer() *gin.Engine {
	server := gin.Default()
	//注册接口
	//解决跨域问题
	/*
			• AllowCrendentials：是否允许带上用户认证 信息（比如 cookie）。
			• AllowHeader：业务请求中可以带上的头。
			• AllowOriginFunc：哪些来源是允许的。
		• 跨域问题是因为发请求的域名+端口和接收请求的域名 + 端口对不上。比如说这里的 localhost:3000 发到 localhost:8080 上。
		• 解决跨域问题的关键是在 preflight 请求里面告诉浏览器自己愿意接收请求。
		• Gin 提供了解决跨域问题的 middleware，可以直接使用。
		• middleware 是一种机制，可以用来解决一些所有业务都关心的问题，使用 Use 方法来注册 middleware。
	*/
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
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("mysession", store))
	server.Use(middleware.NewLoginMiddlewareBuilder().DepositPaths("/users/signup").
		DepositPaths("/users/login").DepositPaths("/users/loginSms/code").BuildSess())

	return server
}

func initUser(serverDB *gorm.DB) *web.UserHandler {
	userDao := dao.NewUserDao(serverDB)
	localCache := cache.NewUserLocalCache()
	localRepository := repository.NewCodeRepository(localCache)
	userRepo := repository.NewUserRepository(userDao, localCache)
	userSvc := service.NewUserService(userRepo)
	sms := memory.NewSmSService()
	cacheSvc := service.NewCodeService(localRepository, sms)
	user := web.NewUserHandler(userSvc, cacheSvc)
	return user
}

func initDB() (db *gorm.DB, err error) {
	//初始化数据库
	type DBConfig struct {
		Dns string `json:"dns"`
	}
	var cfg = DBConfig{Dns: "root:root@tcp(localhost:13316)/webook_debug"} //设置默认值
	err = viper.UnmarshalKey("db.mysql.dsn", &cfg)
	if err != nil {
		panic(err)
	}
	serverDB, err := gorm.Open(mysql.Open(cfg.Dns))
	if err != nil {
		fmt.Printf("Open DB init err [%v]\n", err)
		return nil, err
	}
	err = dao.InitDBTable(serverDB) //创建表
	if err != nil {
		fmt.Printf("DB InitDBTable err [%v]\n", err)
		return nil, err
	}
	return serverDB, nil
}

func initViperConfig() error {
	//配置文件的名字
	viper.SetConfigName("dev")
	//配置文件的类型
	viper.SetConfigType("yaml")
	//配置文件目录
	viper.AddConfigPath("./config")
	//viper.AddConfigPath("./temp/config")
	//读取配置到viper里面，或者加载到内存里面
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return err
}
