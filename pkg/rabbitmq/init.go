package rabbitmq

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"template/pkg/dotenv"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Connection is the connection created
type Connection struct {
	Name         string
	Conn         *amqp.Connection
	Channel      *amqp.Channel
	Queue        map[string]amqp.Queue
	ExchangeName string
	ExchangeType string
	RoutingKey   string
	Err          chan error
}

var (
	connectionPool = make(map[string]*Connection)
)

// NewConnection returns the new connection object
func NewConnection(name, exchangeName string, exchangeType string, routingKey string) RabbitMqTemplate {
	if c, ok := connectionPool[name]; ok {
		return c
	}

	c := &Connection{
		Name: name,
		ExchangeName: exchangeName,
		ExchangeType: exchangeType,
		// RoutingKey:   routingKey,
		Err:          make(chan error),
	}

	connectionPool[name] = c
	c.Queue = make(map[string]amqp.Queue)
	return c
}

// GetConnection returns the connection which was instantiated
func GetConnection(name string) *Connection {
	return connectionPool[name]
}

// GetConn get connection rabbitmq
func (c *Connection) GetConn() *amqp.Connection {
	return c.Conn
}

// GetChan get channel rabbitmq
func (c *Connection) GetChan() *amqp.Channel {
	return c.Channel
}

// Connect connect rabbit mq
func (c *Connection) Connect() error {
	var (
		err     error
		amqpURI = dotenv.GetString("RABBIT_MQ_URL", "")
	)

	if dotenv.GetBool("IS_RABBIT_MQ_TLS", false) {
		c.Conn, err = amqp.DialTLS(amqpURI, &tls.Config{InsecureSkipVerify: true})
	} else {
		c.Conn, err = amqp.Dial(amqpURI)
	}

	if err != nil {
		return fmt.Errorf("Error in creating rabbitmq connection with %s : %s", amqpURI, err.Error())
	}

	go func() {
		<-c.Conn.NotifyClose(make(chan *amqp.Error)) //Listen to NotifyClose
		c.Err <- errors.New("Connection Closed")
	}()

	return nil
}

// ConnectChannel connect channel rabbitmq
func (c *Connection) ConnectChannel() (*amqp.Channel, error) {
	var err error

	c.Channel, err = c.Conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Error in opening Channel: %s", err)
	}

	go func() {
		<-c.Channel.NotifyClose(make(chan *amqp.Error))
		c.Err <- errors.New("Channel Closed")
	}()

	if c.ExchangeName != "" {
		if c.ExchangeType == "" {
			c.ExchangeType = "direct"
		}

		err = c.Channel.ExchangeDeclare(
			c.ExchangeName, // name
			c.ExchangeType, // type
			true,           // durable
			false,          // auto-deleted
			false,          // internal
			false,          // noWait
			nil,            // arguments
		)
		if err != nil {
			return nil, fmt.Errorf("Error in Exchange Declare: %s", err)
		}
	}

	return c.Channel, nil
}

// BindQueue bind queue rabbitmq
func (c *Connection) BindQueue(queue string) error {
	var err error
	c.Queue[queue], err = c.Channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("error in declaring the queue %s", err)
	}

	if c.ExchangeName != "" {
		err = c.Channel.QueueBind(queue, "", c.ExchangeName, false, nil)
		if err != nil {
			return fmt.Errorf("Queue  Bind error: %s", err)
		}
	}

	return nil
}

// Reconnect reconnects the connection
func (c *Connection) Reconnect(queue string) error {
	if err := c.Connect(); err != nil {
		return err
	}

	_, err := c.ConnectChannel()
	if err != nil {
		return err
	}

	if err := c.BindQueue(queue); err != nil {
		return err
	}

	return nil
}

// Consume consume rabbitmq
func (c *Connection) Consume(queue string) (<-chan amqp.Delivery, error) {
	select { //non blocking channel - if there is no error will go to default where we do nothing
	case err := <-c.Err:
		if err != nil {
			err1 := c.Reconnect(queue)
			if err1 != nil {
				c.Err <- err1
			}
		}
	default:
	}

	return c.Channel.Consume(
		c.Queue[queue].Name, // queue
		"",                  // consumer
		true,                // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
}

// Publish publish rabbitmq
func (c *Connection) Publish(ctx context.Context, jsonData []byte, queue string) error {
	select { //non blocking channel - if there is no error will go to default where we do nothing
	case err := <-c.Err:
		if err != nil {
			err1 := c.Reconnect(queue)
			if err1 != nil {
				c.Err <- err1
			}
		}
	default:
	}
	return c.Channel.PublishWithContext(
		ctx,                 // context
		"",                  // exchange
		c.Queue[queue].Name, // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         jsonData,
		},
	)
}
