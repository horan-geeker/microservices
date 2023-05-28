package auth

import (
	"microservices/internal/entity"
	"microservices/internal/request"
)

// Login .
func Login(req request.LoginParams) (*entity.Response, error) {
	return &entity.Response{
		Data: map[string]any{
			"username": req.Username,
		},
		Code: 0,
	}, nil
}
