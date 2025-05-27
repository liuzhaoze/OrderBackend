package broker

import (
	"context"
	"encoding/json"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	EventOrderCreated = "order.created"
	EventOrderPaid    = "order.paid"
)

type SendEventRequest struct {
	Channel  *amqp.Channel
	Routing  RabbitMQRoutingType
	Exchange string
	Queue    string
	Body     any
}

func SendEvent(ctx context.Context, request *SendEventRequest) error {
	switch request.Routing {
	case Direct:
		queue, err := request.Channel.QueueDeclare(request.Queue, true, false, false, false, nil)
		if err != nil {
			return err
		}

		marshalled, err := json.Marshal(request.Body)
		if err != nil {
			return err
		}

		return request.Channel.PublishWithContext(ctx, request.Exchange, queue.Name, false, false, amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         marshalled,
			Headers:      RabbitMQInsertHeaders(ctx),
		})

	case FanOut:
		marshalled, err := json.Marshal(request.Body)
		if err != nil {
			return err
		}

		return request.Channel.PublishWithContext(ctx, request.Exchange, "", false, false, amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         marshalled,
			Headers:      RabbitMQInsertHeaders(ctx),
		})

	default:
		return errors.New("invalid routing type")
	}
}
