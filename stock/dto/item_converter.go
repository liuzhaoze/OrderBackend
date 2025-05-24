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
