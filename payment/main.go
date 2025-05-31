package main

import (
	"common/broker"
	_ "common/config"
	"common/metrics"
	"common/server"
	"common/tracing"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"payment/infrastructure/mq"
)

func main() {
	serviceName := viper.GetString("payment.service-name")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	metrics.NewPrometheusClient(
		viper.GetString("payment.metrics-export-host"),
		viper.GetString("payment.metrics-export-port"),
		serviceName+"Counter",
		serviceName+"Histogram",
	)

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

	server.RunHttpServer(serviceName, func(router *gin.Engine) {
		router.POST("/webhook", NewHttpHandler(mq.NewRabbitMQEventSender(rmqChan)).HandleWebhook)
	})
}
