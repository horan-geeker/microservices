package service

import (
	"context"
	"fmt"
	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/checkout/session"
	"github.com/stripe/stripe-go/v84/customer"
	"github.com/stripe/stripe-go/v84/webhook"
	"microservices/entity/config"
	"microservices/entity/model"
	"microservices/pkg/log"
	"strconv"
)

type Stripe interface {
	CreateCheckoutSession(ctx context.Context, userID *model.User, priceID string, orderId int) (string, error)
	ConstructEvent(payload []byte, header string) (stripe.Event, error)
}

type stripeService struct {
	opt *config.StripeOptions
}

func newStripe(opt *config.StripeOptions) Stripe {
	stripe.Key = opt.ApiKey
	return &stripeService{
		opt: opt,
	}
}

func (s *stripeService) CreateCheckoutSession(ctx context.Context, user *model.User, priceID string, orderId int) (string, error) {
	customerParams := &stripe.CustomerParams{
		Email: stripe.String(user.Email),
		Name:  stripe.String(user.Name), // ✔ 只用于 Stripe 记录
		Metadata: map[string]string{
			"userId": strconv.Itoa(user.ID),
		},
	}
	customer, err := customer.New(customerParams)
	if err != nil {
		return "", err
	}
	params := &stripe.CheckoutSessionParams{
		// 关键设置 1：将账单地址收集设为 Auto（或不设置），不要设为 Required，否则姓名必填
		BillingAddressCollection: stripe.String("auto"),
		// 关键设置 2：使用 PaymentMethodOptions 显式预填（v84 语法）
		PaymentMethodOptions: &stripe.CheckoutSessionPaymentMethodOptionsParams{
			WeChatPay: &stripe.CheckoutSessionPaymentMethodOptionsWeChatPayParams{
				Client: stripe.String("web"),
			},
		},
		// 3、传入 customer
		Customer: stripe.String(customer.ID),
		// --- 关键代码：通过 Metadata 传递 ---
		Metadata: map[string]string{
			"priceId": priceID,               // 存入 PriceID
			"userId":  strconv.Itoa(user.ID), // 顺便存入你的用户 ID，方便回调时给用户发货
		},
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			Metadata: map[string]string{
				"priceId": priceID,
				"userId":  strconv.Itoa(user.ID),
				"orderId": fmt.Sprintf("%d", orderId),
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		ClientReferenceID: stripe.String(fmt.Sprintf("%d", orderId)),
		SuccessURL:        stripe.String("https://" + s.opt.Domain + "/recharge?success=true"),
		CancelURL:         stripe.String("https://" + s.opt.Domain + "/recharge?canceled=true"),
	}

	sess, err := session.New(params)
	if err != nil {
		log.Error(ctx, "stripe_checkout_session_create_failed", err, nil)
		return "", err
	}

	return sess.URL, nil
}

func (s *stripeService) ConstructEvent(payload []byte, header string) (stripe.Event, error) {
	return webhook.ConstructEvent(payload, header, s.opt.WebhookSecret)
}
