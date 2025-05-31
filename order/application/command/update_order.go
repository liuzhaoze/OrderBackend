package command

import (
	"common/cqrs"
	"common/metrics"
	"common/tracing"
	"context"
	"github.com/sirupsen/logrus"
	"order/domain"
)

type UpdateOrderCommand struct {
	Order      *domain.Order
	UpdateFunc func(context.Context, *domain.Order) (*domain.Order, error)
}

type UpdateOrderResult struct {
	Order *domain.Order
}

type UpdateOrderHandler cqrs.CommandHandler[UpdateOrderCommand, UpdateOrderResult]

type updateOrder struct {
	orderRepo domain.OrderRepository
}

func (u updateOrder) Handle(ctx context.Context, command UpdateOrderCommand) (UpdateOrderResult, error) {
	ctx, span := tracing.StartSpan(ctx, "Order/Application/Command: update order")
	defer span.End()

	if updatedOrder, err := u.orderRepo.Update(ctx, command.Order, command.UpdateFunc); err != nil {
		return UpdateOrderResult{Order: nil}, err
	} else {
		return UpdateOrderResult{Order: updatedOrder}, nil
	}
}

func NewUpdateOrderHandler(orderRepo domain.OrderRepository,
	logger *logrus.Logger,
	metricsClient metrics.Client,
) UpdateOrderHandler {
	return cqrs.ApplyCommandDecorator[UpdateOrderCommand, UpdateOrderResult](
		updateOrder{orderRepo: orderRepo},
		logger,
		metricsClient,
	)
}
