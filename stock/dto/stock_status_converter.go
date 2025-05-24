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

func (c *StockStatusConverter) ToStockGrpc(stockStatus consts.StockStatus) (stockpb.StockStatus, error) {
	switch stockStatus {
	case consts.StockStatusSufficient:
		return stockpb.StockStatus_Sufficient, nil
	case consts.StockStatusInsufficient:
		return stockpb.StockStatus_Insufficient, nil
	default:
		return stockpb.StockStatus_Unknown, errors.New("invalid stock status converting to gRPC")
	}
}
