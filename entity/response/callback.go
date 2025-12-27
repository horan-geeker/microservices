package response

import "microservices/entity/model"

// GoogleCallback defines response for Google OAuth callback
type GoogleCallback struct {
	User  *model.User `json:"user"`
	Token string      `json:"token"`
}

type AlipayNotify struct {
	Message string `json:"message"`
}

type AlipayCallback struct {
	Message string `json:"message"`
}

type AppleCallback struct {
	Message string `json:"message"`
}
