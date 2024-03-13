package auth

import (
	"github.com/gin-gonic/gin"
	"microservices/api"
	"microservices/cache"
	"microservices/logic"
	"microservices/model"
	"microservices/service"
	"time"
)

type Controller interface {
	ChangePassword(c *gin.Context, param *api.ChangePasswordParams) (map[string]any, error)
	ChangePasswordByPhone(c *gin.Context, param *api.ChangePasswordByPhoneParams) (map[string]any, error)
	Login(c *gin.Context, params *api.LoginParams) (map[string]any, error)
	Logout(c *gin.Context) (map[string]any, error)
}

type controller struct {
	logic logic.Factory
}

// ChangePassword .
func (a *controller) ChangePassword(c *gin.Context, param *api.ChangePasswordParams) (map[string]any, error) {
	auth, err := a.logic.Auth().GetAuthUser(c.Request.Context())
	if err != nil {
		return nil, err
	}
	if err := a.logic.Auth().ChangePassword(c.Request.Context(), auth.Uid, param.NewPassword, param.OldPassword); err != nil {
		return nil, err
	}
	return nil, nil
}

func (a *controller) ChangePasswordByPhone(c *gin.Context, param *api.ChangePasswordByPhoneParams) (map[string]any, error) {
	if err := a.logic.Auth().ChangePasswordByPhone(c.Request.Context(), param.NewPassword, param.Phone, param.SmsCode); err != nil {
		return nil, err
	}
	return nil, nil
}

// Login .
func (a *controller) Login(c *gin.Context, params *api.LoginParams) (map[string]any, error) {
	user, token, err := a.logic.Auth().Login(c.Request.Context(), params.Name, params.Email, params.Phone, params.Password,
		params.SmsCode, params.EmailCode)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"user": map[string]any{
			"id":          user.ID,
			"status":      user.Status,
			"name":        user.Name,
			"email":       user.Email,
			"phone":       user.Phone,
			"lastLoginAt": user.LoginAt.Format(time.DateTime),
		},
		"token": token,
	}, nil
}

// Logout .
func (a *controller) Logout(c *gin.Context) (map[string]any, error) {
	auth, err := a.logic.Auth().GetAuthUser(c.Request.Context())
	if err != nil {
		return nil, err
	}
	if err := a.logic.Auth().Logout(c.Request.Context(), auth.Uid); err != nil {
		return nil, err
	}
	return nil, nil
}

func NewController(model model.Factory, cache cache.Factory, serviceFactory service.Factory) Controller {
	return &controller{
		logic: logic.NewLogic(model, cache, serviceFactory),
	}
}
