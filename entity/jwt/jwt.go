package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"microservices/entity/consts"
	"microservices/entity/meta"
	"time"
)

type Jwt struct {
	opt *Options
}

func (j *Jwt) DecodeToken(tokenString string) (*meta.AuthClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &meta.AuthClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(j.opt.Key), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*meta.AuthClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("token is invalid")
}

func (j *Jwt) GenerateToken(id uint64) (string, error) {
	c := meta.AuthClaims{
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(consts.UserTokenExpiredIn).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(j.opt.Key))
}

func NewJwt(opt *Options) *Jwt {
	return &Jwt{opt}
}
