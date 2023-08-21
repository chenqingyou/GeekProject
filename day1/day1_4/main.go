package main

import (
	"GeekProject/day1/day1_4/internal/repository"
	"GeekProject/day1/day1_4/internal/repository/dao"
	"GeekProject/day1/day1_4/internal/service"
	"GeekProject/day1/day1_4/internal/web"
	"GeekProject/day1/day1_4/internal/web/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	serverDB, err := initDB("root:root@tcp(localhost:13316)/webook")
	if err != nil {
		fmt.Printf("Open DB init err [%v]\n", err)
		return
	}
	server := initServer()
	//初始化数据库以业务
	webserver := initWebServer(serverDB)
	//初始化路由
	webserver.RegisterRoutes(server)
	server.Run(":8080")
}

func initServer() *gin.Engine {
	//启动gin框架
	server := gin.Default()
	//跨域问题
	middleware.NewMiddleware().CrossDomain(server)
	//设置session
	middleware.NewMiddleware().PathAdd("/users/signup").PathAdd("/users/login").Sess(server)
	return server
}

func initWebServer(db *gorm.DB) *web.UserWebHandler {
	ud := dao.NewUserDao(db)
	ur := repository.NewUserRep(ud)
	us := service.NewUserService(ur)
	return web.NewUserWebHandler(us)
}

func initDB(mysqlStr string) (db *gorm.DB, err error) {
	//初始化数据库
	serverDB, err := gorm.Open(mysql.Open(mysqlStr))
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
