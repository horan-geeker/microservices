package ecode

import (
	"errors"
	"microservices/pkg/ecode"
	"net/http"
)

var (
	ErrParamInvalid    = errors.New("invalid param")
	ErrTooManyRequests = errors.New("too many requests")
)

const (
	_ = iota + 4000
	ParamInvalid
	TooManyRequests
)

func init() {
	ecode.Register(ErrParamInvalid, ParamInvalid, http.StatusBadRequest)
	ecode.Register(ErrTooManyRequests, TooManyRequests, http.StatusTooManyRequests)
}
