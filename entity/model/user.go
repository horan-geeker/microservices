package model

import (
	"time"
)

// User represents a user gorm model.
type User struct {
	ID        int       `json:"id,omitempty" gorm:"primary_key;AUTO_INCREMENT;column:id"`
	Status    int       `json:"status" gorm:"column:status"`
	Name      string    `json:"name" gorm:"column:name"`
	Password  string    `json:"-" gorm:"column:password"`
	Email     string    `json:"email" gorm:"column:email"`
	Phone     string    `json:"phone" gorm:"column:phone"`
	IsAdmin   int       `json:"isAdmin,omitempty" gorm:"column:is_admin"`
	LoginAt   time.Time `json:"loginAt,omitempty" gorm:"column:login_at"`
	CreatedAt time.Time `json:"createdAt,omitempty" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" gorm:"column:updated_at"`
	// DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deletedAt;index:idx_deletedAt"`
}
