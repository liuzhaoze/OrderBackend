package domain

import (
	"common/consts"
	"errors"
)

type Order struct {
	OrderID     string
	CustomerID  string
	Items       []*Item
	Status      consts.OrderStatus
	PaymentLink string
}

func NewOrder(orderID string, customerID string, items []*Item, status consts.OrderStatus, paymentLink string) (*Order, error) {
	if customerID == "" {
		return nil, errors.New("customerID is empty")
	}
	if items == nil {
		return nil, errors.New("items is empty")
	}
	return &Order{OrderID: orderID, CustomerID: customerID, Items: items, Status: status, PaymentLink: paymentLink}, nil
}

func NewPendingOrder(customerID string, items []*Item) (*Order, error) {
	return NewOrder("", customerID, items, consts.OrderStatusPending, "")
}
