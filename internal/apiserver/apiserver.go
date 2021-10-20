package apiserver

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/alisavch/image-service/internal/broker"
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

	db, err := newDB(conf.DBConfig)
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
	services := service.NewService(repos)
	currentService := NewService(services)
	rabbit := broker.NewAMQPBroker()

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
		logger.Printf("%s:%s", "Failed to load .env", err)
	}
}

func newDB(config utils.DBConfig) (*sql.DB, error) {
	URL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.User, config.Password, config.Host, config.Port, config.Name)
	db, err := sql.Open("postgres", URL)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres")
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}
	return db, nil
}
