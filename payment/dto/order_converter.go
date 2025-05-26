package dto

import (
	"common/protobuf/orderpb"
	"payment/domain"
)

type OrderConverter struct{}

func NewOrderConverter() *OrderConverter {
	return &OrderConverter{}
}

func (c *OrderConverter) ToOrderGrpc(order *domain.Order) *orderpb.Order {
	status, _ := NewOrderStatusConverter().ToOrderGrpc(order.Status)
	return &orderpb.Order{
		OrderID:     order.OrderID,
		CustomerID:  order.CustomerID,
		Items:       NewItemConverter().ToOrderGrpcBatch(order.Items),
		Status:      status,
		PaymentLink: order.PaymentLink,
	}
}
