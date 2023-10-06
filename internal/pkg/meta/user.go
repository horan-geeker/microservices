package meta

import (
	"microservices/pkg/meta"
	"time"
)

// User represents a user gorm model.
type User struct {
	// May add TypeMeta in the future.

	// Standard object's metadata.
	meta.ObjectMeta `json:"metadata,omitempty"`

	Status int `json:"status" gorm:"column:status"`

	Name string `json:"name" gorm:"column:name"`

	Password string `json:"-" gorm:"column:password"`

	Email string `json:"email" gorm:"column:email"`

	Phone string `json:"phone" gorm:"column:phone"`

	IsAdmin int `json:"isAdmin,omitempty" gorm:"column:is_admin"`

	LoginAt time.Time `json:"loginAt,omitempty" gorm:"column:login_at"`
}
