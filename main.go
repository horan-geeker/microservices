package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"microservices/internal/config"
	"microservices/internal/pkg/options"
	_ "microservices/internal/router"
	"microservices/internal/service"
	"microservices/internal/store/mysql"
	"microservices/internal/store/redis"
	"microservices/pkg/app"
	"os"
)

// 初始化资源
var (
	env = config.NewConfig()
	_   = mysql.GetMysqlInstance(options.NewMysqlOptions())
	_   = redis.GetRedisInstance(options.NewRedisOptions())
	_   = service.GetServiceInstance(options.NewTencentOptions(), options.NewAliyunOptions())
)

func main() {
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000000",
	})
	log.SetOutput(os.Stdout)
	// disable gin log
	if env.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	gin.DefaultWriter = io.Discard
	server := app.NewApp(app.NewServerOptions(3))
	log.Info("server run ", env.ServerHost, ":", env.ServerPort)
	if err := server.Run(fmt.Sprintf("%s:%d", env.ServerHost, env.ServerPort)); err != nil {
		panic(err)
	}
}
