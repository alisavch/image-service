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
)

//go:generate swagger generate spec
func main() {
	logger := apiserver.NewLogger()

	flag.Parse()
	logger.Info("The server is running")
	logger.Info("v 1.0.0")
	if err := apiserver.Start(); err != nil {
		logger.Fatalf("error starting server: %s", err.Error())
	}
}
