package request

// GoogleCallbackParam defines parameters for Google OAuth callback
type GoogleCallbackParam struct {
	Code string `form:"code" binding:"required"`
}

// AppleVerifyReceiptParam represents the payload uploaded by client to verify a receipt.
type AppleVerifyReceiptParam struct {
	ReceiptData            string `json:"receiptData" binding:"required"`
	ExcludeOldTransactions bool   `json:"excludeOldTransactions"`
}

// AppleIAPNotification represents the payload from Apple's in-app purchase server notification.
// According to Apple's App Store Server Notifications V2, the root payload only contains signedPayload.
// Other fields (notificationType, subtype, environment) are inside the decoded JWS payload.
type AppleIAPNotification struct {
	SignedPayload string `json:"signedPayload" binding:"required"`
}
