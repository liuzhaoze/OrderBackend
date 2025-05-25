package main

import (
	"common/broker"
	client "common/client/stock"
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"order/application"
	"order/application/command"
	"order/application/query"
	"order/infrastructure/database"
	"order/infrastructure/mq"
)

func NewApplication(ctx context.Context) (*application.Application, func()) {
	logger := logrus.StandardLogger()
	orderRepo := database.NewMemoryDatabase()
	stockGrpcClient, closeStockGrpcClient, err := client.NewStockGrpcClient(ctx)
	if err != nil {
		logrus.Panicln(err)
	}

	rmqConn, closeRmqConn := broker.RabbitMQConnect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)
	rmqChan := broker.RabbitMQChannel(rmqConn)
	eventSender := mq.NewRabbitMQEventSender(rmqChan)

	return &application.Application{
			Commands: application.Commands{
				CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGrpcClient, eventSender, logger),
				UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger),
			},
			Queries: application.Queries{
				GetOrder: query.NewGetOrderHandler(orderRepo, logger),
			},
		}, func() {
			_ = closeStockGrpcClient()
			_ = rmqChan.Close()
			_ = closeRmqConn()
		}
}
