package router

import (
	"microservices/cache"
	"microservices/controller"
	"microservices/model"
	"microservices/router/middleware"
	"microservices/service"
)

func init() {
	orderController := controller.NewOrderController(model.NewFactory(), cache.NewFactory(), service.NewFactory())
	router.GET("/orders/:id", middleware.Authenticate(), orderController.GetDetail)
	router.GET("/orders", middleware.Authenticate(), orderController.GetList)
	router.POST("/orders/pay/alipay", middleware.Authenticate(), orderController.CreateAlipayPrepay)
	router.POST("/orders/pay/apple-verify-receipt", middleware.Authenticate(), orderController.VerifyAppleReceipt)
}
