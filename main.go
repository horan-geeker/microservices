package main

import (
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"microservices/internal/config"
	"microservices/internal/middleware"
	"microservices/internal/router"
	"os"
	"time"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000000",
	})
	log.SetOutput(os.Stdout)
	// disable gin log
	if config.AppConfig.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	gin.DefaultWriter = io.Discard
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression)) // 注意顺序需要在 log response 之前压缩
	r.Use(middleware.RequestLogger)
	r.Use(gin.Recovery())
	router.Register(r)
	log.Println("App started at ", time.Now())
	if err := r.Run(fmt.Sprintf("%s:%d", config.AppConfig.ServerHost, config.AppConfig.ServerPort)); err != nil {
		panic(err)
	}
}
