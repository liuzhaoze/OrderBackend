package dto

import (
	"common/protobuf/stockpb"
	"order/domain"
)

type ItemWithQuantityConverter struct{}

func NewItemWithQuantityConverter() *ItemWithQuantityConverter {
	return &ItemWithQuantityConverter{}
}

func (c *ItemWithQuantityConverter) ToStockGrpc(item *domain.ItemWithQuantity) *stockpb.ItemWithQuantity {
	return &stockpb.ItemWithQuantity{ItemID: item.ItemID, Quantity: item.Quantity}
}
