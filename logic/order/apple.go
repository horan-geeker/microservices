package order

import (
	"context"
	"fmt"
	"microservices/entity/config"
	"microservices/entity/response"
	"microservices/pkg/log"
)

func (p *logic) VerifyAppleReceipt(ctx context.Context, userId int, receiptData string, excludeOldTransactions bool) (*response.AppleVerifyReceipt, error) {
	appleOptions := config.NewAppleOptions()
	if appleOptions.SharedSecret == "" {
		return nil, fmt.Errorf("未配置 Apple Shared Secret")
	}

	result, err := p.srv.Apple().VerifyReceipt(ctx, receiptData, excludeOldTransactions)
	if err != nil {
		return nil, err
	}

	receiptInfo := make([]response.AppleReceiptInfo, 0, len(result.LatestReceiptInfo))
	for _, info := range result.LatestReceiptInfo {
		receiptInfo = append(receiptInfo, response.AppleReceiptInfo{
			ProductID:             info.ProductID,
			TransactionID:         info.TransactionID,
			OriginalTransactionID: info.OriginalTransactionID,
			PurchaseDateMs:        info.PurchaseDateMs,
			ExpiresDateMs:         info.ExpiresDateMs,
			CancellationDateMs:    info.CancellationDateMs,
		})
	}

	log.Info(ctx, "apple_verify_receipt", map[string]any{
		"userId":      userId,
		"status":      result.Status,
		"environment": result.Environment,
	})

	return &response.AppleVerifyReceipt{
		Status:        result.Status,
		Environment:   result.Environment,
		LatestReceipt: result.LatestReceipt,
		ReceiptInfo:   receiptInfo,
	}, nil
}
