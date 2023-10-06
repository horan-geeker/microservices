package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt"
)

type Jwt[T jwt.Claims] struct {
	opt *Options
}

func (j *Jwt[T]) ParseJWTToken(tokenString string) (*T, error) {
	// 解析token
	var claims T
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(j.opt.Key), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(T); ok && token.Valid { // 校验token
		return &claims, nil
	}
	return nil, errors.New("token is invalid")
}

func NewJwt[T jwt.Claims](opt *Options) *Jwt[T] {
	return &Jwt[T]{opt}
}
