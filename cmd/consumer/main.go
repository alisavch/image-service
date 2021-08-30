package main

import (
	"flag"

	"github.com/alisavch/image-service/internal/apiconsumer"

	_ "github.com/alisavch/image-service/internal/log"
	"github.com/sirupsen/logrus"
)

func main() {
	flag.Parse()
	logrus.Info("The consumer is running")
	apiconsumer.Consume()
}
