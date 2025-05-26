package main

import (
	"common/broker"
	_ "common/config" // import for side effect to load configuration
	"common/discovery"
	"common/protobuf/orderpb"
	"common/server"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"order/infrastructure/mq"
	"order/ports"
)

func main() {
	serviceName := viper.GetString("order.service-name")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application, cleanup := NewApplication(ctx)
	defer cleanup()

	deregisterConsul, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer deregisterConsul()

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

	eventReceiver := mq.NewRabbitMQEventReceiver()
	go eventReceiver.Listen(rmqChan)

	go server.RunGrpcServer(serviceName, func(s *grpc.Server) {
		orderpb.RegisterOrderServiceServer(s, NewGrpcHandler(application))
	})

	server.RunHttpServer(serviceName, func(router *gin.Engine) {
		router.StaticFile("/payment", "../public/payment.html")
		ports.RegisterHandlersWithOptions(router, NewHttpHandler(application), ports.GinServerOptions{
			BaseURL:      "/api",
			Middlewares:  nil, // 中间件在 RunHttpServer 中统一设置
			ErrorHandler: nil,
		})
	})
}
