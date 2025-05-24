package main

import (
	client "common/client/stock"
	"context"
	"github.com/sirupsen/logrus"
	"order/application"
	"order/application/command"
	"order/application/query"
	"order/infrastructure/database"
)

func NewApplication(ctx context.Context) (*application.Application, func()) {
	logger := logrus.StandardLogger()
	orderRepo := database.NewMemoryDatabase()
	stockGrpcClient, closeStockGrpcClient, err := client.NewStockGrpcClient(ctx)
	if err != nil {
		logrus.Panicln(err)
	}
	return &application.Application{
			Commands: application.Commands{
				CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGrpcClient, logger),
				UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger),
			},
			Queries: application.Queries{
				GetOrder: query.NewGetOrderHandler(orderRepo, logger),
			},
		}, func() {
			_ = closeStockGrpcClient()
		}
}
