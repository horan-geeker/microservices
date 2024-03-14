package auth

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"microservices/cache"
	"microservices/entity/config"
	"microservices/entity/consts"
	"microservices/entity/ecode"
	"microservices/entity/jwt"
	entity "microservices/entity/model"
	"microservices/model"
	"microservices/pkg/util"
	"microservices/service"
	"time"
)

// Logic defines functions used to handle user api.
type Logic interface {
	Login(ctx context.Context, name, email, phone, password, smsCode, emailCode *string) (*entity.User, string, error)
	Logout(ctx context.Context, uid int) error
	Register(ctx context.Context, name, email, phone, password *string) (*entity.User, string, error)
	ChangePassword(ctx context.Context, uid int, newPassword string, oldPassword string) error
	ChangePasswordByPhone(ctx context.Context, newPassword, phone, smsCode string) error
	VerifyPassword(password string, inputPasswd string) bool
	VerifySmsCode(ctx context.Context, phone string, smsCode string) error
	GeneratePasswordHash(password string) string
	GetAuthUser(ctx context.Context) (*jwt.AuthClaims, error)
}

type logic struct {
	model model.Factory
	cache cache.Factory
	srv   service.Factory
}

func NewAuth(model model.Factory, cache cache.Factory, service service.Factory) Logic {
	return &logic{
		model: model,
		cache: cache,
		srv:   service,
	}
}

// GetUserByIdentity .
func (a *logic) GetUserByIdentity(ctx context.Context, name, email, phone *string) (*entity.User, error) {
	if name == nil && email == nil && phone == nil {
		return nil, errors.New("必须传入一个标识用户的参数")
	}
	var user *entity.User
	var err error
	switch {
	case name != nil:
		user, err = a.model.User().GetByName(ctx, *name)
	case email != nil:
		user, err = a.model.User().GetByEmail(ctx, *email)
	case phone != nil:
		user, err = a.model.User().GetByPhone(ctx, *phone)
	}
	return user, err
}

// Login .
func (a *logic) Login(ctx context.Context, name, email, phone, password, smsCode, emailCode *string) (*entity.User,
	string, error) {
	user, err := a.GetUserByIdentity(ctx, name, email, phone)
	if err != nil {
		return nil, "", err
	}
	if password != nil {
		if !a.VerifyPassword(user.Password, *password) {
			return nil, "", ecode.ErrInvalidPassword
		}
	}
	if phone != nil && smsCode != nil {
		if err := a.VerifySmsCode(ctx, user.Phone, *smsCode); err != nil {
			return nil, "", err
		}
	}
	if email != nil && emailCode != nil {
		if err := a.VerifyEmailCode(ctx, user.Email, *emailCode); err != nil {
			return nil, "", err
		}
	}
	token, err := jwt.NewJwt(config.NewJwtOptions()).GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}
	if err := a.cache.User().SetToken(ctx, user.ID, token); err != nil {
		return nil, "", err
	}
	if err := a.model.User().Update(ctx, user.ID, map[string]any{
		"login_at": time.Now(),
	}); err != nil {
		return nil, "", err
	}
	return user, token, nil
}

// Logout .
func (a *logic) Logout(ctx context.Context, uid int) error {
	if err := a.cache.User().DeleteToken(ctx, uid); err != nil {
		return err
	}
	return nil
}

// Register .
func (a *logic) Register(ctx context.Context, name, email, phone, password *string) (*entity.User, string, error) {
	var userName, userEmail, userPhone, userpPassword string
	if name != nil {
		userName = *name
		exist, err := a.model.User().GetByName(ctx, *name)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", err
		}
		if exist.ID != 0 {
			return nil, "", ecode.ErrUserNameAlreadyExist
		}
	}
	if email != nil {
		userEmail = *email
		exist, err := a.model.User().GetByEmail(ctx, *email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", err
		}
		if exist.ID != 0 {
			return nil, "", ecode.ErrUserEmailAlreadyExist
		}
	}
	if phone != nil {
		userPhone = *phone
		exist, err := a.model.User().GetByPhone(ctx, *phone)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", err
		}
		if exist.ID != 0 {
			return nil, "", ecode.ErrUserPhoneAlreadyExist
		}
	}
	if password != nil {
		userpPassword = *password
	}
	user := &entity.User{
		Name:     userName,
		Email:    userEmail,
		Phone:    userPhone,
		Password: a.GeneratePasswordHash(userpPassword),
		Status:   consts.UserStatusNormal,
		LoginAt:  time.Now(),
	}
	if err := a.model.User().Create(ctx, user); err != nil {
		return nil, "", err
	}
	token, err := jwt.NewJwt(config.NewJwtOptions()).GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}
	if err := a.cache.User().SetToken(ctx, user.ID, token); err != nil {
		return nil, "", err
	}
	return user, token, nil
}

// ChangePassword .
func (a *logic) ChangePassword(ctx context.Context, uid int, newPassword string, oldPassword string) error {
	user, err := a.model.User().GetByUid(ctx, uid)
	if err != nil {
		return err
	}
	if !a.VerifyPassword(user.Password, oldPassword) {
		return ecode.ErrInvalidPassword
	}
	if err := a.model.User().Update(ctx, uid, map[string]any{
		"password": a.GeneratePasswordHash(newPassword),
	}); err != nil {
		return err
	}
	return nil
}

// ChangePasswordByPhone .
func (a *logic) ChangePasswordByPhone(ctx context.Context, newPassword, phone, smsCode string) error {
	if err := a.VerifySmsCode(ctx, phone, smsCode); err != nil {
		return err
	}
	user, err := a.model.User().GetByPhone(ctx, phone)
	if err != nil {
		return err
	}
	if err := a.model.User().Update(ctx, user.ID, map[string]any{
		"password": a.GeneratePasswordHash(newPassword),
	}); err != nil {
		return err
	}
	return nil
}

// VerifyPassword .
func (a *logic) VerifyPassword(password string, inputPasswd string) bool {
	if password != a.GeneratePasswordHash(inputPasswd) {
		return false
	}
	return true
}

// VerifySmsCode .
func (a *logic) VerifySmsCode(ctx context.Context, phone string, smsCode string) error {
	code, err := a.cache.Auth().GetSmsCode(ctx, phone)
	if err != nil {
		return err
	}
	if smsCode != code {
		return ecode.ErrUserSmsCodeError
	}
	if err := a.cache.Auth().DeleteSmsCode(ctx, phone); err != nil {
		return err
	}
	return nil
}

// VerifyEmailCode .
func (a *logic) VerifyEmailCode(ctx context.Context, email string, code string) error {
	cachedCode, err := a.cache.Auth().GetEmailCode(ctx, code)
	if err != nil {
		return err
	}
	if cachedCode != code {
		return ecode.ErrUserEmailCodeError
	}
	if err := a.cache.Auth().DeleteEmailCode(ctx, email); err != nil {
		return err
	}
	return nil
}

// GeneratePasswordHash .
func (a *logic) GeneratePasswordHash(password string) string {
	return util.MD5(password)
}

func (a *logic) GetAuthUser(ctx context.Context) (*jwt.AuthClaims, error) {
	authValue := ctx.Value("auth")
	if authValue == nil {
		return nil, ecode.ErrTokenInternalNotSet
	}
	authClaims, ok := authValue.(*jwt.AuthClaims)
	if !ok {
		return nil, ecode.ErrTokenInternalNotSet
	}
	return authClaims, nil
}
