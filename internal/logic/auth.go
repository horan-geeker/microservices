package logic

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
	"microservices/internal/model"
	"microservices/internal/pkg/consts"
	"microservices/internal/pkg/ecode"
	"microservices/internal/pkg/options"
	"microservices/internal/store"
	"microservices/pkg/meta"
	"microservices/pkg/util"
	"time"
)

// AuthLogicInterface defines functions used to handle user api.
type AuthLogicInterface interface {
	Login(ctx context.Context, name, email, phone, password, smsCode, emailCode *string) (*model.User, string, error)
	Logout(ctx context.Context, uid uint64) error
	Register(ctx context.Context, name, email, phone, password *string) (*model.User, string, error)
	ChangePassword(ctx context.Context, uid uint64, newPassword string, oldPassword, smsCode *string) error
	VerifyPassword(password string, inputPasswd string) bool
	VerifySmsCode(ctx context.Context, uid uint64, smsCode string) error
	GeneratePasswordHash(password string) string
	GenerateJWTToken(id uint64) (string, error)
	GetAuthUser(ctx context.Context) (*meta.AuthClaims, error)
}

type authLogic struct {
	store store.DataFactory
	cache store.CacheFactory
}

func newAuth(l *logic) AuthLogicInterface {
	return &authLogic{
		store: l.store,
		cache: l.cache,
	}
}

// GetUserByIdentity .
func (a *authLogic) GetUserByIdentity(ctx context.Context, name, email, phone *string) (*model.User, error) {
	if name == nil && email == nil && phone == nil {
		return nil, errors.New("必须传入一个标识用户的参数")
	}
	var user *model.User
	var err error
	switch {
	case name != nil:
		user, err = a.store.Users().GetByName(ctx, *name)
	case email != nil:
		user, err = a.store.Users().GetByEmail(ctx, *email)
	case phone != nil:
		user, err = a.store.Users().GetByPhone(ctx, *phone)
	}
	return user, err
}

// Login .
func (a *authLogic) Login(ctx context.Context, name, email, phone, password, smsCode, emailCode *string) (*model.User, string, error) {
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
		if err := a.VerifySmsCode(ctx, user.ID, *smsCode); err != nil {
			return nil, "", err
		}
	}
	if email != nil && emailCode != nil {
		if err := a.VerifyEmailCode(ctx, user.ID, *emailCode); err != nil {
			return nil, "", err
		}
	}
	token, err := a.GenerateJWTToken(user.ID)
	if err != nil {
		return nil, "", err
	}
	if err := a.cache.Users().SetToken(ctx, user.ID, token); err != nil {
		return nil, "", err
	}
	if err := a.store.Users().Update(ctx, user.ID, map[string]any{
		"login_at": time.Now(),
	}); err != nil {
		return nil, "", err
	}
	return user, token, nil
}

// Logout .
func (a *authLogic) Logout(ctx context.Context, uid uint64) error {
	if err := a.cache.Users().DeleteToken(ctx, uid); err != nil {
		return err
	}
	return nil
}

// Register .
func (a *authLogic) Register(ctx context.Context, name, email, phone, password *string) (*model.User, string, error) {
	var userName, userEmail, userPhone, userpPassword string
	if name != nil {
		userName = *name
		exist, err := a.store.Users().GetByName(ctx, *name)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", err
		}
		if exist.ID != 0 {
			return nil, "", ecode.ErrUserNameAlreadyExist
		}
	}
	if email != nil {
		userEmail = *email
		exist, err := a.store.Users().GetByEmail(ctx, *email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", err
		}
		if exist.ID != 0 {
			return nil, "", ecode.ErrUserEmailAlreadyExist
		}
	}
	if phone != nil {
		userPhone = *phone
		exist, err := a.store.Users().GetByPhone(ctx, *phone)
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
	user := &model.User{
		Name:     userName,
		Email:    userEmail,
		Phone:    userPhone,
		Password: a.GeneratePasswordHash(userpPassword),
		Status:   consts.UserStatusNormal,
		LoginAt:  time.Now(),
	}
	if err := a.store.Users().Create(ctx, user); err != nil {
		return nil, "", err
	}
	token, err := a.GenerateJWTToken(user.ID)
	if err != nil {
		return nil, "", err
	}
	if err := a.cache.Users().SetToken(ctx, user.ID, token); err != nil {
		return nil, "", err
	}
	return user, token, nil
}

// ChangePassword .
func (a *authLogic) ChangePassword(ctx context.Context, uid uint64, newPassword string, oldPassword *string, smsCode *string) error {
	if oldPassword != nil {
		user, err := a.store.Users().GetByUid(ctx, uid)
		if err != nil {
			return err
		}
		if !a.VerifyPassword(user.Password, *oldPassword) {
			return ecode.ErrInvalidPassword
		}
	}
	if smsCode != nil {
		if err := a.VerifySmsCode(ctx, uid, *smsCode); err != nil {
			return err
		}
	}
	if err := a.store.Users().Update(ctx, uid, map[string]any{
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
func (a *authLogic) VerifySmsCode(ctx context.Context, uid uint64, smsCode string) error {
	code, err := a.cache.Auth().GetSmsCode(ctx, uid)
	if err != nil {
		return err
	}
	if smsCode != code {
		return ecode.ErrUserSmsCodeError
	}
	if err := a.cache.Auth().DeleteSmsCode(ctx, uid); err != nil {
		return err
	}
	return nil
}

// VerifyEmailCode .
func (a *authLogic) VerifyEmailCode(ctx context.Context, uid uint64, emailCode string) error {
	code, err := a.cache.Auth().GetEmailCode(ctx, uid)
	if err != nil {
		return err
	}
	if emailCode != code {
		return ecode.ErrUserEmailCodeError
	}
	if err := a.cache.Auth().DeleteEmailCode(ctx, uid); err != nil {
		return err
	}
	return nil
}

// GeneratePasswordHash .
func (a *authLogic) GeneratePasswordHash(password string) string {
	return util.MD5(password)
}

func (a *authLogic) GenerateJWTToken(id uint64) (string, error) {
	c := meta.AuthClaims{
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(consts.UserTokenExpiredIn).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	opts := options.NewJwtOptions()
	return token.SignedString([]byte(opts.Key))
}

func (a *authLogic) GetAuthUser(ctx context.Context) (*meta.AuthClaims, error) {
	authValue := ctx.Value("auth")
	if authValue == nil {
		return nil, ecode.ErrTokenInternalNotSet
	}
	authClaims, ok := authValue.(*meta.AuthClaims)
	if !ok {
		return nil, ecode.ErrTokenInternalNotSet
	}
	return authClaims, nil
}
