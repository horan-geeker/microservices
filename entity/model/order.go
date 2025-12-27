package model

import (
	"database/sql"
	"microservices/entity/consts"
	"time"
)

// Order represents an order gorm model.
type Order struct {
	ID          int                 `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`
	MchId       string              `json:"mchId" gorm:"column:mchid"`
	AppId       string              `json:"appId" gorm:"column:appid"`
	OutTradeNo  string              `json:"outTradeNo" gorm:"column:out_trade_no"`
	UserId      int                 `json:"userId" gorm:"column:user_id"`
	TotalAmount uint64              `json:"totalAmount" gorm:"column:total_amount"`
	Currency    string              `json:"currency" gorm:"column:currency;default:CNY"`
	Subject     string              `json:"subject" gorm:"column:subject"`
	Description *string             `json:"description" gorm:"column:description"`
	Status      consts.OrderStatus  `json:"status" gorm:"column:status;default:0"`
	TradeType   string              `json:"tradeType" gorm:"column:trade_type"`
	Platform    string              `json:"platform" gorm:"column:platform"`
	ClientIp    *string             `json:"clientIp" gorm:"column:client_ip"`
	ExpireAt    *time.Time          `json:"expireAt" gorm:"column:expire_at"`
	PaidAt      sql.Null[time.Time] `json:"paidAt" gorm:"column:paid_at"`
	ClosedAt    sql.Null[time.Time] `json:"closedAt" gorm:"column:closed_at"`
	Version     int                 `json:"version" gorm:"column:version;default:0"`
	CreatedAt   time.Time           `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt   time.Time           `json:"updatedAt" gorm:"column:updated_at"`
}

// TableName returns the table name
func (Order) TableName() string {
	return "orders"
}
