package main

import (
	"context"
	_ "microservices/command"
	"microservices/entity/config"
	"microservices/pkg/app"
	_ "microservices/router"
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
