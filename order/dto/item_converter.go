package dto

import (
	"common/protobuf/stockpb"
	"order/domain"
)

type ItemConverter struct{}

func NewItemConverter() *ItemConverter {
	return &ItemConverter{}
}

func (c *ItemConverter) FromStockGrpc(item *stockpb.Item) *domain.Item {
	return &domain.Item{ItemID: item.ItemID, Name: item.Name, Quantity: item.Quantity, PriceID: item.PriceID}
}

func (c *ItemConverter) FromStockGrpcBatch(items []*stockpb.Item) []*domain.Item {
	result := make([]*domain.Item, len(items))
	for i, item := range items {
		result[i] = c.FromStockGrpc(item)
	}
	return result
}
