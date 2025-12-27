package consts

// OrderStatus 订单状态
type OrderStatus int

const (
	OrderStatusPending        OrderStatus = 0 // 待支付
	OrderStatusPaid           OrderStatus = 1 // 已支付
	OrderStatusProcessSuccess OrderStatus = 2 // 已完成充值，续费
	OrderStatusCreateError    OrderStatus = 3 // 订单创建失败

	TradeTypeAlipay = "alipay"
	TradeTypeApple  = "apple"
)
