package dto

import (
	"common/protobuf/stockpb"
	"order/domain"
	"order/ports"
)

type ItemWithQuantityConverter struct{}

func NewItemWithQuantityConverter() *ItemWithQuantityConverter {
	return &ItemWithQuantityConverter{}
}

func (c *ItemWithQuantityConverter) ToStockGrpc(item *domain.ItemWithQuantity) *stockpb.ItemWithQuantity {
	return &stockpb.ItemWithQuantity{ItemID: item.ItemID, Quantity: item.Quantity}
}

func (c *ItemWithQuantityConverter) FromHttp(item ports.ItemWithQuantity) *domain.ItemWithQuantity {
	return &domain.ItemWithQuantity{ItemID: item.ItemID, Quantity: item.Quantity}
}
