package broker

// AMQPBrokerAPI contains interfaces.
type AMQPBrokerAPI struct {
	*ProcessMessage
}

// NewAMQPBrokerAPI configures AMQP.
func NewAMQPBrokerAPI() *AMQPBrokerAPI {
	return &AMQPBrokerAPI{NewProcessMessageAPI()}
}

// AMQPBrokerConsumer contains interfaces for rabbitmq and message handling.
type AMQPBrokerConsumer struct {
	*ProcessMessage
}

// NewAMQPBrokerConsumer configures AMQPBrokerConsumer.
func NewAMQPBrokerConsumer(image Image, bucket S3Bucket) *AMQPBrokerConsumer {
	return &AMQPBrokerConsumer{ProcessMessage: NewProcessMessageConsumer(NewService(image, bucket))}
}
