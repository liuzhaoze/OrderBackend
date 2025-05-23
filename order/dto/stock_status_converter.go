package dto

import (
	"common/consts"
	"common/protobuf/stockpb"
	"errors"
)

type StockStatusConverter struct{}

func NewStockStatusConverter() *StockStatusConverter {
	return &StockStatusConverter{}
}

func (c *StockStatusConverter) FromStockGrpc(stockStatus stockpb.StockStatus) (consts.StockStatus, error) {
	switch stockStatus {
	case stockpb.StockStatus_Sufficient:
		return consts.StockStatusSufficient, nil
	case stockpb.StockStatus_Insufficient:
		return consts.StockStatusInsufficient, nil
	default:
		return consts.StockStatusUnknown, errors.New("invalid stock status converting from gRPC")
	}
}
