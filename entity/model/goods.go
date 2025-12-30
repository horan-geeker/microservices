package model

import "time"

// Goods represents a goods gorm model.
type Goods struct {
	ID            int       `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`
	Name          string    `json:"name" gorm:"column:name"`
	Description   string    `json:"description" gorm:"column:description"`
	Price         int       `json:"price" gorm:"column:price"`
	OriginalPrice int       `json:"originalPrice" gorm:"column:original_price"`
	Discount      int       `json:"discount" gorm:"column:discount"`
	Currency      string    `json:"currency" gorm:"column:currency"`
	StripePriceID string    `json:"stripePriceId" gorm:"column:stripe_price_id"`
	Credit        int       `json:"credit" gorm:"column:credit"`
	Status        int       `json:"status" gorm:"column:status"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

// TableName returns the table name
func (Goods) TableName() string {
	return "goods"
}
