package mq

import (
	"common/broker"
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"payment/application"
	"payment/application/command"
	"payment/domain"
)

type RabbitMQEventReceiver struct {
	app *application.Application
}

func NewRabbitMQEventReceiver(app *application.Application) *RabbitMQEventReceiver {
	return &RabbitMQEventReceiver{app: app}
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

	if _, err := r.app.Commands.CreatePayment.Handle(context.TODO(), command.CreatePaymentCommand{Order: order}); err != nil {
		logrus.Warnf("failed to create payment for order %s: %s", order.OrderID, err)
		// TODO: retry
	}

	_ = msg.Ack(false)
}
