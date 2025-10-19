package ecode

import (
	"microservices/pkg/ecode"
	"net/http"
)

var (
	ErrParamInvalid            = &ecode.CustomError{Code: ParamInvalid, HttpStatus: http.StatusBadRequest}
	ErrTooManyRequests         = &ecode.CustomError{Code: TooManyRequests, HttpStatus: http.StatusTooManyRequests, Message: "Too many requests"}
	ErrGenerationFailed        = &ecode.CustomError{Code: GenerationFailed, HttpStatus: http.StatusInternalServerError}
	ErrGenerationQuotaExceeded = &ecode.CustomError{Code: GenerationQuotaExceeded, HttpStatus: http.StatusOK}
)

const (
	_ = iota + 4000
	ParamInvalid
	TooManyRequests
	GenerationFailed
	GenerationQuotaExceeded
)
