package main

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"net/http"
)

func main() {
	if err := initViper(); err != nil {
		panic(err)
	}
	server := InitWebServer()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "你好 你来了")
	})
	server.Run(":8080")
}

// 配置中心远程加载
func initViperRemote() {
	viper.SetConfigType("yaml")
	//remote不支持key的切割，比如db.mysql
	err := viper.AddRemoteProvider("etcd3", "http://127.0.0.1:12378", "/webook")
	if err != nil {
		panic(err)
	}
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

func initViperReader() {
	viper.SetConfigType("yaml")
	cfg := `
db.mysql:
  dsn: "root:root@tcp(localhost:13316)/webook"
redis:
  addr: "localhost:6379"`
	err := viper.ReadConfig(bytes.NewReader([]byte(cfg)))
	if err != nil {
		panic(err)
	}

}

func initViperV2() error {
	//设置默认值
	viper.SetDefault("db.mysql.dsn", "root:root@tcp(localhost:13316)/webook")
	viper.SetConfigFile("cofig/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return err
}

func initViper() error {
	//使用启动命令
	cfile := pflag.String("config", "config/config.yaml", "指定配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	viper.WatchConfig()                            //实时监听配置文件变更
	viper.OnConfigChange(func(in fsnotify.Event) { //文件通知，只有本地，远程配置中心无法使用
		fmt.Println(in.Name, in.Op) //这里代表配置文件发生变化，但是不能告诉你文件的那些内容变了
	})
	//配置文件的名字
	//viper.SetConfigName("dev")
	////配置文件的类型
	//viper.SetConfigType("yaml")
	////配置文件目录
	//viper.AddConfigPath("./config")
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
