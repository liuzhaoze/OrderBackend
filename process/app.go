package main

import (
	client "common/client/order"
	"context"
	"github.com/sirupsen/logrus"
	"process/application"
	"process/application/command"
)

func NewApplication(ctx context.Context) (*application.Application, func()) {
	logger := logrus.StandardLogger()
	orderGrpcClient, closeOrderGrpcClient, err := client.NewOrderGrpcClient(ctx)
	if err != nil {
		logger.Panicln(err)
	}

	return &application.Application{
			Commands: application.Commands{
				ProcessOrder: command.NewProcessOrderHandler(orderGrpcClient, logger),
			},
			Queries: application.Queries{},
		}, func() {
			_ = closeOrderGrpcClient()
		}
}
