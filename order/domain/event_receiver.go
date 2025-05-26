package domain

import amqp "github.com/rabbitmq/amqp091-go"

type EventReceiver interface {
	OrderPaidEventHandler(msg *amqp.Delivery)
}
