package main

import (
	"common/protobuf/orderpb"
	"context"
	"fmt"
)

type GrpcHandler struct {
}

func NewGrpcHandler() *GrpcHandler {
	return &GrpcHandler{}
}

func (g *GrpcHandler) UpdateOrder(ctx context.Context, request *orderpb.UpdateOrderRequest) (*orderpb.UpdateOrderResponse, error) {
	// TODO: implement real logic
	fmt.Println(request.Order)
	return &orderpb.UpdateOrderResponse{Order: request.Order}, nil
}
