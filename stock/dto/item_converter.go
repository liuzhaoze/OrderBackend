package dto

import (
	"common/protobuf/stockpb"
	"stock/domain"
)

type ItemConverter struct{}

func NewItemConverter() *ItemConverter {
	return &ItemConverter{}
}

func (c *ItemConverter) ToStockGrpc(item *domain.Item) *stockpb.Item {
	return &stockpb.Item{ItemID: item.ItemID, Name: item.Name, Quantity: item.Quantity, PriceID: item.PriceID}
}

func (c *ItemConverter) ToStockGrpcBatch(items []*domain.Item) []*stockpb.Item {
	result := make([]*stockpb.Item, len(items))
	for i, item := range items {
		result[i] = c.ToStockGrpc(item)
	}
	return result
}
