package router

import (
	"microservices/cache"
	"microservices/controller"
	"microservices/model"
	"microservices/router/middleware"
	"microservices/service"
)

func init() {
	fileController := controller.NewFileController(model.NewFactory(), cache.NewFactory(), service.NewFactory())

	// 文件上传接口 - 需要认证和每秒1次的用户限频
	router.POST("/files/upload",
		middleware.Authenticate(),
		middleware.ReqRateLimit(10, 60, true, false),
		fileController.Upload)

	router.GET("/files",
		middleware.Authenticate(),
		fileController.List)

	router.GET("/files/:id",
		middleware.Authenticate(),
		fileController.Detail)

}
