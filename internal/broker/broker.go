package broker

// AMQPBroker contains interfaces.
type AMQPBroker struct {
	*RabbitMQ
}

// NewAMQPBroker is the AMQP constructor.
func NewAMQPBroker(image Image, bucket S3Bucket) *AMQPBroker {
	return &AMQPBroker{NewRabbitMQ(NewService(image, bucket))}
}
