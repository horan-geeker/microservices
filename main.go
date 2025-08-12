// @title 微服务接口文档
// @version 1.0
// @description 这是微服务的接口文档
// @termsOfService http://swagger.io/terms/
// @contact.name horanhe
// @contact.url https://github.com/horan-geeker
// @contact.email horanhe@email.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:80
// @BasePath /
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
