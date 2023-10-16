package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func main() {
	if err := initViperV2(); err != nil {
		panic(err)
	}
	server := InitWebServer()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "你好 你来了")
	})
	server.Run(":8080")
}

func initViperV2() error {
	viper.SetConfigFile("cofig/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return err
}

func initViper() error {
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

//func initUser(serverDB *gorm.DB, client redis.Cmdable) *web.UserHandler {
//	userDao := dao.NewUserDao(serverDB)
//	userCache := cache.NewUserCache(client)
//	userRepo := repository.NewUserRepository(userDao, userCache)
//	userSvc := service.NewUserService(userRepo)
//	//验证码缓存
//	codeCache := cache.NewCodeCache(client)
//	//初始化验证码服务
//	codeRepo := repository.NewCodeRepository(codeCache)
//	//
//	codeSms := memory.NewService()
//	codeSvc := service.NewCodeService(codeRepo, codeSms)
//	user := web.NewUserHandler(userSvc, codeSvc)
//	return user
//}
