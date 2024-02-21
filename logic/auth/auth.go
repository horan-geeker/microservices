package auth

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"microservices/entity/consts"
	ecode2 "microservices/entity/ecode"
	"microservices/entity/jwt"
	meta2 "microservices/entity/meta"
	"microservices/entity/options"
	"microservices/pkg/util"
	"microservices/repository"
	"time"
)

// AuthLogicInterface defines functions used to handle user api.
type AuthLogicInterface interface {
	Login(ctx context.Context, name, email, phone, password, smsCode, emailCode *string) (*meta2.User, string, error)
	Logout(ctx context.Context, uid uint64) error
	Register(ctx context.Context, name, email, phone, password *string) (*meta2.User, string, error)
	ChangePassword(ctx context.Context, uid uint64, newPassword string, oldPassword string) error
	ChangePasswordByPhone(ctx context.Context, newPassword, phone, smsCode string) error
	VerifyPassword(password string, inputPasswd string) bool
	VerifySmsCode(ctx context.Context, phone string, smsCode string) error
	GeneratePasswordHash(password string) string
	GetAuthUser(ctx context.Context) (*meta2.AuthClaims, error)
}

type authLogic struct {
	repository repository.Factory
}

func NewAuth(repository repository.Factory) AuthLogicInterface {
	return &authLogic{
		repository: repository,
	}
}

// GetUserByIdentity .
func (a *authLogic) GetUserByIdentity(ctx context.Context, name, email, phone *string) (*meta2.User, error) {
	if name == nil && email == nil && phone == nil {
		return nil, errors.New("必须传入一个标识用户的参数")
	}
	var user *meta2.User
	var err error
	switch {
	case name != nil:
		user, err = a.repository.Users().GetByName(ctx, *name)
	case email != nil:
		user, err = a.repository.Users().GetByEmail(ctx, *email)
	case phone != nil:
		user, err = a.repository.Users().GetByPhone(ctx, *phone)
	}
	return user, err
}

// Login .
func (a *authLogic) Login(ctx context.Context, name, email, phone, password, smsCode, emailCode *string) (*meta2.User, string, error) {
	user, err := a.GetUserByIdentity(ctx, name, email, phone)
	if err != nil {
		return nil, "", err
	}
	if password != nil {
		if !a.VerifyPassword(user.Password, *password) {
			return nil, "", ecode2.ErrInvalidPassword
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
	token, err := jwt.NewJwt(options.NewJwtOptions()).GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}
	if err := a.repository.Users().SetToken(ctx, user.ID, token); err != nil {
		return nil, "", err
	}
	if err := a.repository.Users().Update(ctx, user.ID, map[string]any{
		"login_at": time.Now(),
	}); err != nil {
		return nil, "", err
	}
	return user, token, nil
}

// Logout .
func (a *authLogic) Logout(ctx context.Context, uid uint64) error {
	if err := a.repository.Users().DeleteToken(ctx, uid); err != nil {
		return err
	}
	return nil
}

// Register .
func (a *authLogic) Register(ctx context.Context, name, email, phone, password *string) (*meta2.User, string, error) {
	var userName, userEmail, userPhone, userpPassword string
	if name != nil {
		userName = *name
		exist, err := a.repository.Users().GetByName(ctx, *name)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", err
		}
		if exist.ID != 0 {
			return nil, "", ecode2.ErrUserNameAlreadyExist
		}
	}
	if email != nil {
		userEmail = *email
		exist, err := a.repository.Users().GetByEmail(ctx, *email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", err
		}
		if exist.ID != 0 {
			return nil, "", ecode2.ErrUserEmailAlreadyExist
		}
	}
	if phone != nil {
		userPhone = *phone
		exist, err := a.repository.Users().GetByPhone(ctx, *phone)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", err
		}
		if exist.ID != 0 {
			return nil, "", ecode2.ErrUserPhoneAlreadyExist
		}
	}
	if password != nil {
		userpPassword = *password
	}
	user := &meta2.User{
		Name:     userName,
		Email:    userEmail,
		Phone:    userPhone,
		Password: a.GeneratePasswordHash(userpPassword),
		Status:   consts.UserStatusNormal,
		LoginAt:  time.Now(),
	}
	if err := a.repository.Users().Create(ctx, user); err != nil {
		return nil, "", err
	}
	token, err := jwt.NewJwt(options.NewJwtOptions()).GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}
	if err := a.repository.Users().SetToken(ctx, user.ID, token); err != nil {
		return nil, "", err
	}
	return user, token, nil
}

// ChangePassword .
func (a *authLogic) ChangePassword(ctx context.Context, uid uint64, newPassword string, oldPassword string) error {
	user, err := a.repository.Users().GetByUid(ctx, uid)
	if err != nil {
		return err
	}
	if !a.VerifyPassword(user.Password, oldPassword) {
		return ecode2.ErrInvalidPassword
	}
	if err := a.repository.Users().Update(ctx, uid, map[string]any{
		"password": a.GeneratePasswordHash(newPassword),
	}); err != nil {
		return err
	}
	return nil
}

// ChangePasswordByPhone .
func (a *authLogic) ChangePasswordByPhone(ctx context.Context, newPassword, phone, smsCode string) error {
	if err := a.VerifySmsCode(ctx, phone, smsCode); err != nil {
		return err
	}
	user, err := a.repository.Users().GetByPhone(ctx, phone)
	if err != nil {
		return err
	}
	if err := a.repository.Users().Update(ctx, user.ID, map[string]any{
		"password": a.GeneratePasswordHash(newPassword),
	}); err != nil {
		return err
	}
	return nil
}

// VerifyPassword .
func (a *authLogic) VerifyPassword(password string, inputPasswd string) bool {
	if password != a.GeneratePasswordHash(inputPasswd) {
		return false
	}
	return true
}

// VerifySmsCode .
func (a *authLogic) VerifySmsCode(ctx context.Context, phone string, smsCode string) error {
	code, err := a.repository.Auth().GetSmsCode(ctx, phone)
	if err != nil {
		return err
	}
	if smsCode != code {
		return ecode2.ErrUserSmsCodeError
	}
	if err := a.repository.Auth().DeleteSmsCode(ctx, phone); err != nil {
		return err
	}
	return nil
}

// VerifyEmailCode .
func (a *authLogic) VerifyEmailCode(ctx context.Context, email string, code string) error {
	cachedCode, err := a.repository.Auth().GetEmailCode(ctx, code)
	if err != nil {
		return err
	}
	if cachedCode != code {
		return ecode2.ErrUserEmailCodeError
	}
	if err := a.repository.Auth().DeleteEmailCode(ctx, email); err != nil {
		return err
	}
	return nil
}

// GeneratePasswordHash .
func (a *authLogic) GeneratePasswordHash(password string) string {
	return util.MD5(password)
}

func (a *authLogic) GetAuthUser(ctx context.Context) (*meta2.AuthClaims, error) {
	authValue := ctx.Value("auth")
	if authValue == nil {
		return nil, ecode2.ErrTokenInternalNotSet
	}
	authClaims, ok := authValue.(*meta2.AuthClaims)
	if !ok {
		return nil, ecode2.ErrTokenInternalNotSet
	}
	return authClaims, nil
}
