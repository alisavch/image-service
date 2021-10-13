package consumer

// Rabbit unites interfaces.
type Rabbit struct {
	AMQP
}

// NewRabbit configures Rabbit.
func NewRabbit(mq AMQP) *Rabbit {
	return &Rabbit{
		AMQP: mq,
	}
}
