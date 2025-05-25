package stripe

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"payment/domain"
)

type Stripe struct {
	paymentBaseURL string
}

func NewStripe(apiKey, paymentBaseURL string) *Stripe {
	stripe.Key = apiKey // 为 Stripe SDK 设置全局 Stripe API key
	return &Stripe{paymentBaseURL: paymentBaseURL}
}

func (s *Stripe) CreatePaymentLink(ctx context.Context, order *domain.Order) (string, error) {
	items := make([]*stripe.CheckoutSessionLineItemParams, len(order.Items))
	for i, item := range order.Items {
		items[i] = &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(item.PriceID),
			Quantity: stripe.Int64(item.Quantity),
		}
	}

	marshalled, err := json.Marshal(order.Items)
	if err != nil {
		return "", err
	}

	metadata := map[string]string{
		"order_id":    order.OrderID,
		"customer_id": order.CustomerID,
		"status":      string(order.Status),
		"items":       string(marshalled),
	}

	params := &stripe.CheckoutSessionParams{
		LineItems:  items,
		Metadata:   metadata,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(fmt.Sprintf("%s/payment?customer-id=%s&order-id=%s", s.paymentBaseURL, order.CustomerID, order.OrderID)),
	}

	result, err := session.New(params)
	if err != nil {
		return "", err
	}

	return result.URL, nil
}
