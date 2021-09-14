package broker

import "github.com/streadway/amqp"

// AMQP contains methods for working with message broker.
type AMQP interface {
	Connect() error
	Publish(exchange, key string, body string) error
	DeclareQueue(name string) (amqp.Queue, error)
	ConsumeQueue(queue string) error
	QosQueue() error
}

// AMQPBroker contains interfaces.
type AMQPBroker struct {
	AMQP
}

// NewAMQPBroker is the AMQP constructor.
func NewAMQPBroker(conn *amqp.Connection, ch *amqp.Channel, done chan error) *AMQPBroker {
	return &AMQPBroker{
		AMQP: NewRabbitMQ(conn, ch, done),
	}
}
