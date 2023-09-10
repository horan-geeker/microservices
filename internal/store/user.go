package store

import (
	"context"
	"microservices/internal/model"
)

// UserStore defines the user storage interface.
type UserStore interface {
	//Create(ctx context.Context, user *model.User) error
	GetByUserName(ctx context.Context, username string) (*model.User, error)
	//List(ctx context.Context) ([]model.User, error)
}
