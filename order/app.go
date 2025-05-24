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
	orderRepo := database.NewMemoryDatabase()
	stockGrpcClient, closeStockGrpcClient, err := client.NewStockGrpcClient(ctx)
	if err != nil {
		logrus.Panicln(err)
	}
	return &application.Application{
			Commands: application.Commands{
				CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGrpcClient),
				UpdateOrder: command.NewUpdateOrderHandler(orderRepo),
			},
			Queries: application.Queries{
				GetOrder: query.NewGetOrderHandler(orderRepo),
			},
		}, func() {
			_ = closeStockGrpcClient()
		}
}
