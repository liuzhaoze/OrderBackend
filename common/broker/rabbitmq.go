package broker

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

type RabbitMQRoutingType string

const (
	Direct RabbitMQRoutingType = "direct"
	FanOut RabbitMQRoutingType = "fanout"
)

func RabbitMQConnect(user, password, host, port string) (*amqp.Connection, func() error) {
	address := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)
	conn, err := amqp.Dial(address)
	if err != nil {
		logrus.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	return conn, conn.Close
}

func RabbitMQChannel(conn *amqp.Connection) *amqp.Channel {
	channel, err := conn.Channel()
	if err != nil {
		logrus.Fatalf("Failed to open a channel: %s", err)
	}

	err = channel.ExchangeDeclare(EventOrderCreated, string(Direct), true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	err = channel.ExchangeDeclare(EventOrderPaid, string(FanOut), true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	return channel
}

type RabbitMQCarrier map[string]interface{}

func (r RabbitMQCarrier) Get(key string) string {
	if val, ok := r[key]; ok {
		return val.(string)
	} else {
		return ""
	}
}

func (r RabbitMQCarrier) Set(key string, value string) {
	r[key] = value
}

func (r RabbitMQCarrier) Keys() []string {
	keys := make([]string, len(r))
	i := 0
	for k := range r {
		keys[i] = k
	}
	return keys
}

func RabbitMQInsertHeaders(ctx context.Context) map[string]interface{} {
	headers := make(RabbitMQCarrier)
	otel.GetTextMapPropagator().Inject(ctx, headers)
	return headers
}

func RabbitMQExtractHeaders(ctx context.Context, headers map[string]interface{}) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, RabbitMQCarrier(headers))
}
