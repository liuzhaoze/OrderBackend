package main

import (
	"common/broker"
	_ "common/config"
	"common/tracing"
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"process/infrastructure/mq"
)

func main() {
	serviceName := viper.GetString("process.service-name")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application, cleanup := NewApplication(ctx)
	defer cleanup()

	shutdown, err := tracing.OTelTracer(
		viper.GetString("zipkin.host"),
		viper.GetString("zipkin.port"),
		serviceName,
	)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = shutdown(ctx)
	}()

	rmqConn, closeRmqConn := broker.RabbitMQConnect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)
	rmqChan := broker.RabbitMQChannel(rmqConn)
	defer func() {
		_ = rmqChan.Close()
		_ = closeRmqConn()
	}()

	eventReceiver := mq.NewRabbitMQEventReceiver(application)
	go eventReceiver.Listen(rmqChan)

	select {}
}
