package ecode

import (
	"errors"
	errors2 "microservices/pkg/ecode"
	"net/http"
)

var (
	ErrInvalidPassword       = errors.New("密码错误")
	ErrUserPasswordDuplicate = errors.New("设置新密码不能与旧密码相同")
	ErrUserSmsCodeError      = errors.New("短信验证码错误")
	ErrUserEmailCodeError    = errors.New("邮箱验证码错误")
	ErrUserNameAlreadyExist  = errors.New("用户名已存在")
	ErrUserEmailAlreadyExist = errors.New("邮箱已存在")
	ErrUserPhoneAlreadyExist = errors.New("手机号已存在")
	ErrAccountFrozen         = errors.New("密码失败次数过多，账号已冻结，请两个小时后重试")
)

const (
	_ = iota + 1000
	UserPasswordError
	UserPasswordDuplicate
	UserSmsCodeError
	UserEmailCodeError
	UserNameAlreadyExist
	UserEmailAlreadyExist
	UserPhoneAlreadyExist
	AccountFrozen
)

func init() {
	errors2.Register(ErrInvalidPassword, UserPasswordError, http.StatusOK)
	errors2.Register(ErrUserPasswordDuplicate, UserPasswordDuplicate, http.StatusOK)
	errors2.Register(ErrUserSmsCodeError, UserSmsCodeError, http.StatusOK)
	errors2.Register(ErrUserEmailCodeError, UserEmailCodeError, http.StatusOK)
	errors2.Register(ErrUserNameAlreadyExist, UserNameAlreadyExist, http.StatusOK)
	errors2.Register(ErrUserEmailAlreadyExist, UserEmailAlreadyExist, http.StatusOK)
	errors2.Register(ErrUserPhoneAlreadyExist, UserPhoneAlreadyExist, http.StatusOK)
	errors2.Register(ErrAccountFrozen, AccountFrozen, http.StatusForbidden)
}
