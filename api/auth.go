package api

// LoginParams .
type LoginParams struct {
	Name      *string `binding:"required_without_all=Email Phone" json:"name"`
	Email     *string `binding:"required_without_all=Name Phone,omitempty,email" json:"email"`
	Phone     *string `binding:"required_without_all=Email Name" json:"phone"`
	Password  *string `binding:"required_without=SmsCode" json:"password"`
	SmsCode   *string `binding:"required_without=Password" json:"smsCode"`
	EmailCode *string `binding:"required_without=Password" json:"emailCode"`
}

// ChangePasswordParams .
type ChangePasswordParams struct {
	OldPassword string `binding:"required" json:"oldPassword"`
	NewPassword string `binding:"required" json:"newPassword"`
}

// ChangePasswordByPhoneParams .
type ChangePasswordByPhoneParams struct {
	Phone       string `binding:"required" json:"phone"`
	SmsCode     string `binding:"required" json:"smsCode"`
	NewPassword string `binding:"required" json:"newPassword"`
}

// RegisterParams .
type RegisterParams struct {
	Name     *string `binding:"required_without_all=Email Phone" json:"name"`
	Email    *string `binding:"required_without_all=Name Phone,omitempty,email" json:"email"`
	Phone    *string `binding:"required_without_all=Email Name" json:"phone"`
	Password *string `binding:"required_without=SmsCode" json:"password"`
	SmsCode  *string `binding:"required_without=Password" json:"smsCode"`
}
