package controller

import (
	"github.com/gin-gonic/gin"
	"microservices/cache"
	_ "microservices/entity"
	"microservices/entity/request"
	"microservices/entity/response"
	"microservices/logic"
	"microservices/model"
	"microservices/service"
)

type UserController interface {
	Edit(c *gin.Context, param *request.EditUserParam) (*response.EditUser, error)
	GetUserInfo(c *gin.Context) (*response.GetUser, error)
}

type userController struct {
	logic logic.Factory
}

// Edit godoc
// @Summary 编辑用户信息
// @Description 编辑用户详细信息
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} entity.Response[response.EditUser]
// @Failure 401 {object} entity.Response[any] "用户登录态校验失败(code: 4)"
// @Router /users/{id}/edit [post]
func (u *userController) Edit(c *gin.Context, param *request.EditUserParam) (*response.EditUser, error) {
	auth, err := u.logic.Auth().GetAuthUser(c)
	if err != nil {
		return nil, err
	}
	err = u.logic.User().Edit(c.Request.Context(), auth.Uid, param.Name, param.Email, param.Phone)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// GetUserInfo godoc
// @Summary 获取用户信息
// @Description 获取用户详细信息
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} entity.Response[response.GetUser]
// @Failure 400 {object} entity.Response[any]
// @Router /userinfo [get]
func (u *userController) GetUserInfo(c *gin.Context) (*response.GetUser, error) {
	auth, err := u.logic.Auth().GetAuthUser(c.Request.Context())
	if err != nil {
		return nil, err
	}
	userinfo, err := u.logic.User().GetByUid(c.Request.Context(), auth.Uid)
	if err != nil {
		return nil, err
	}
	return &response.GetUser{
		User: userinfo,
	}, nil
}

// NewUserController .
func NewUserController(model model.Factory, cache cache.Factory, service service.Factory) UserController {
	return &userController{
		logic: logic.NewLogic(model, cache, service),
	}
}
