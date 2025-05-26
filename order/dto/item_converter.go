package dto

import (
	"common/protobuf/orderpb"
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

func (c *ItemConverter) FromOrderGrpc(item *orderpb.Item) *domain.Item {
	return &domain.Item{ItemID: item.ItemID, Name: item.Name, Quantity: item.Quantity, PriceID: item.PriceID}
}

func (c *ItemConverter) FromStockGrpcBatch(items []*stockpb.Item) []*domain.Item {
	result := make([]*domain.Item, len(items))
	for i, item := range items {
		result[i] = c.FromStockGrpc(item)
	}
	return result
}

func (c *ItemConverter) FromOrderGrpcBatch(items []*orderpb.Item) []*domain.Item {
	result := make([]*domain.Item, len(items))
	for i, item := range items {
		result[i] = c.FromOrderGrpc(item)
	}
	return result
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
