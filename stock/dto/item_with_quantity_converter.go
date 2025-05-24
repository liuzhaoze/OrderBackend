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

func (c *ItemWithQuantityConverter) FromStockGrpcBatch(items []*stockpb.ItemWithQuantity) []*domain.ItemWithQuantity {
	result := make([]*domain.ItemWithQuantity, len(items))
	for i, item := range items {
		result[i] = c.FromStockGrpc(item)
	}
	return result
}
