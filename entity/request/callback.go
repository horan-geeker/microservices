package request

// GoogleCallbackParam defines parameters for Google OAuth callback
type GoogleCallbackParam struct {
	Code string `form:"code" binding:"required"`
}
