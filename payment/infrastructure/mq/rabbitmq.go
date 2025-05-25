package mq

import (
	"common/broker"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"payment/domain"
)

type RabbitMQEventReceiver struct {
}

func NewRabbitMQEventReceiver() *RabbitMQEventReceiver {
	return &RabbitMQEventReceiver{}
}

func (r *RabbitMQEventReceiver) Listen(channel *amqp.Channel) {
	queue, err := channel.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	messages, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Warnf("failed to consume message from queue %s, err: %s", queue.Name, err)
	}

	forever := make(chan struct{})
	go func() {
		for msg := range messages {
			r.OrderCreatedEventHandler(&msg)
		}
	}()
	<-forever
}

func (r *RabbitMQEventReceiver) OrderCreatedEventHandler(msg *amqp.Delivery) {
	logrus.Infof("received message: %s", msg.Body)

	order := &domain.Order{}
	if err := json.Unmarshal(msg.Body, order); err != nil {
		logrus.Warnf("failed to unmarshal order: %s", err)
		_ = msg.Nack(false, false)
		return
	}

	// TODO: implement real logic

	_ = msg.Ack(false)
}
