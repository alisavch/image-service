package apiserver

import (
	"fmt"
	"github.com/alisavch/image-service/internal/broker"
	_ "github.com/lib/pq" // Registers database.
	"github.com/sirupsen/logrus"
)

// Consume starts the server.
func Consume() {
	rabbit := new(broker.RabbitMQ)
	if err := rabbit.Connect(); err != nil {
		logrus.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer rabbit.Close()

	q, err := rabbit.DeclareQueue("publisher")
	if err != nil {
		logrus.Fatalf("Failed to declare queue: %s", err)
	}

	fmt.Println("Channel and Queue established")

	err = rabbit.ConsumeQueue(q.Name)
	if err != nil {
		logrus.Fatalf("consume: %s", err)
	}
}
