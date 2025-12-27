package controller

import (
	"bytes"
	"github.com/go-pay/gopay/alipay"
	"io"
	"microservices/entity/request"
	"microservices/entity/response"
	"microservices/logic"
	"microservices/pkg/log"

	"github.com/gin-gonic/gin"
)

type CallbackController interface {
	GoogleAuthCallback(ctx *gin.Context, param *request.GoogleCallbackParam) (*response.GoogleCallback, error)
	AlipayNotify(c *gin.Context) (*string, error)
	AlipayCallback(c *gin.Context) (*response.AlipayCallback, error)
	AppleCallback(c *gin.Context, payload *request.AppleIAPNotification) (*response.AppleCallback, error)
}

type callbackController struct {
	logic logic.Factory
}

func NewCallbackController(l logic.Factory) CallbackController {
	return &callbackController{logic: l}
}

// GoogleAuthCallback handles the Google OAuth callback.
// @Summary Google OAuth Callback
// @Description Handle Google OAuth callback code
// @Tags callback
// @Accept  json
// @Produce  json
// @Param code query string true "OAuth Code"
// @Success 200 {object} entity.Response[response.GoogleCallback]
// @Failure 400 {object} entity.Response[any]
// @Router /callback/google-auth [get]
func (c *callbackController) GoogleAuthCallback(ctx *gin.Context, param *request.GoogleCallbackParam) (*response.GoogleCallback, error) {
	user, token, err := c.logic.Callback().GoogleCallback(ctx, param.Code)
	if err != nil {
		return nil, err
	}

	return &response.GoogleCallback{
		User:  user,
		Token: token,
	}, nil
}

// AlipayNotify godoc
// @Summary 支付宝异步通知回调
// @Description 接收支付宝网关推送的服务器端通知消息
// @Tags 回调
// @Accept  application/x-www-form-urlencoded
// @Produce json
// @Success 200 {string} success
// @Failure 400 {string} fail
// @Router /callback/alipay-notify [post]
func (c *callbackController) AlipayNotify(ctx *gin.Context) (*string, error) {
	// 读取原始 body
	body, err := ctx.GetRawData()
	if err != nil {
		return nil, err
	}
	log.Info(ctx, "pay", map[string]any{
		"params": string(body),
	})
	// 重新放回去，防止后续无法再读 body
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	// 转成 BodyMap
	bm, err := alipay.ParseNotifyToBodyMap(ctx.Request)
	if err != nil {
		return nil, err
	}
	payload := make(map[string]string)
	for key, values := range ctx.Request.PostForm {
		if len(values) > 0 {
			payload[key] = values[0]
		}
	}

	if err := c.logic.Callback().VerifyAlipayNotifySign(ctx.Request.Context(), bm); err != nil {
		return nil, err
	}

	if err := c.logic.Callback().HandleAlipayNotify(ctx.Request.Context(), payload); err != nil {
		return nil, err
	}
	success := "success"
	return &success, nil
}

// AlipayCallback godoc
// @Summary 支付宝授权回调
// @Description 处理支付宝授权完成后的浏览器回跳
// @Tags 回调
// @Produce json
// @Success 200 {object} entity.Response[response.AlipayCallback]
// @Router /callback/alipay [get]
func (c *callbackController) AlipayCallback(ctx *gin.Context) (*response.AlipayCallback, error) {
	params := make(map[string]string)
	for key, values := range ctx.Request.URL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	if err := c.logic.Callback().HandleAlipayCallback(ctx.Request.Context(), params); err != nil {
		return nil, err
	}

	return &response.AlipayCallback{Message: "success"}, nil
}

// AppleCallback godoc
// @Summary Apple 内购通知回调
// @Description 接收来自苹果内购的服务器回调通知(App Store Server Notifications V2)
// @Tags 回调
// @Accept json
// @Produce json
// @Param payload body request.AppleIAPNotification true "苹果内购回调 payload (包含 signedPayload JWS)"
// @Success 200 {object} entity.Response[response.AppleCallback]
// @Failure 400 {string} string "invalid JWS signature or payload"
// @Router /callback/apple [post]
func (c *callbackController) AppleCallback(ctx *gin.Context, payload *request.AppleIAPNotification) (*response.AppleCallback, error) {
	if err := c.logic.Callback().HandleAppleCallback(ctx.Request.Context(), payload); err != nil {
		return nil, err
	}

	return &response.AppleCallback{Message: "success"}, nil
}
