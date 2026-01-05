package callback

import (
	"context"
	"encoding/json"
	"fmt"
	"microservices/entity/consts"
	"microservices/pkg/log"
	"strconv"
	"time"

	"github.com/stripe/stripe-go/v84"
	"gorm.io/gorm"
)

func (l *logic) HandleStripeCallback(ctx context.Context, payload []byte, header string) error {
	event, err := l.srv.Stripe().ConstructEvent(payload, header)
	if err != nil {
		log.Error(ctx, "stripe_webhook_verify_failed", err, nil)
		return fmt.Errorf("Webhook signature verification failed")
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			log.Error(ctx, "stripe_webhook_parse_failed", err, nil)
			return err
		}

		orderId, err := strconv.Atoi(session.ClientReferenceID)
		if err != nil {
			log.Error(ctx, "stripe_callback_invalid_client_ref", err,
				map[string]any{"clientRef": session.ClientReferenceID})
			return fmt.Errorf("Invalid ClientReferenceID")
		}

		priceId := session.Metadata["priceId"]
		userIdStr := session.Metadata["userId"]
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			log.Error(ctx, "stripe_callback_order_metadata_err", err, map[string]any{"metadata": session.Metadata})
			return err
		}
		order, err := l.model.Order().GetByIdAndUserId(ctx, orderId, userId)
		if err != nil {
			log.Error(ctx, "stripe_callback_order_not_found", err, map[string]any{"orderId": orderId})
			return err
		}

		// Update Order Status to 1 (Paid)
		err = l.model.Order().Update(ctx, order.ID, map[string]interface{}{
			"status":  consts.OrderStatusPaid,
			"paid_at": time.Now(),
		})
		if err != nil {
			log.Error(ctx, "stripe_callback_update_order_failed", err, map[string]any{"orderId": order.ID})
			return err
		}

		goods, err := l.model.Goods().GetByPriceId(ctx, priceId)
		if err != nil {
			return err
		}
		err = l.model.User().Update(ctx, order.UserId, map[string]any{
			"credit": gorm.Expr("credit + ?", goods.Credit),
		})
		if err != nil {
			log.Error(ctx, "stripe_callback_update_credit_failed", err, map[string]any{"userId": order.UserId})
			return err
		}

		log.Info(ctx, "stripe_payment_success", map[string]any{"session": session})

	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			log.Error(ctx, "stripe_webhook_parse_failed", err, nil)
			return err
		}

		priceId := paymentIntent.Metadata["priceId"]
		userIdStr := paymentIntent.Metadata["userId"]
		orderIdStr := paymentIntent.Metadata["orderId"]

		if priceId == "" || userIdStr == "" || orderIdStr == "" {
			log.Error(ctx, "stripe_pi_missing_metadata", fmt.Errorf("missing metadata"), map[string]any{"metadata": paymentIntent.Metadata})
			return fmt.Errorf("missing metadata")
		}

		userId, _ := strconv.Atoi(userIdStr)
		orderId, _ := strconv.Atoi(orderIdStr)

		order, err := l.model.Order().GetByIdAndUserId(ctx, orderId, userId)
		if err != nil {
			log.Error(ctx, "stripe_pi_order_not_found", err, map[string]any{"orderId": orderId})
			return err
		}

		if order.Status == consts.OrderStatusPaid {
			log.Info(ctx, "stripe_pi_order_already_paid", map[string]any{"orderId": orderId})
			return nil
		}

		// Update Order Status to 1 (Paid)
		err = l.model.Order().Update(ctx, order.ID, map[string]interface{}{
			"status":  consts.OrderStatusPaid,
			"paid_at": time.Now(),
		})
		if err != nil {
			log.Error(ctx, "stripe_pi_update_order_failed", err, map[string]any{"orderId": order.ID})
			return err
		}

		goods, err := l.model.Goods().GetByPriceId(ctx, priceId)
		if err != nil {
			return err
		}
		err = l.model.User().Update(ctx, order.UserId, map[string]any{
			"credit": gorm.Expr("credit + ?", goods.Credit),
		})
		if err != nil {
			log.Error(ctx, "stripe_pi_update_credit_failed", err, map[string]any{"userId": order.UserId})
			return err
		}

		log.Info(ctx, "stripe_pi_payment_success", map[string]any{"pi": paymentIntent.ID})

	default:
		log.Info(ctx, "stripe_event_unhandled", map[string]any{"eventType": event.Type})
	}

	return nil
}
