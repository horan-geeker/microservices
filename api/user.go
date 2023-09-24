package api

// GetUserParam .
type GetUserParam struct {
	Uid int `json:"uid" validate:"required"`
}

// EditUserParam .
type EditUserParam struct {
	Name  *string `json:"name" validate:"omitempty,min=2,max=20"`
	Email *string `json:"email" validate:"omitempty,email"`
	Phone *string `json:"phone" validate:"omitempty,phone"`
}
