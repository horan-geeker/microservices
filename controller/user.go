package controller

import (
	"github.com/gin-gonic/gin"
	"microservices/api"
	"microservices/entity/ecode"
	"microservices/logic"
	"microservices/repository"
	"microservices/service"
)

type UserApi interface {
	Edit(c *gin.Context, param *api.EditUserParam) (map[string]any, error)
	Get(c *gin.Context, uid uint64) (map[string]any, error)
	Register(c *gin.Context, params *api.RegisterParams) (map[string]any, error)
}

type UserController struct {
	logic logic.LogicInterface
}

func (u *UserController) Edit(c *gin.Context, param *api.EditUserParam) (map[string]any, error) {
	token := c.GetHeader("Authorization")
	if len(token) == 0 {
		return nil, ecode.ErrTokenIsEmpty
	}
	auth, err := u.logic.Auth().GetAuthUser(c)
	if err != nil {
		return nil, err
	}
	err = u.logic.Users().Edit(c.Request.Context(), auth.Uid, param.Name, param.Email, param.Phone)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Get .
func (u *UserController) Get(c *gin.Context, uid uint64) (map[string]any, error) {
	userinfo, err := u.logic.Users().GetByUid(c.Request.Context(), uid)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"user": userinfo,
	}, nil
}

// Register .
func (u *UserController) Register(c *gin.Context, params *api.RegisterParams) (map[string]any, error) {
	user, token, err := u.logic.Auth().Register(c.Request.Context(), params.Name, params.Email, params.Phone, params.Password)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"user": map[string]any{
			"id": user.ID,
		},
		"token": token,
	}, nil
}

// NewUserController .
func NewUserController(repositoryFactory repository.Factory, serviceFactory service.Factory) UserApi {
	return &UserController{
		logic: logic.NewLogic(repositoryFactory, serviceFactory),
	}
}
