package main

import (
	"common/consts"
	"common/protobuf/stockpb"
	"context"
	"errors"
	"stock/application"
	"stock/application/command"
	"stock/application/query"
	"stock/domain"
	"stock/dto"
)

type GrpcHandler struct {
	app *application.Application
}

func NewGrpcHandler(app *application.Application) *GrpcHandler {
	return &GrpcHandler{app: app}
}

func (g *GrpcHandler) CheckAndFetchItems(ctx context.Context, request *stockpb.CheckAndFetchItemsRequest) (*stockpb.CheckAndFetchItemsResponse, error) {
	requestItems := make([]*domain.ItemWithQuantity, len(request.Items))
	for i, item := range request.Items {
		requestItems[i] = dto.NewItemWithQuantityConverter().FromStockGrpc(item)
	}
	checkResult, checkErr := g.app.Queries.CheckItems.Handle(ctx, query.CheckItemsQuery{
		Items: requestItems,
	})
	if checkErr != nil {
		return nil, checkErr
	}

	grpcItems := make([]*stockpb.Item, len(checkResult.Items))
	for i, item := range checkResult.Items {
		grpcItems[i] = dto.NewItemConverter().ToStockGrpc(item)
	}

	switch checkResult.StockStatus {
	case consts.StockStatusInsufficient:
		statusCode, err := dto.NewStockStatusConverter().ToStockGrpc(checkResult.StockStatus)
		if err != nil {
			return nil, err
		}

		return &stockpb.CheckAndFetchItemsResponse{
			StatusCode: statusCode,
			Items:      grpcItems,
		}, nil
	case consts.StockStatusSufficient:
		fetchResult, fetchErr := g.app.Commands.FetchItems.Handle(ctx, command.FetchItemsCommand{
			Items: requestItems,
		})
		if fetchErr != nil {
			return nil, fetchErr
		}

		grpcItems = make([]*stockpb.Item, len(fetchResult.Items))
		for i, item := range fetchResult.Items {
			grpcItems[i] = dto.NewItemConverter().ToStockGrpc(item)
		}

		statusCode, err := dto.NewStockStatusConverter().ToStockGrpc(checkResult.StockStatus)
		if err != nil {
			return nil, err
		}

		return &stockpb.CheckAndFetchItemsResponse{
			StatusCode: statusCode,
			Items:      grpcItems,
		}, nil
	default:
		return &stockpb.CheckAndFetchItemsResponse{StatusCode: stockpb.StockStatus_Unknown, Items: make([]*stockpb.Item, 0)}, errors.New("unknown stock status in gRPC: CheckAndFetchItems")
	}
}
