package main

import (
	"common/broker"
	_ "common/config"
	"github.com/spf13/viper"
	"payment/infrastructure/mq"
)

func main() {
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

	select {}
}
