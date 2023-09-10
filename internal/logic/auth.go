package logic

import (
	"context"
	"microservices/internal/store"
)

// AuthLogicInterface defines functions used to handle user api.
type AuthLogicInterface interface {
	AuthorizeCookie(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error)
	ChangePassword(ctx context.Context, user *v1.User) error
}

type authLogic struct {
	store store.Factory
}

// AuthorizeCookie .
func (a *authLogic) AuthorizeCookie(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error) {
	return a.store.Users().List(ctx, opts)
}

// ChangePassword .
func (a *authLogic) ChangePassword(ctx context.Context, user *v1.User) error {
	return nil
}

func newAuth(l *logic) AuthLogicInterface {
	return &authLogic{
		store: l.store,
	}
}
