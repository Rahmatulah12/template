package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqTemplate interface {
	Connect() error
	ConnectChannel() (*amqp.Channel, error)
	BindQueue(queue string) error
	Reconnect(queue string) error
	Consume(queue string) (<-chan amqp.Delivery, error)
	Publish(ctx context.Context, jsonData []byte, queue string) error
	GetConn() *amqp.Connection
	GetChan() *amqp.Channel
}
