package router

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"microservices/internal/config"
	"microservices/internal/middleware"
	"microservices/internal/pkg/options"
	"microservices/internal/service"
	"microservices/internal/store/mysql"
	"microservices/internal/store/redis"
	"microservices/pkg/app"
)

var (
	// 注入全局中间件，注意压缩 response 的中间件顺序需要在 log response 之后压缩
	env    = config.GetConfig()
	router = app.NewApp(app.NewServerOptions(env.AppEnv, env.ServerHost, env.ServerPort, env.ServerTimeout))
	_      = router.Use(gzip.Gzip(gzip.DefaultCompression), middleware.RequestLogger, gin.Recovery())

	dataFactory    = mysql.GetMysqlInstance(options.NewMysqlOptions())
	cacheFactory   = redis.GetRedisInstance(options.NewRedisOptions())
	serviceFactory = service.GetServiceInstance(options.NewTencentOptions(), options.NewAliyunOptions())
)
