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

func (c *ItemWithQuantityConverter) ToStockGrpcBatch(items []*domain.ItemWithQuantity) []*stockpb.ItemWithQuantity {
	result := make([]*stockpb.ItemWithQuantity, len(items))
	for i, item := range items {
		result[i] = c.ToStockGrpc(item)
	}
	return result
}

func (c *ItemWithQuantityConverter) FromHttp(item ports.ItemWithQuantity) *domain.ItemWithQuantity {
	return &domain.ItemWithQuantity{ItemID: item.ItemID, Quantity: item.Quantity}
}

func (c *ItemWithQuantityConverter) FromHttpBatch(items []ports.ItemWithQuantity) []*domain.ItemWithQuantity {
	result := make([]*domain.ItemWithQuantity, len(items))
	for i, item := range items {
		result[i] = c.FromHttp(item)
	}
	return result
}
