package apiconsumer

import (
	"github.com/alisavch/image-service/internal/broker"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Consume starts the message consumer.
func Consume() {
	rabbit := new(broker.RabbitMQ)

	err := godotenv.Load(".env")
	if err != nil {
		logrus.Fatalf("%s:%s", "Error loading .env file", err)
	}

	err = rabbit.Connect()
	if err != nil {
		logrus.Fatalf("%s: %s", "Failed to open a channel", err)
	}

	q, err := rabbit.DeclareQueue("publisher")
	if err != nil {
		logrus.Fatalf("%s: %s", "Failed to declare a queue", err)
	}

	err = rabbit.ConsumeQueue(q.Name)
	if err != nil {
		logrus.Fatalf("%s: %s", "Failed to consume a queue", err)
	}
}
