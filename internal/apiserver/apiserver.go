package apiserver

import (
	"database/sql"
	"net/http"

	"github.com/alisavch/image-service/internal/broker"
	"github.com/alisavch/image-service/internal/bucket"
	"github.com/alisavch/image-service/internal/repository"
	"github.com/alisavch/image-service/internal/service"
	"github.com/alisavch/image-service/internal/utils"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Registers database.
)

// Start starts the server.
func Start() error {
	logger := NewLogger()
	initEnvironments()

	conf := utils.NewConfig()

	db, err := repository.NewDB(conf.DBConfig)
	if err != nil {
		logger.Fatalf("%s: %s", "Failed to initialize database", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Fatalf("%s: %s", "Failed to close database", err)
		}
	}(db)

	repos := repository.NewRepository(db)
	aws := bucket.NewAWS()
	services := service.NewService(repos, aws)
	currentService := NewAPI(services, aws)
	rabbit := broker.NewAMQPBrokerAPI()

	err = rabbit.Connect()
	if err != nil {
		logger.Fatalf("%s: %s", "Failed to connect to Rabbitmq", err)
	}

	srv := NewServer(rabbit, currentService)

	return http.ListenAndServe(
		":8080",
		srv,
	)
}

func initEnvironments() {
	logger := NewLogger()
	if err := godotenv.Load(); err != nil {
		logger.Printf("%s:%s", "The remote environment is used", err)
	}
}
