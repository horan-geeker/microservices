package request

// LoginParams .
type LoginParams struct {
	Username string `binding:"required" json:"username"`
	Password string `binding:"required" json:"password"`
}
