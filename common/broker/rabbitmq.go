package broker

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"time"
)

type RabbitMQRoutingType string

const (
	Direct RabbitMQRoutingType = "direct"
	FanOut RabbitMQRoutingType = "fanout"
)

const (
	DLX            = "dlx" // Dead Letter Exchange
	DLQ            = "dlq" // Dead Letter Queue
	retryHeaderKey = "x-retry-count"
)

var maxRetryCount int

func RabbitMQConnect(user, password, host, port string, maxRetryNumber int) (*amqp.Connection, func() error) {
	maxRetryCount = maxRetryNumber
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

	// Dead Letter Exchange
	err = channel.ExchangeDeclare(DLX, string(FanOut), true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	// Dead Letter Queue
	dlq, err := channel.QueueDeclare(DLQ, true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	err = channel.QueueBind(dlq.Name, "", DLX, false, nil)
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

func RabbitMQRetry(ctx context.Context, channel *amqp.Channel, delivery *amqp.Delivery) error {
	if delivery.Headers == nil {
		delivery.Headers = amqp.Table{}
	}

	retryCount, exists := delivery.Headers[retryHeaderKey].(int)
	if !exists {
		retryCount = 0
	}
	retryCount++
	delivery.Headers[retryHeaderKey] = retryCount

	if retryCount > maxRetryCount {
		logrus.Warnf("max retry count reached, sending to DLQ")
		_ = delivery.Nack(false, false)
		return channel.PublishWithContext(ctx, DLX, "", false, false, amqp.Publishing{
			Headers:      delivery.Headers,
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         delivery.Body,
		})
	}

	logrus.Warnf("retrying message, attempt %d/%d", retryCount, maxRetryCount)
	time.Sleep(time.Second * time.Duration(retryCount))
	return delivery.Nack(false, true)
}
