package mq

import (
	"common/broker"
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
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
