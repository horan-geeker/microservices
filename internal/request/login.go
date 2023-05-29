package request

// LoginParams .
type LoginParams struct {
	Username string `binding:"required" json:"username"`
	Password string `binding:"required" json:"password"`
}

// ChangePasswordParams .
type ChangePasswordParams struct {
	OldPassword string `json:"oldPassword,omitempty"`
	SmsCode     string `json:"smsCode,omitempty"`
	NewPassword string `binding:"required" json:"newPassword"`
}
