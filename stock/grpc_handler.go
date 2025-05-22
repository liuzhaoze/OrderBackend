package main

import (
	"common/protobuf/stockpb"
	"context"
	"fmt"
)

type GrpcHandler struct {
}

func NewGrpcHandler() *GrpcHandler {
	return &GrpcHandler{}
}

func (g *GrpcHandler) CheckAndFetchItems(ctx context.Context, request *stockpb.CheckAndFetchItemsRequest) (*stockpb.CheckAndFetchItemsResponse, error) {
	// TODO: implement real logic
	for _, item := range request.Items {
		fmt.Println((*item).ItemID, (*item).Quantity)
	}
	return &stockpb.CheckAndFetchItemsResponse{StatusCode: stockpb.StockStatus_Insufficient, Items: make([]*stockpb.Item, 0)}, nil
}
