package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	server := InitWebServer()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "你好 你来了")
	})
	server.Run(":8080")
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
