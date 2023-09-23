package main

import (
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"microservices/internal/config"
	_ "microservices/internal/config"
	"microservices/internal/middleware"
	_ "microservices/internal/router"
	"microservices/pkg/meta"
	"os"
)

var app = meta.GetEnginInstance()

func main() {
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000000",
	})
	log.SetOutput(os.Stdout)
	env := config.NewEnvConfig()
	// disable gin log
	if env.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	gin.DefaultWriter = io.Discard
	app.Engine.Use(gzip.Gzip(gzip.DefaultCompression)) // 注意顺序需要在 log response 之后压缩
	app.Engine.Use(middleware.RequestLogger)
	app.Engine.Use(gin.Recovery())
	log.Info("server run ", env.ServerHost, ":", env.ServerPort)
	if err := app.Run(fmt.Sprintf("%s:%d", env.ServerHost, env.ServerPort)); err != nil {
		panic(err)
	}
}
