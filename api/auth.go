package api

// LoginParams .
type LoginParams struct {
	Name     *string `binding:"required_without_all=Email Phone" json:"name"`
	Email    *string `binding:"required_without_all=Name Phone,omitempty,email" json:"email"`
	Phone    *string `binding:"required_without_all=Email Name" json:"phone"`
	Password *string `binding:"required_without=SmsCode" json:"password"`
	SmsCode  *string `binding:"required_without=Password" json:"smsCode"`
}

// ChangePasswordParams .
type ChangePasswordParams struct {
	Uid         uint64  `binding:"required"`
	OldPassword *string `binding:"required_without:SmsCode" json:"oldPassword,omitempty"`
	SmsCode     *string `binding:"required_without:OldPassword" json:"smsCode,omitempty"`
	NewPassword string  `binding:"required" json:"newPassword"`
}

// RegisterParams .
type RegisterParams struct {
	Name     *string `binding:"required_without_all=Email Phone" json:"name"`
	Email    *string `binding:"required_without_all=Name Phone,omitempty,email" json:"email"`
	Phone    *string `binding:"required_without_all=Email Name" json:"phone"`
	Password *string `binding:"required_without=SmsCode" json:"password"`
	SmsCode  *string `binding:"required_without=Password" json:"smsCode"`
}
