package main

import (
	"common/broker"
	_ "common/config"
	"common/server"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"payment/infrastructure/mq"
)

func main() {
	serviceName := viper.GetString("payment.service-name")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application, cleanup := NewApplication(ctx)
	defer cleanup()

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
