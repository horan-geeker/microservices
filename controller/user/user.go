package user

import (
	"github.com/gin-gonic/gin"
	"microservices/cache"
	"microservices/entity/ecode"
	"microservices/entity/request"
	"microservices/entity/response"
	"microservices/logic"
	"microservices/model"
	"microservices/service"
)

type Controller interface {
	Edit(c *gin.Context, param *request.EditUserParam) (*response.EditUser, error)
	Get(c *gin.Context, uid int) (*response.GetUser, error)
}

type controller struct {
	logic logic.Factory
}

func (u *controller) Edit(c *gin.Context, param *request.EditUserParam) (*response.EditUser, error) {
	token := c.GetHeader("Authorization")
	if len(token) == 0 {
		return nil, ecode.ErrTokenIsEmpty
	}
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

// Get .
func (u *controller) Get(c *gin.Context, uid int) (*response.GetUser, error) {
	userinfo, err := u.logic.User().GetByUid(c.Request.Context(), uid)
	if err != nil {
		return nil, err
	}
	return &response.GetUser{
		User: userinfo,
	}, nil
}

// NewController .
func NewController(model model.Factory, cache cache.Factory, service service.Factory) Controller {
	return &controller{
		logic: logic.NewLogic(model, cache, service),
	}
}
