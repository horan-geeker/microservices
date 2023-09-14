package api

type GetUserInfoParam struct {
	Uid int `json:"uid" validate:"required"`
}
