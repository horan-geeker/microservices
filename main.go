package main

import (
	"context"
	_ "microservices/internal/command"
	"microservices/internal/config"
	_ "microservices/internal/router"
	"microservices/pkg/app"
)

// 初始化资源
var (
	env         = config.GetConfig()
	application = app.NewApp(app.NewServerOptions(env.AppEnv, env.ServerHost, env.ServerPort, env.ServerTimeout))
)

func main() {
	if err := application.Running(context.Background()); err != nil {
		panic(err)
	}
}
