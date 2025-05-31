package command

import (
	"common/broker"
	"common/consts"
	"common/cqrs"
	"common/metrics"
	"common/protobuf/stockpb"
	"common/tracing"
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
	"order/domain"
	"order/dto"
)

type CreateOrderCommand struct {
	CustomerID string
	Items      []*domain.ItemWithQuantity
}

type CreateOrderResult struct {
	OrderID string
}

type CreateOrderHandler cqrs.CommandHandler[CreateOrderCommand, CreateOrderResult]

type createOrder struct {
	orderRepo   domain.OrderRepository
	stockGrpc   stockpb.StockServiceClient
	eventSender domain.EventSender
}

func (c createOrder) Handle(ctx context.Context, command CreateOrderCommand) (CreateOrderResult, error) {
	ctx, span := tracing.StartSpan(ctx, "Order/Application/Command: create order")
	defer span.End()

	packedItems := packItems(command.Items)
	grpcRequest := &stockpb.CheckAndFetchItemsRequest{Items: dto.NewItemWithQuantityConverter().ToStockGrpcBatch(packedItems)}
	// 通过 stock gRPC 校验库存是否充足
	// 库存充足：扣减订单对应的库存，并返回剩余库存
	// 库存不足：直接返回订单对应的库存，不扣减库存
	grpcResponse, grpcErr := c.stockGrpc.CheckAndFetchItems(ctx, grpcRequest)
	if grpcErr != nil {
		return CreateOrderResult{OrderID: ""}, status.Convert(grpcErr).Err()
	}

	stockStatus, convertErr := dto.NewStockStatusConverter().FromStockGrpc(grpcResponse.StatusCode)
	if convertErr != nil {
		return CreateOrderResult{OrderID: ""}, convertErr
	}
	stockItems := dto.NewItemConverter().FromStockGrpcBatch(grpcResponse.Items)

	switch stockStatus {
	case consts.StockStatusInsufficient:
		span.AddEvent("库存不足")
		// 库存不足，将当前库存添加到错误信息中返回
		text := "insufficient stock to create order, current stock:"
		for _, item := range stockItems {
			text += fmt.Sprintf(" %+v", *item)
		}
		return CreateOrderResult{OrderID: ""}, errors.New(text)

	case consts.StockStatusSufficient:
		span.AddEvent("库存充足")
		// 库存充足，使用请求物品数量覆盖库存物品数量，目的是保留 Name 和 PriceID
		for _, from := range packedItems {
			for _, to := range stockItems {
				if from.ItemID == to.ItemID {
					to.Quantity = from.Quantity
				}
			}
		}

		pendingOrder, err := domain.NewPendingOrder(command.CustomerID, stockItems)
		if err != nil {
			return CreateOrderResult{OrderID: ""}, err
		}

		order, err := c.orderRepo.Create(ctx, pendingOrder)
		if err != nil {
			return CreateOrderResult{OrderID: ""}, err
		}

		err = c.eventSender.Direct(ctx, domain.Event{Destination: broker.EventOrderCreated, Data: order})
		if err != nil {
			return CreateOrderResult{OrderID: ""}, err
		}
		span.AddEvent("send to order.created MQ (direct)")

		return CreateOrderResult{OrderID: order.OrderID}, nil

	default:
		return CreateOrderResult{OrderID: ""}, errors.New("unknown stock status from stock gRPC")
	}
}

func NewCreateOrderHandler(orderRepo domain.OrderRepository, stockGrpc stockpb.StockServiceClient, eventSender domain.EventSender,
	logger *logrus.Logger,
	metricsClient metrics.Client,
) CreateOrderHandler {
	return cqrs.ApplyCommandDecorator[CreateOrderCommand, CreateOrderResult](
		createOrder{orderRepo: orderRepo, stockGrpc: stockGrpc, eventSender: eventSender},
		logger,
		metricsClient,
	)
}

// packItems 将重复物品的记录合并在一起
func packItems(items []*domain.ItemWithQuantity) []*domain.ItemWithQuantity {
	packed := make(map[string]int64)
	for _, item := range items {
		packed[item.ItemID] += item.Quantity
	}

	packedItems := make([]*domain.ItemWithQuantity, 0, len(packed))
	for id, quantity := range packed {
		packedItems = append(packedItems, &domain.ItemWithQuantity{
			ItemID:   id,
			Quantity: quantity,
		})
	}

	return packedItems
}
