package main

import (
	"flag"

	"github.com/alisavch/image-service/internal/apiserver"
	_ "github.com/alisavch/image-service/internal/log"
	_ "github.com/lib/pq" // Registers database.
	"github.com/sirupsen/logrus"
)

func main() {
	flag.Parse()
	logrus.Info("The consumer is running")
	apiserver.Consume()
}
