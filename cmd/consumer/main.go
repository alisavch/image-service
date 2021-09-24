package main

import (
	"flag"

	"github.com/alisavch/image-service/internal/consumer"
	"github.com/alisavch/image-service/internal/log"

	_ "github.com/alisavch/image-service/internal/log"
)

func main() {
	logger := log.NewCustomLogger()

	flag.Parse()
	logger.Info("The consumer is running")
	consumer.Consume()
}
