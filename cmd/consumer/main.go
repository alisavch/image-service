package main

import (
	"flag"

	"github.com/alisavch/image-service/internal/consumer"
	_ "github.com/alisavch/image-service/internal/log"
)

func main() {
	logger := consumer.NewLogger()

	flag.Parse()
	logger.Info("The consumer is running")
	consumer.Consume()
	logger.Info("The consumer has stopped receiving messages")
}
