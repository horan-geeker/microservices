package response

import "microservices/entity/model"

type GetOrderListResponse struct {
	Items []*model.Order `json:"items"`
	Total int64          `json:"total"`
}

type GetOrderDetailResponse struct {
	*model.Order
}

// AppleVerifyReceipt represents verification result of an Apple in-app purchase receipt.
type AppleVerifyReceipt struct {
	Status        int                `json:"status"`
	Environment   string             `json:"environment"`
	LatestReceipt string             `json:"latestReceipt,omitempty"`
	ReceiptInfo   []AppleReceiptInfo `json:"receiptInfo,omitempty"`
}

// AppleReceiptInfo includes important fields from the latest receipt info returned by Apple.
type AppleReceiptInfo struct {
	ProductID             string `json:"productId"`
	TransactionID         string `json:"transactionId"`
	OriginalTransactionID string `json:"originalTransactionId"`
	PurchaseDateMs        string `json:"purchaseDateMs"`
	ExpiresDateMs         string `json:"expiresDateMs,omitempty"`
	CancellationDateMs    string `json:"cancellationDateMs,omitempty"`
}

type CreateAlipayPrepay struct {
	OrderId     int    `json:"orderId"`
	OrderString string `json:"orderString"`
}
