package mq

import (
	"common/broker"
	"common/tracing"
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"order/application"
	"order/application/command"
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
	app *application.Application
}

func NewRabbitMQEventReceiver(app *application.Application) *RabbitMQEventReceiver {
	return &RabbitMQEventReceiver{app: app}
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
	ctx, span := tracing.StartSpan(broker.RabbitMQExtractHeaders(context.Background(), msg.Headers), "Order/MQ: 处理订单支付完成事件")
	defer span.End()

	logrus.Infof("received message: %s", msg.Body)

	order := &domain.Order{}
	if err := json.Unmarshal(msg.Body, order); err != nil {
		logrus.Warnf("failed to unmarshal order: %s", err)
		_ = msg.Nack(false, false)
		return
	}

	_, err := r.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrderCommand{
		Order: order,
		UpdateFunc: func(c context.Context, o *domain.Order) (*domain.Order, error) {
			if err := o.UpdateStatus(order.Status); err != nil {
				return nil, err
			}
			return o, nil
		},
	})
	if err != nil {
		logrus.Warnf("failed to update order %s: %s", order.OrderID, err)
		// TODO: retry
	}

	_ = msg.Ack(false)
}
