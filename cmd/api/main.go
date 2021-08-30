// Golang SwaggerUI image-service
//
// Documentation of our awesome API
//
//     Schemes: http
//     BasePath: /
//     Version: 1.0.0
//     Host: localhost:8080
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - basic
//
//    SecurityDefinitions:
//    basic:
//      type: basic
//
// swagger:meta
package main

import (
	"flag"

	"github.com/alisavch/image-service/internal/apiserver"

	_ "github.com/alisavch/image-service/internal/log"
	"github.com/sirupsen/logrus"
)

//go:generate swagger generate spec -o ./swagger.json
func main() {
	flag.Parse()
	logrus.Info("The server is running")
	if err := apiserver.Start(); err != nil {
		logrus.Fatalf("error starting server: %s", err.Error())
	}
}
