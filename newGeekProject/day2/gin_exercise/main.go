package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	server := gin.Default()
	//get请求，就执行这段代码
	server.GET("/hello", func(ctx *gin.Context) {
		//返回值
		ctx.JSON(http.StatusOK, "hello,go")
	})
	//go func() {  可以监听多个端口
	//	server1 := gin.Default()
	//	server1.Run(":8080")
	//}()

	//静态路由，完全匹配
	server.POST("/post", func(ctx *gin.Context) { //所有的都是两个参数，一个是路由，一个是方法
		ctx.JSON(http.StatusOK, "hello go")
	})
	//参数路由，不能单独使用:/user/:
	server.GET("/user/:name", func(ctx *gin.Context) { //所有的都是两个参数，一个是路由，一个是方法
		name := ctx.Param("name") //获取传进来的参数
		ctx.JSON(http.StatusOK, "hello,这是参数路由"+name)
	})
	//通配符路由，*不能单独出现，比如/views/*
	server.GET("/views/*.html", func(ctx *gin.Context) { //所有的都是两个参数，一个是路由，一个是方法
		ctx.JSON(http.StatusOK, "hello,这是通配符路由")
	})

	//查询参数
	server.GET("/order", func(ctx *gin.Context) {
		oid := ctx.Query("id")
		ctx.JSON(http.StatusOK, "hello,这是查询参数"+oid)
	})

	//监听端口默认为8080
	err := server.Run(":8080")
	if err != nil {
		return
	}
}
