package broker

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
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
