package domain

import (
	"common/consts"
	"errors"
	"fmt"
	"slices"
)

type Order struct {
	OrderID     string
	CustomerID  string
	Items       []*Item
	Status      consts.OrderStatus
	PaymentLink string
}

var validStatusTransition = map[consts.OrderStatus][]consts.OrderStatus{
	consts.OrderStatusUnknown:           {},
	consts.OrderStatusPending:           {consts.OrderStatusWaitingForPayment},
	consts.OrderStatusWaitingForPayment: {consts.OrderStatusPaid},
	consts.OrderStatusPaid:              {consts.OrderStatusFinished},
	consts.OrderStatusFinished:          {},
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

func (o *Order) UpdateStatus(status consts.OrderStatus) error {
	if slices.Contains(validStatusTransition[o.Status], status) {
		o.Status = status
		return nil
	} else {
		return fmt.Errorf("cannot update order status from %s to %s", o.Status, status)
	}
}

func (o *Order) UpdatePaymentLink(paymentLink string) error {
	o.PaymentLink = paymentLink
	return nil
}
