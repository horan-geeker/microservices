package callback

import (
	"context"
	"fmt"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"microservices/entity/config"
	"microservices/pkg/log"
)

func (l *logic) HandleAlipayNotify(ctx context.Context, payload map[string]string) error {
	tradeStatus := payload["trade_status"]
	if tradeStatus != "TRADE_SUCCESS" && tradeStatus != "TRADE_FINISHED" {
		log.Info(ctx, "alipay_notify_ignore", map[string]any{"tradeStatus": tradeStatus})
		return nil
	}

	outTradeNo := payload["out_trade_no"]
	if outTradeNo == "" {
		return fmt.Errorf("异步通知缺少 out_trade_no")
	}
	order, err := l.model.Order().GetByOutTradeNo(ctx, outTradeNo)
	if err != nil {
		return err
	}
	log.Info(ctx, "alipay_notify_processed", map[string]any{
		"orderId":    order.ID,
		"outTradeNo": outTradeNo,
	})
	// todo 购买成功激活后续业务
	return nil
}

func (l *logic) HandleAlipayCallback(ctx context.Context, params map[string]string) error {
	log.Info(ctx, "alipay_callback_logic", map[string]any{
		"params": params,
	})
	return nil
}

func (l *logic) VerifyAlipayNotifySign(ctx context.Context, bm gopay.BodyMap) error {
	alipayOption := config.NewAlipayOptions()
	var (
		ok  bool
		err error
	)
	if alipayOption.AlipayPublicKey != "" {
		ok, err = alipay.VerifySign(alipayOption.AlipayPublicKey, bm)
	} else if alipayOption.AlipayPublicCertPath != "" {
		ok, err = alipay.VerifySignWithCert(alipayOption.AlipayPublicCertPath, bm)
	} else {
		return fmt.Errorf("未配置支付宝公钥或公钥证书")
	}

	if err != nil {
		log.Error(ctx, "alipay_notify_verify_error", err, map[string]any{"bodyMap": bm})
		return err
	}

	if !ok {
		return fmt.Errorf("支付宝验签失败")
	}

	return nil
}
