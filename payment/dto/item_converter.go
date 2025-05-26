package dto

import (
	"common/protobuf/orderpb"
	"payment/domain"
)

type ItemConverter struct{}

func NewItemConverter() *ItemConverter {
	return &ItemConverter{}
}

func (c *ItemConverter) ToOrderGrpc(item *domain.Item) *orderpb.Item {
	return &orderpb.Item{ItemID: item.ItemID, Name: item.Name, Quantity: item.Quantity, PriceID: item.PriceID}
}

func (c *ItemConverter) ToOrderGrpcBatch(items []*domain.Item) []*orderpb.Item {
	result := make([]*orderpb.Item, len(items))
	for i, item := range items {
		result[i] = c.ToOrderGrpc(item)
	}
	return result
}
