package mq

import (
	"common/broker"
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"order/domain"
)

type RabbitMQEventSender struct {
	channel *amqp.Channel
}

func NewRabbitMQEventSender(channel *amqp.Channel) *RabbitMQEventSender {
	return &RabbitMQEventSender{channel: channel}
}

func (r *RabbitMQEventSender) Direct(ctx context.Context, event domain.Event) error {
	return broker.SendEvent(ctx, &broker.SendEventRequest{
		Channel:  r.channel,
		Routing:  broker.Direct,
		Exchange: "",
		Queue:    event.Destination,
		Body:     event.Data,
	})
}

type RabbitMQEventReceiver struct {
}

func NewRabbitMQEventReceiver() *RabbitMQEventReceiver {
	return &RabbitMQEventReceiver{}
}

func (r *RabbitMQEventReceiver) Listen(channel *amqp.Channel) {
	queue, err := channel.QueueDeclare("", true, false, true, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	err = channel.QueueBind(queue.Name, "", broker.EventOrderPaid, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	messages, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	forever := make(chan struct{})
	go func() {
		for msg := range messages {
			r.OrderPaidEventHandler(&msg)
		}
	}()
	<-forever
}

func (r *RabbitMQEventReceiver) OrderPaidEventHandler(msg *amqp.Delivery) {
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
