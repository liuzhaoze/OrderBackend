package dto

import (
	"common/protobuf/orderpb"
	"process/domain"
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

func (c *OrderConverter) FromOrderGrpc(order *orderpb.Order) *domain.Order {
	status, _ := NewOrderStatusConverter().FromOrderGrpc(order.Status)
	return &domain.Order{
		OrderID:     order.OrderID,
		CustomerID:  order.CustomerID,
		Items:       NewItemConverter().FromOrderGrpcBatch(order.Items),
		Status:      status,
		PaymentLink: order.PaymentLink,
	}
}
