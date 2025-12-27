package service

import (
	"context"
	"fmt"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/sirupsen/logrus"
	"microservices/entity/config"
	"microservices/pkg/log"
)

type Alipay interface {
	CreateAppPrepay(ctx context.Context, outTradeNo string, amount int64, description string) (*AppPrepayResult, error)
	QueryOrder(ctx context.Context, orderId string) (*alipay.TradeQueryResponse, error)
}

type AppPrepayResult struct {
	OutTradeNo  string
	OrderString string
}

type ali struct {
	client *alipay.Client
	opt    *config.AlipayOptions
}

func newAlipay(opt *config.AlipayOptions) Alipay {
	client, err := alipay.NewClient(opt.AppId, opt.PrivateKey, opt.IsProduction)
	if err != nil {
		log.Error(context.Background(), "pay", err, map[string]any{
			"appId":        opt.AppId,
			"privateKey":   opt.PrivateKey,
			"isProduction": opt.IsProduction,
		})
		panic("Failed to create Alipay legacy client")
	}
	if !opt.IsProduction {
		// 打开Debug开关，输出日志，默认关闭
		client.DebugSwitch = gopay.DebugOn
	}
	// 设置支付宝请求 公共参数
	// 注意：具体设置哪些参数，根据不同的方法而不同，此处列举出所有设置参数
	client.SetLocation(alipay.LocationShanghai). // 设置时区，不设置或出错均为默认服务器时间
							SetCharset(alipay.UTF8).     // 设置字符编码，不设置默认 utf-8
							SetSignType(alipay.RSA2).    // 设置签名类型，不设置默认 RSA2
							SetReturnUrl(opt.ReturnUrl). // 设置返回URL
							SetNotifyUrl(opt.NotifyUrl)  // 设置异步通知URL
	if opt.AppPublicCertPath != "" && opt.AlipayRootCertPath != "" && opt.AlipayPublicCertPath != "" {
		if err = client.SetCertSnByPath(opt.AppPublicCertPath, opt.AlipayRootCertPath, opt.AlipayPublicCertPath); err != nil {
			logrus.WithError(err).Error("Failed to set Alipay certificates for legacy client")
		}
	}
	//if opt.AlipayPublicKey != "" {
	//	if err = client.SetAliPayPublicCertSN(opt.AlipayPublicKey); err != nil {
	//		logrus.WithError(err).Error("Failed to set Alipay public key for legacy client")
	//	}
	//}

	return &ali{
		client: client,
		opt:    opt,
	}
}

func (a *ali) CreateAppPrepay(ctx context.Context, outTradeNo string, amount int64, subject string) (*AppPrepayResult, error) {
	totalAmount := fmt.Sprintf("%.2f", float64(amount)/100.0)

	bm := make(gopay.BodyMap)
	bm.Set("subject", subject).
		Set("out_trade_no", outTradeNo).
		Set("total_amount", totalAmount).
		Set("product_code", "QUICK_MSECURITY_PAY")

	orderString, err := a.client.TradeAppPay(ctx, bm)
	if err != nil {
		log.Error(ctx, "pay", err, nil)
		return nil, fmt.Errorf("创建支付宝预支付订单失败: %w", err)
	}

	if orderString == "" {
		return nil, fmt.Errorf("支付宝未返回有效的预支付信息")
	}

	return &AppPrepayResult{
		OutTradeNo:  outTradeNo,
		OrderString: orderString,
	}, nil
}

func (a *ali) QueryOrder(ctx context.Context, outTradeNo string) (*alipay.TradeQueryResponse, error) {
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", outTradeNo)

	aliRsp, err := a.client.TradeQuery(ctx, bm)
	if err != nil {
		logrus.WithError(err).Error("Failed to query Alipay order")
		return nil, fmt.Errorf("查询支付宝订单失败: %w", err)
	}

	//if aliRsp.StatusCode != alipayv3.Success {
	//	return nil, fmt.Errorf("支付宝查询订单失败: %s", gconv.String(aliRsp.ErrResponse))
	//}

	return aliRsp, nil
}
