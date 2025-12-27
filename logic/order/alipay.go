package order

import (
	"context"
	"fmt"
	"microservices/entity/config"
	"microservices/entity/consts"
	entity "microservices/entity/model"
	"microservices/pkg/util"
	"time"
)

// CreateAliPayPrepayOrder price 分为单位
func (p *logic) CreateAliPayPrepayOrder(ctx context.Context, userId, price int, platform, channel, clientIP,
	description string) (*entity.Order, string, error) {
	// 根据支付渠道生成不同格式的 outTradeNo
	var outTradeNo string
	// 其他支付渠道使用时间戳格式
	outTradeNo = fmt.Sprintf("%s%d%010d", time.Now().Format(consts.TimeFormatSeq), util.RandomN(3), userId)

	var mchId, appId, tradeType, orderString string
	if channel == consts.TradeTypeAlipay {
		alipayConfig := config.NewAlipayOptions()
		mchId = alipayConfig.AppId // 支付宝使用 AppId 作为商户标识
		appId = alipayConfig.AppId
		tradeType = consts.TradeTypeAlipay
	} else if channel == consts.TradeTypeApple {
		tradeType = consts.TradeTypeApple
	}

	expiredAt := time.Now().Add(30 * time.Minute)
	// 创建订单记录
	order := &entity.Order{
		MchId:       mchId,
		AppId:       appId,
		OutTradeNo:  outTradeNo,
		UserId:      userId,
		TotalAmount: uint64(price),
		Currency:    "",
		Subject:     "",
		Status:      consts.OrderStatusPending,
		TradeType:   tradeType, // 支付类型
		Platform:    platform,
		ClientIp:    &clientIP,
		ExpireAt:    &expiredAt,
	}

	if err := p.model.Order().Create(ctx, order); err != nil {
		return nil, "", fmt.Errorf("创建订单记录失败: %w", err)
	}

	if channel == consts.TradeTypeAlipay {
		// 调用支付宝支付
		result, err := p.srv.Alipay().CreateAppPrepay(ctx, outTradeNo, int64(price), description)
		if err != nil {
			// 如果支付宝支付失败，更新订单状态
			_ = p.model.Order().Update(ctx, order.ID, map[string]interface{}{
				"status": consts.OrderStatusCreateError,
			})
			return nil, "", err
		}
		orderString = result.OrderString
	}

	return order, orderString, nil
}
