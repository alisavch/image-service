package consumer

import (
	"github.com/alisavch/image-service/internal/broker"
)

// Consume starts the message consumer.
func Consume() {
	rabbit := broker.NewAMQPBroker()
	mq := NewRabbit(rabbit)
	logger := NewLogger()

	err := mq.Connect()
	if err != nil {
		logger.Fatalf("%s: %s", "Failed to open a channel", err)
	}

	q, err := mq.DeclareQueue("publisher")
	if err != nil {
		logger.Fatalf("%s: %s", "Failed to declare a queue", err)
	}

	err = mq.ConsumeQueue(q.Name)
	if err != nil {
		logger.Fatalf("%s: %s", "Failed to consume a queue", err)
	}
}
