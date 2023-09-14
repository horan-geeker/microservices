package store

import (
	"context"
	"microservices/internal/model"
)

// UserStore defines the user storage interface.
type UserStore interface {
	//Create(ctx context.Context, user *model.User) error
	GetByUid(ctx context.Context, uid int) (*model.User, error)
	//List(ctx context.Context) ([]model.User, error)
}
