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
