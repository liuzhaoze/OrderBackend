package query

import (
	"common/cqrs"
	"context"
	"order/domain"
)

type GetOrderQuery struct {
	OrderID    string
	CustomerID string
}

type GetOrderResult struct {
	Order *domain.Order
}

type GetOrderHandler cqrs.QueryHandler[GetOrderQuery, GetOrderResult]

type getOrder struct {
	orderRepo domain.OrderRepository
}

func (g getOrder) Handle(ctx context.Context, query GetOrderQuery) (GetOrderResult, error) {
	order, err := g.orderRepo.Get(ctx, query.OrderID, query.CustomerID)
	if err != nil {
		return GetOrderResult{Order: nil}, err
	}
	return GetOrderResult{Order: order}, nil
}

func NewGetOrderHandler(orderRepo domain.OrderRepository,
) GetOrderHandler {
	return cqrs.ApplyQueryDecorator[GetOrderQuery, GetOrderResult](
		getOrder{orderRepo: orderRepo},
	)
}
