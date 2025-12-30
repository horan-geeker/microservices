package request

type GetOrderListRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

type CreateAlipayPrepay struct {
	Price       int    `json:"price" binding:"required,min=1,max=3"`
	Platform    string `json:"platform" binding:"required,oneof=APP"`
	Channel     string `json:"channel" binding:"required"`
	Description string `json:"description"`
}

type CreateStripeCheckoutRequest struct {
	PriceID string `json:"priceId" binding:"required"`
}
