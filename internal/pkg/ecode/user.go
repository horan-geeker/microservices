package ecode

import (
	"errors"
	errors2 "microservices/errors"
	"net/http"
)

var (
	ErrInvalidPassword       = errors.New("密码错误")
	ErrUserPasswordDuplicate = errors.New("设置新密码不能与旧密码相同")
	ErrUserSmsCodeError      = errors.New("短信验证码错误")
	ErrUserNameAlreadyExist  = errors.New("用户名已存在")
	ErrUserEmailAlreadyExist = errors.New("邮箱已存在")
	ErrUserPhoneAlreadyExist = errors.New("手机号已存在")
)

const (
	_ = iota + 1000
	UserPasswordError
	UserPasswordDuplicate
	UserSmsCodeError
	UserNameAlreadyExist
	UserEmailAlreadyExist
	UserPhoneAlreadyExist
)

func init() {
	errors2.Register(ErrInvalidPassword, UserPasswordError, http.StatusOK)
	errors2.Register(ErrUserPasswordDuplicate, UserPasswordDuplicate, http.StatusOK)
	errors2.Register(ErrUserSmsCodeError, UserSmsCodeError, http.StatusOK)
	errors2.Register(ErrUserNameAlreadyExist, UserNameAlreadyExist, http.StatusOK)
	errors2.Register(ErrUserEmailAlreadyExist, UserEmailAlreadyExist, http.StatusOK)
	errors2.Register(ErrUserPhoneAlreadyExist, UserPhoneAlreadyExist, http.StatusOK)
}
