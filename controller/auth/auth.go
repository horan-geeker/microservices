package auth

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

type Controller interface {
	ChangePassword(c *gin.Context, param *request.ChangePassword) (*response.ChangePassword, error)
	ChangePasswordByPhone(c *gin.Context, param *request.ChangePasswordByPhone) (*response.ChangePasswordByPhone, error)
	Login(c *gin.Context, params *request.Login) (*response.Login, error)
	Logout(c *gin.Context) (*response.Logout, error)
	Register(c *gin.Context, params *request.Register) (*response.Register, error)
}

type controller struct {
	logic logic.Factory
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 用户修改自己的登录密码
// @Tags 用户
// @Accept json
// @Produce json
// @Param 请求体 body request.ChangePassword true "修改密码请求体"
// @Success 200 {object} entity.Response[response.ChangePassword]
// @Failure 400 {object} entity.Response[any]
// @Router /auth/change-password [post]
func (a *controller) ChangePassword(c *gin.Context, param *request.ChangePassword) (*response.ChangePassword, error) {
	auth, err := a.logic.Auth().GetAuthUser(c.Request.Context())
	if err != nil {
		return nil, err
	}
	if err := a.logic.Auth().ChangePassword(c.Request.Context(), auth.Uid, param.NewPassword, param.OldPassword); err != nil {
		return nil, err
	}
	return nil, nil
}

func (a *controller) ChangePasswordByPhone(c *gin.Context, param *request.ChangePasswordByPhone) (*response.ChangePasswordByPhone, error) {
	if err := a.logic.Auth().ChangePasswordByPhone(c.Request.Context(), param.NewPassword, param.Phone, param.SmsCode); err != nil {
		return nil, err
	}
	return nil, nil
}

// Login .
func (a *controller) Login(c *gin.Context, params *request.Login) (*response.Login, error) {
	user, token, err := a.logic.Auth().Login(c.Request.Context(), params.Name, params.Email, params.Phone, params.Password,
		params.SmsCode, params.EmailCode)
	if err != nil {
		return nil, err
	}
	return &response.Login{
		UserId:      user.ID,
		Status:      user.Status,
		Name:        user.Name,
		Email:       user.Email,
		Phone:       user.Phone,
		LastLoginAt: user.LoginAt.Unix(),
		Token:       token,
	}, nil
}

// Logout .
func (a *controller) Logout(c *gin.Context) (*response.Logout, error) {
	auth, err := a.logic.Auth().GetAuthUser(c.Request.Context())
	if err != nil {
		return nil, err
	}
	if err := a.logic.Auth().Logout(c.Request.Context(), auth.Uid); err != nil {
		return nil, err
	}
	return nil, nil
}

// Register .
func (a *controller) Register(c *gin.Context, params *request.Register) (*response.Register, error) {
	user, token, err := a.logic.Auth().Register(c.Request.Context(), params.Name, params.Email, params.Phone, params.Password)
	if err != nil {
		return nil, err
	}
	return &response.Register{
		UserId: user.ID,
		Token:  token,
	}, nil
}

func NewController(model model.Factory, cache cache.Factory, service service.Factory) Controller {
	return &controller{
		logic: logic.NewLogic(model, cache, service),
	}
}
