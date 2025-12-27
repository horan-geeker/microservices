package callback

import (
	"context"
	"encoding/json"
	"fmt"
	"microservices/entity/consts"
	"microservices/entity/ecode"
	"microservices/entity/request"
	"microservices/pkg/log"
	"microservices/util"
)

func (l *logic) HandleAppleCallback(ctx context.Context, notification *request.AppleIAPNotification) error {
	// Verify and decode the signedPayload (JWS) using Apple service
	appleSvc := l.srv.Apple()
	payload, err := appleSvc.DecodeSignedPayload(notification.SignedPayload)
	if err != nil {
		log.Error(ctx, "apple_callback_decode_failed", err, map[string]any{
			"signedPayload": notification.SignedPayload,
		})
		return ecode.ErrParamInvalid.WithMessage(fmt.Sprintf("无效的 Apple 回调签名: %v", err))
	}

	log.Info(ctx, "apple_iap_callback_received", map[string]any{
		"notificationType": payload.NotificationType,
		"subtype":          payload.Subtype,
		"notificationUUID": payload.NotificationUUID,
		"version":          payload.Version,
		"signedDate":       payload.SignedDate,
	})
	if payload.NotificationType != consts.IOSNotificationTypeSubscribed && // 订阅初次购买
		payload.NotificationType != consts.IOSNotificationTypeOneTimeCharge && // 一次性消费
		payload.NotificationType != consts.IOSNotificationTypeDidRenew && // 订阅触发续订
		payload.NotificationType != consts.IOSNotificationTypeDidChangeRenewalPref { // 订阅变更
		return ecode.ErrParamInvalid.WithMessage("无法处理该 Apple 消息类型")
	}
	// Parse the data field which contains signedTransactionInfo and signedRenewalInfo
	var data struct {
		Environment           string `json:"environment"`
		AppAppleId            int64  `json:"appAppleId"`
		BundleId              string `json:"bundleId"`
		BundleVersion         string `json:"bundleVersion"`
		SignedTransactionInfo string `json:"signedTransactionInfo"`
		SignedRenewalInfo     string `json:"signedRenewalInfo"`
	}

	if len(payload.Data) > 0 {
		if err := json.Unmarshal(payload.Data, &data); err != nil {
			log.Error(ctx, "apple_callback_data_unmarshal_failed", err, map[string]any{
				"data": string(payload.Data),
			})
			return ecode.ErrParamInvalid.WithMessage("无法解析 Apple 回调数据")
		}

		log.Info(ctx, "apple_callback_data", map[string]any{
			"environment":   data.Environment,
			"appAppleId":    data.AppAppleId,
			"bundleId":      data.BundleId,
			"bundleVersion": data.BundleVersion,
		})

		// Decode the signed transaction info if present
		if data.SignedTransactionInfo != "" {
			transactionInfo, err := appleSvc.DecodeSignedTransaction(data.SignedTransactionInfo)
			if err != nil {
				log.Error(ctx, "apple_callback_transaction_decode_failed", err, nil)
			} else {
				txdata, _ := util.StructToMap(*transactionInfo)
				log.Info(ctx, "apple_transaction_info", txdata)
				order, err := l.model.Order().GetByOutTradeNo(ctx, transactionInfo.AppAccountToken)
				if err != nil {
					return err
				}
				// 核心开通会员逻辑
				if payload.NotificationType == consts.IOSNotificationTypeDidRenew ||
					payload.NotificationType == consts.IOSNotificationTypeDidChangeRenewalPref {
					// todo 订阅续订业务
					_ = order
				} else if payload.NotificationType == consts.IOSNotificationTypeSubscribed ||
					payload.NotificationType == consts.IOSNotificationTypeOneTimeCharge {
					// todo 单次购买业务
				}

			}
		}

		// Decode the signed renewal info if present
		if data.SignedRenewalInfo != "" {
			renewalInfo, err := appleSvc.DecodeSignedRenewalInfo(data.SignedRenewalInfo)
			if err != nil {
				log.Error(ctx, "apple_callback_renewal_decode_failed", err, nil)
			} else {
				log.Info(ctx, "apple_renewal_info", map[string]any{
					"originalTransactionId": renewalInfo.OriginalTransactionId,
					"autoRenewProductId":    renewalInfo.AutoRenewProductId,
					"autoRenewStatus":       renewalInfo.AutoRenewStatus,
					"environment":           renewalInfo.Environment,
				})
			}
		}
	}

	return nil
}
