package ecode

import (
	"errors"
	errors2 "microservices/errors"
	"net/http"
)

var (
	ErrUserPasswordDuplicate = errors.New("设置新密码不能与旧密码相同")
)

const (
	_ = iota + 1000
	UserPasswordDuplicate
)

func init() {
	errors2.Register(ErrUserPasswordDuplicate, UserPasswordDuplicate, http.StatusOK)
}
