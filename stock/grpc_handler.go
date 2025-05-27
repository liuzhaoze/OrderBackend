package main

import (
	"common/consts"
	"common/protobuf/stockpb"
	"common/tracing"
	"context"
	"errors"
	"stock/application"
	"stock/application/command"
	"stock/application/query"
	"stock/dto"
)

type GrpcHandler struct {
	app *application.Application
}

func NewGrpcHandler(app *application.Application) *GrpcHandler {
	return &GrpcHandler{app: app}
}

func (g *GrpcHandler) CheckAndFetchItems(ctx context.Context, request *stockpb.CheckAndFetchItemsRequest) (*stockpb.CheckAndFetchItemsResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "Stock/gRPC: 校验库存和扣减库存")
	defer span.End()

	requestItems := dto.NewItemWithQuantityConverter().FromStockGrpcBatch(request.Items)

	checkResult, checkErr := g.app.Queries.CheckItems.Handle(ctx, query.CheckItemsQuery{
		Items: requestItems,
	})
	if checkErr != nil {
		return nil, checkErr
	}

	switch checkResult.StockStatus {
	case consts.StockStatusInsufficient:
		span.AddEvent("库存不足")
		statusCode, err := dto.NewStockStatusConverter().ToStockGrpc(checkResult.StockStatus)
		if err != nil {
			return nil, err
		}

		return &stockpb.CheckAndFetchItemsResponse{
			StatusCode: statusCode,
			Items:      dto.NewItemConverter().ToStockGrpcBatch(checkResult.Items),
		}, nil
	case consts.StockStatusSufficient:
		span.AddEvent("库存充足")
		fetchResult, fetchErr := g.app.Commands.FetchItems.Handle(ctx, command.FetchItemsCommand{
			Items: requestItems,
		})
		if fetchErr != nil {
			return nil, fetchErr
		}
		span.AddEvent("已扣减库存")

		statusCode, err := dto.NewStockStatusConverter().ToStockGrpc(checkResult.StockStatus)
		if err != nil {
			return nil, err
		}

		return &stockpb.CheckAndFetchItemsResponse{
			StatusCode: statusCode,
			Items:      dto.NewItemConverter().ToStockGrpcBatch(fetchResult.Items),
		}, nil
	default:
		return &stockpb.CheckAndFetchItemsResponse{StatusCode: stockpb.StockStatus_Unknown, Items: make([]*stockpb.Item, 0)}, errors.New("unknown stock status in gRPC: CheckAndFetchItems")
	}
}
