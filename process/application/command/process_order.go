package command

import (
	"common/consts"
	"common/cqrs"
	"common/metrics"
	"common/protobuf/orderpb"
	"common/tracing"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"process/domain"
	"process/dto"
	"time"
)

type ProcessOrderCommand struct {
	Order *domain.Order
}

type ProcessOrderResult struct {
	Order *domain.Order
}

type ProcessOrderHandler cqrs.CommandHandler[ProcessOrderCommand, ProcessOrderResult]

type processOrder struct {
	orderGrpc orderpb.OrderServiceClient
}

func (p processOrder) Handle(ctx context.Context, command ProcessOrderCommand) (ProcessOrderResult, error) {
	ctx, span := tracing.StartSpan(ctx, "Process/Application/Command: process order")
	defer span.End()

	order := command.Order
	if order.Status != consts.OrderStatusPaid {
		return ProcessOrderResult{Order: nil}, fmt.Errorf("order %s is not paid", command.Order.OrderID)
	}

	time.Sleep(5 * time.Second) // 模拟处理订单
	order.Status = consts.OrderStatusFinished
	span.AddEvent("订单处理完成")

	result, err := p.orderGrpc.UpdateOrder(ctx, &orderpb.UpdateOrderRequest{
		UpdateOptions: orderpb.UpdateOption_Status,
		Order:         dto.NewOrderConverter().ToOrderGrpc(order),
	})
	if err != nil {
		return ProcessOrderResult{Order: nil}, err
	}

	return ProcessOrderResult{Order: dto.NewOrderConverter().FromOrderGrpc(result.Order)}, nil
}

func NewProcessOrderHandler(orderGrpc orderpb.OrderServiceClient,
	logger *logrus.Logger,
	metricsClient metrics.Client,
) ProcessOrderHandler {
	return cqrs.ApplyCommandDecorator[ProcessOrderCommand, ProcessOrderResult](
		processOrder{orderGrpc: orderGrpc},
		logger,
		metricsClient,
	)
}
