package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"microservices/internal/config"
	"microservices/internal/router"
	"time"
)

func main() {
	r := gin.Default()
	router.Register(r)
	log.Println("App started at ", time.Now())
	if err := r.Run(fmt.Sprintf("%s:%d", config.AppConfig.ServerHost, config.AppConfig.ServerPort)); err != nil {
		panic(err)
	}
}
