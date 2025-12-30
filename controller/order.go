package controller

import (
	"microservices/cache"
	entity "microservices/entity/model"
	"microservices/entity/request"
	"microservices/entity/response"
	"microservices/logic"
	"microservices/model"
	"microservices/service"
	"microservices/util"

	"github.com/gin-gonic/gin"
)

type OrderController interface {
	GetDetail(c *gin.Context, id int) (*entity.Order, error)
	GetList(c *gin.Context, req *request.GetOrderListRequest) (*response.GetOrderListResponse, error)
	CreateAlipayPrepay(c *gin.Context, req *request.CreateAlipayPrepay) (*response.CreateAlipayPrepay, error)
	VerifyAppleReceipt(c *gin.Context, param *request.AppleVerifyReceiptParam) (*response.AppleVerifyReceipt, error)
	CreateStripeCheckout(c *gin.Context, req *request.CreateStripeCheckoutRequest) (*response.CreateStripeCheckoutResponse, error)
}

type orderController struct {
	logic logic.Factory
}

func NewOrderController(model model.Factory, cache cache.Factory, service service.Factory) OrderController {
	return &orderController{
		logic: logic.NewLogic(model, cache, service),
	}
}

func (ctrl *orderController) GetDetail(c *gin.Context, id int) (*entity.Order, error) {
	auth, err := ctrl.logic.Auth().GetAuthUser(c.Request.Context())
	if err != nil {
		return nil, err
	}

	order, err := ctrl.logic.Order().GetDetail(c.Request.Context(), id, auth.Uid)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, nil
	}

	if order.UserId != auth.Uid {
		return nil, nil // Treat as not found for security
	}

	return order, nil
}

func (ctrl *orderController) GetList(c *gin.Context, req *request.GetOrderListRequest) (*response.GetOrderListResponse, error) {
	auth, err := ctrl.logic.Auth().GetAuthUser(c.Request.Context())
	if err != nil {
		return nil, err
	}

	return ctrl.logic.Order().GetList(c.Request.Context(), auth.Uid, req)
}

// CreatePrepayOrder godoc
// @Summary 创建预支付订单
// @Description 创建微信支付预支付订单
// @Tags 支付
// @Accept  json
// @Produce  json
// @Param param body request.CreatePrepayOrderParam true "预支付订单参数"
// @Success 200 {object} entity.Response[response.CreatePrepayOrder]
// @Failure 400 {object} entity.Response[any]
// @Failure 401 {object} entity.Response[any] "用户登录态校验失败"
// @Router /order/pay/alipay [post]
func (p *orderController) CreateAlipayPrepay(c *gin.Context, req *request.CreateAlipayPrepay) (*response.CreateAlipayPrepay, error) {
	clientIP := util.GetUserIP(c)
	auth, err := p.logic.Auth().GetAuthUser(c.Request.Context())
	if err != nil {
		return nil, err
	}
	user, err := p.logic.User().GetByUid(c.Request.Context(), auth.Uid)
	if err != nil {
		return nil, err
	}
	order, orderString, err := p.logic.Order().CreateAliPayPrepayOrder(c.Request.Context(), user.ID, req.Price, req.Platform,
		req.Channel, clientIP, req.Description)
	if err != nil {
		return nil, err
	}
	return &response.CreateAlipayPrepay{
		OrderId:     order.ID,
		OrderString: orderString,
	}, nil
}

// VerifyAppleReceipt godoc
// @Summary 验证苹果内购收据
// @Description 调用苹果 verifyReceipt 接口校验客户端上传的内购收据
// @Tags 支付
// @Accept  json
// @Produce  json
// @Param param body request.AppleVerifyReceiptParam true "苹果内购收据参数"
// @Success 200 {object} entity.Response[response.AppleVerifyReceipt]
// @Failure 400 {object} entity.Response[any]
// @Failure 401 {object} entity.Response[any] "用户登录态校验失败"
// @Router /order/pay/apple-verify-receipt [post]
func (p *orderController) VerifyAppleReceipt(c *gin.Context, param *request.AppleVerifyReceiptParam) (*response.AppleVerifyReceipt, error) {
	auth, err := p.logic.Auth().GetAuthUser(c.Request.Context())
	if err != nil {
		return nil, err
	}
	user, err := p.logic.User().GetByUid(c.Request.Context(), auth.Uid)
	if err != nil {
		return nil, err
	}

	return p.logic.Order().VerifyAppleReceipt(c.Request.Context(), user.ID, param.ReceiptData, param.ExcludeOldTransactions)
}

// CreateStripeCheckout godoc
// @Summary Create Stripe Checkout Session
// @Description Initiate Stripe checkout process
// @Tags Payment
// @Accept json
// @Produce json
// @Param param body request.CreateStripeCheckoutRequest true "Stripe Checkout Parameters"
// @Success 200 {object} entity.Response[response.CreateStripeCheckoutResponse]
// @Router /order/create-stripe-checkout [post]
func (p *orderController) CreateStripeCheckout(c *gin.Context, req *request.CreateStripeCheckoutRequest) (*response.CreateStripeCheckoutResponse, error) {
	clientIP := util.GetUserIP(c)
	auth, err := p.logic.Auth().GetAuthUser(c.Request.Context())
	if err != nil {
		return nil, err
	}

	order, url, err := p.logic.Order().CreateStripeCheckout(c.Request.Context(), auth.Uid, req.PriceID, clientIP)
	if err != nil {
		return nil, err
	}

	return &response.CreateStripeCheckoutResponse{
		OrderId: order.ID,
		Url:     url,
	}, nil
}
