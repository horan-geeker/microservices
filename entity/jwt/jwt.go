package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"microservices/entity/consts"
	"time"
)

// Jwt .
type Jwt struct {
	opt *Options
}

// AuthClaims .
type AuthClaims struct {
	Uid uint64 `json:"uid"`
	jwt.StandardClaims
}

// DecodeToken .
func (j *Jwt) DecodeToken(tokenString string) (*AuthClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(j.opt.Key), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("token is invalid")
}

// GenerateToken .
func (j *Jwt) GenerateToken(id uint64) (string, error) {
	c := AuthClaims{
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(consts.UserTokenExpiredIn).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(j.opt.Key))
}

// NewJwt .
func NewJwt(opt *Options) *Jwt {
	return &Jwt{opt}
}
