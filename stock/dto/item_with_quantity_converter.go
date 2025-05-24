package dto

import (
	"common/protobuf/stockpb"
	"stock/domain"
)

type ItemWithQuantityConverter struct{}

func NewItemWithQuantityConverter() *ItemWithQuantityConverter {
	return &ItemWithQuantityConverter{}
}

func (c *ItemWithQuantityConverter) FromStockGrpc(item *stockpb.ItemWithQuantity) *domain.ItemWithQuantity {
	return &domain.ItemWithQuantity{ItemID: item.ItemID, Quantity: item.Quantity}
}
