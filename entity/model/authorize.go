package model

import (
	"time"
)

// Authorize represents the authorized user info from third-party providers.
type Authorize struct {
	ID           int       `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`
	UserID       int       `json:"userId" gorm:"column:user_id;index:idx_user_id"`
	Provider     string    `json:"provider" gorm:"column:provider;index:idx_provider_uid,unique"`      // google, github, etc.
	ProviderID   string    `json:"providerId" gorm:"column:provider_id;index:idx_provider_uid,unique"` // openid or unique id from provider
	Email        string    `json:"email" gorm:"column:email"`
	Name         string    `json:"name" gorm:"column:name"`
	Avatar       string    `json:"avatar" gorm:"column:avatar"`
	AccessToken  string    `json:"-" gorm:"column:access_token"`
	RefreshToken string    `json:"-" gorm:"column:refresh_token"`
	ExpiresAt    time.Time `json:"-" gorm:"column:expires_at"`
	JSON         string    `json:"-" gorm:"column:json;type:text"` // Raw response if needed
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"column:updated_at"`
}
