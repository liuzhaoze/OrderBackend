package domain

import (
	"context"
	"fmt"
)

type OrderRepository interface {
	Create(ctx context.Context, order *Order) (*Order, error)
	Get(ctx context.Context, orderID string, customerID string) (*Order, error)
	Update(ctx context.Context, order *Order, updateFunc func(context.Context, *Order) (*Order, error)) (*Order, error)
}

type NotFoundError struct {
	OrderID    string
	CustomerID string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Not fount order: %s for customer: %s", e.OrderID, e.CustomerID)
}
