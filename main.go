package main

import (
	"context"
	_ "microservices/command"
	"microservices/pkg/app"
	_ "microservices/router"
)

func main() {
	if err := app.GetApp().Running(context.Background()); err != nil {
		panic(err)
	}
}
