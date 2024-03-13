package user

import (
	"context"
	"microservices/cache"
	entity "microservices/entity/model"
	"microservices/model"
	"microservices/service"
)

// Logic defines functions used to handle user api.
type Logic interface {
	Create(ctx context.Context, user *entity.User) error
	GetByUid(ctx context.Context, uid uint64) (*entity.User, error)
	List(ctx context.Context) ([]entity.User, error)
	Edit(ctx context.Context, id uint64, name, email, phone *string) error
}

type logic struct {
	model model.Factory
	cache cache.Factory
	srv   service.Factory
}

func NewLogic(model model.Factory, cache cache.Factory, service service.Factory) Logic {
	return &logic{model: model, cache: cache, srv: service}
}

// Create .
func (u *logic) Create(ctx context.Context, user *entity.User) error {
	//return u.model.User().Create(ctx, user)
	return nil
}

// GetByUid .
func (u *logic) GetByUid(ctx context.Context, uid uint64) (*entity.User, error) {
	return u.model.User().GetByUid(ctx, uid)
}

// List .
func (u *logic) List(ctx context.Context) ([]entity.User, error) {
	//return u.model.User().List(ctx)
	return nil, nil
}

// Edit .
func (u *logic) Edit(ctx context.Context, id uint64, name, email, phone *string) error {
	user, err := u.model.User().GetByUid(ctx, id)
	if err != nil {
		return err
	}
	data := make(map[string]any)
	if name != nil {
		data["name"] = *name
	}
	if email != nil {
		data["email"] = *email
	}
	if phone != nil {
		data["phone"] = *phone
	}
	return u.model.User().Update(ctx, user.ID, data)
}
