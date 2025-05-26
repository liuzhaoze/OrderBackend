package dto

import (
	"common/consts"
	"common/protobuf/orderpb"
	"errors"
)

type OrderStatusConverter struct{}

func NewOrderStatusConverter() *OrderStatusConverter {
	return &OrderStatusConverter{}
}

func (c *OrderStatusConverter) FromOrderGrpc(orderStatus orderpb.OrderStatus) (consts.OrderStatus, error) {
	switch orderStatus {
	case orderpb.OrderStatus_Pending:
		return consts.OrderStatusPending, nil
	case orderpb.OrderStatus_WaitingForPayment:
		return consts.OrderStatusWaitingForPayment, nil
	case orderpb.OrderStatus_Paid:
		return consts.OrderStatusPaid, nil
	case orderpb.OrderStatus_Finished:
		return consts.OrderStatusFinished, nil
	default:
		return consts.OrderStatusUnknown, errors.New("invalid order status converting from gRPC")
	}
}

func (c *OrderStatusConverter) ToOrderGrpc(orderStatus consts.OrderStatus) (orderpb.OrderStatus, error) {
	switch orderStatus {
	case consts.OrderStatusPending:
		return orderpb.OrderStatus_Pending, nil
	case consts.OrderStatusWaitingForPayment:
		return orderpb.OrderStatus_WaitingForPayment, nil
	case consts.OrderStatusPaid:
		return orderpb.OrderStatus_Paid, nil
	case consts.OrderStatusFinished:
		return orderpb.OrderStatus_Finished, nil
	default:
		return orderpb.OrderStatus_Unknown, errors.New("invalid order status converting to gRPC")
	}
}
