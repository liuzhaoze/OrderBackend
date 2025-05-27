package main

import (
	"common/protobuf/orderpb"
	"common/tracing"
	"context"
	"order/application"
	"order/application/command"
	"order/domain"
	"order/dto"
)

type GrpcHandler struct {
	app *application.Application
}

func NewGrpcHandler(app *application.Application) *GrpcHandler {
	return &GrpcHandler{app: app}
}

func (g *GrpcHandler) UpdateOrder(ctx context.Context, request *orderpb.UpdateOrderRequest) (*orderpb.UpdateOrderResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "Order/gRPC/update order: 更新订单")
	defer span.End()

	order := dto.NewOrderConverter().FromOrderGrpc(request.Order)
	result, err := g.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrderCommand{
		Order: order,
		UpdateFunc: func(c context.Context, o *domain.Order) (*domain.Order, error) {
			if (request.UpdateOptions & orderpb.UpdateOption_Status) != 0 {
				if err := o.UpdateStatus(order.Status); err != nil {
					return nil, err
				}
			}
			if (request.UpdateOptions & orderpb.UpdateOption_PaymentLink) != 0 {
				if err := o.UpdatePaymentLink(order.PaymentLink); err != nil {
					return nil, err
				}
			}
			return o, nil
		},
	})
	if err != nil {
		return nil, err
	}
	return &orderpb.UpdateOrderResponse{Order: dto.NewOrderConverter().ToOrderGrpc(result.Order)}, nil
}
