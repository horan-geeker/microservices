package response

import "microservices/entity/model"

type EditUser struct{}
type GetUser struct {
	User *model.User `json:"user"`
}
