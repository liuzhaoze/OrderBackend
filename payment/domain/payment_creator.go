package domain

import "context"

type PaymentCreator interface {
	CreatePaymentLink(ctx context.Context, order *Order) (string, error)
}
