package consumer

import "github.com/streadway/amqp"

// AMQP contains methods for working with message broker.
type AMQP interface {
	Connect() error
	DeclareQueue(name string) (amqp.Queue, error)
	ConsumeQueue(queue string) error
}

// DisplayLog contains methods for log display.
type DisplayLog interface {
	Info(args ...interface{})
	Fatalf(format string, args ...interface{})
}
