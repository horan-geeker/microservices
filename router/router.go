package router

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"microservices/entity/config"
	"microservices/pkg/app"
	"microservices/router/middleware"
)

// 初始化资源
var (
	env    = config.GetEnvConfig()
	router = app.NewApp(app.NewServerOptions(env.AppEnv, env.ServerHost, env.ServerPort, env.ServerTimeout))
	_      = router.Use(gzip.Gzip(gzip.DefaultCompression), middleware.RequestLogger, middleware.ReqRateLimit(100, 1, false, false), gin.Recovery())
)
