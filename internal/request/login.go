package request

// LoginParams .
type LoginParams struct {
	Username string `validate:"required" json:"username"`
	Password string `validate:"required" json:"password"`
}
