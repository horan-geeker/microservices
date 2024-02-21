package router

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"microservices/entity/config"
	"microservices/pkg/app"
	"microservices/repository"
	"microservices/router/middleware"
)

var (
	// 注入全局中间件，注意压缩 response 的中间件顺序需要在 log response 之后压缩
	env    = config.GetConfig()
	router = app.NewApp(app.NewServerOptions(env.AppEnv, env.ServerHost, env.ServerPort, env.ServerTimeout))
	_      = router.Use(gzip.Gzip(gzip.DefaultCompression), middleware.RequestLogger, gin.Recovery())

	repositoryFactory = repository.NewFactory()
)
