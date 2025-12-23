package response

import "microservices/entity/model"

// GoogleCallback defines response for Google OAuth callback
type GoogleCallback struct {
	User  *model.User `json:"user"`
	Token string      `json:"token"`
}
