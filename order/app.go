package main

import (
	"common/broker"
	client "common/client/stock"
	"common/metrics"
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
	orderRepo, closeOrderRepo, err := database.NewMongoDatabase(
		viper.GetString("mongo.user"),
		viper.GetString("mongo.password"),
		viper.GetString("mongo.host"),
		viper.GetString("mongo.port"),
		viper.GetString("mongo.db-name"),
		viper.GetString("mongo.collection-name"),
	)
	if err != nil {
		logrus.Panicln(err)
	}

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

	metricsClient := metrics.GetPrometheusClient()

	return &application.Application{
			Commands: application.Commands{
				CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGrpcClient, eventSender, logger, metricsClient),
				UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricsClient),
			},
			Queries: application.Queries{
				GetOrder: query.NewGetOrderHandler(orderRepo, logger, metricsClient),
			},
		}, func() {
			_ = closeOrderRepo(ctx)
			_ = closeStockGrpcClient()
			_ = rmqChan.Close()
			_ = closeRmqConn()
		}
}
