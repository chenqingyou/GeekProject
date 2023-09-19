//go:build wireinject

package intergration

import (
	"GeekProject/newGeekProject/day2/webook/internal/repository"
	"GeekProject/newGeekProject/day2/webook/internal/repository/cache"
	"GeekProject/newGeekProject/day2/webook/internal/repository/dao"
	"GeekProject/newGeekProject/day2/webook/internal/service"
	"GeekProject/newGeekProject/day2/webook/internal/web"
	"GeekProject/newGeekProject/day2/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		//最基础的第三方依赖
		ioc.InitDB, ioc.InitRedis,
		dao.NewUserDao,
		cache.NewUserCache, cache.NewCodeCache,
		repository.NewUserRepository, repository.NewCodeRepository,
		service.NewUserService, service.NewCodeService,
		web.NewUserHandler,
		ioc.InitSMSService,
		ioc.InitGin,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}
