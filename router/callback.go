package router

import (
	"microservices/cache"
	"microservices/controller"
	"microservices/logic"
	"microservices/model"
	"microservices/service"
)

func init() {
	l := logic.NewLogic(model.NewFactory(), cache.NewFactory(), service.NewFactory())
	c := controller.NewCallbackController(l)
	router.GET("/callback/google-auth", c.GoogleAuthCallback)
	router.POST("/callback/alipay-notify", c.AlipayNotify)
	router.GET("/callback/alipay", c.AlipayCallback)
	router.POST("/callback/apple", c.AppleCallback)
	router.POST("/callback/stripe", c.StripeCallback)
}
