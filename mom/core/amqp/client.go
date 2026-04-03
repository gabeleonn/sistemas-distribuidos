package amqp

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn         *amqp.Connection
	ch           *amqp.Channel
	exchangeName string
}

func New() *Client {
	amqpURL := os.Getenv("AMQP_URL")
	exchangeName := os.Getenv("EXCHANGE_NAME")

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		panic("falha ao estabelecer conexão com o cliente: " + err.Error())
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		panic("falha ao criar canal do cliente: " + err.Error())
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		ch.Close()
		conn.Close()
		panic("falha ao declarar exchange: " + err.Error())
	}

	return &Client{
		conn:         conn,
		ch:           ch,
		exchangeName: exchangeName,
	}
}

func (c *Client) DeclareQueue(name string, durable bool, autoDelete bool, exclusive bool) (amqp.Queue, error) {
	return c.ch.QueueDeclare(name, durable, autoDelete, exclusive, false, nil)
}

func (c *Client) BindQueue(queueName string, key string) error {
	return c.ch.QueueBind(queueName, key, c.exchangeName, false, nil)
}

func (c *Client) Publish(key string, body []byte) error {
	return c.ch.Publish(
		c.exchangeName,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}

func (c *Client) Consume(queueName string, consumer string) (<-chan amqp.Delivery, error) {
	return c.ch.Consume(queueName, consumer, true, false, false, false, nil)
}

func (c *Client) Close() {
	if c.ch != nil {
		c.ch.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
