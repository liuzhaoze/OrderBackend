package mq

import (
	"common/broker"
	"common/tracing"
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
			r.OrderCreatedEventHandler(channel, &msg)
		}
	}()
	<-forever
}

func (r *RabbitMQEventReceiver) OrderCreatedEventHandler(channel *amqp.Channel, msg *amqp.Delivery) {
	ctx, span := tracing.StartSpan(broker.RabbitMQExtractHeaders(context.Background(), msg.Headers), "Payment/MQ: 处理订单创建完成事件")
	defer span.End()

	logrus.Infof("received message: %s", msg.Body)

	order := &domain.Order{}
	if err := json.Unmarshal(msg.Body, order); err != nil {
		logrus.Warnf("failed to unmarshal order: %s", err)
		_ = msg.Nack(false, false)
		return
	}

	if _, err := r.app.Commands.CreatePayment.Handle(ctx, command.CreatePaymentCommand{Order: order}); err != nil {
		logrus.Warnf("failed to create payment for order %s: %s", order.OrderID, err)
		if err = broker.RabbitMQRetry(ctx, channel, msg); err != nil {
			logrus.Errorf("failed to retry message: %v", err)
			_ = msg.Nack(false, false)
		}
		return
	}

	logrus.Info("message processed successfully")
	_ = msg.Ack(false)
}

type RabbitMQEventSender struct {
	channel *amqp.Channel
}

func NewRabbitMQEventSender(channel *amqp.Channel) *RabbitMQEventSender {
	return &RabbitMQEventSender{channel: channel}
}

func (r *RabbitMQEventSender) Broadcast(ctx context.Context, event domain.Event) error {
	return broker.SendEvent(ctx, &broker.SendEventRequest{
		Channel:  r.channel,
		Routing:  broker.FanOut,
		Exchange: event.Destination,
		Queue:    "",
		Body:     event.Data,
	})
}
