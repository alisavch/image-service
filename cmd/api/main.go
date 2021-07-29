package main

import (
	"flag"

	"github.com/sirupsen/logrus"

	"github.com/alisavch/image-service/internal/apiserver"
	_ "github.com/alisavch/image-service/internal/log"
)

func main() {
	flag.Parse()
	logrus.Info("The server is running")
	if err := apiserver.Start(); err != nil {
		logrus.Fatalf("error starting server: %s", err.Error())
	}
}
