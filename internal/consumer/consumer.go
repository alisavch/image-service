package consumer

import (
	"github.com/alisavch/image-service/internal/broker"
	"github.com/alisavch/image-service/internal/log"
)

var logger log.Logger = log.NewCustomLogger()

// Consume starts the message consumer.
func Consume() {
	rabbit := new(broker.RabbitMQ)

	err := rabbit.Connect()
	if err != nil {
		logger.Fatalf("%s: %s", "Failed to open a channel", err)
	}

	q, err := rabbit.DeclareQueue("publisher")
	if err != nil {
		logger.Fatalf("%s: %s", "Failed to declare a queue", err)
	}

	err = rabbit.ConsumeQueue(q.Name)
	if err != nil {
		logger.Fatalf("%s: %s", "Failed to consume a queue", err)
	}
}
