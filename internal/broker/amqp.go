package broker

// AMQPBroker contains interfaces.
type AMQPBroker struct {
	*RabbitMQ
}

// NewAMQPBroker is the AMQP constructor.
func NewAMQPBroker() *AMQPBroker {
	return &AMQPBroker{NewRabbitMQ()}
}
