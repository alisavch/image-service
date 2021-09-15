package apiserver

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/alisavch/image-service/internal/broker"
	"github.com/alisavch/image-service/internal/repository"
	"github.com/alisavch/image-service/internal/service"
	"github.com/alisavch/image-service/internal/utils"
	_ "github.com/lib/pq" // Registers database.
	"github.com/sirupsen/logrus"
)

// Start starts the server.
func Start() error {
	user, pass, host, port, dbname, err := utils.GetDBEnvironments(utils.NewConfig(".env"))
	if err != nil {
		logrus.Fatalf("%s: %s", "Failed to find variables", err)
	}

	db, err := newDB(user, pass, host, port, dbname)
	if err != nil {
		logrus.Fatalf("%s: %s", "Failed to initialize database", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logrus.Fatalf("%s: %s", "Failed to close database", err)
		}
	}(db)

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	rabbit := broker.NewAMQPBroker()

	err = rabbit.Connect()
	if err != nil {
		logrus.Fatalf("%s: %s", "Failed to connect to Rabbitmq", err)
	}

	srv := NewServer(services, rabbit)

	return http.ListenAndServe(
		":8080",
		srv,
	)
}

func newDB(user, pass, host, port, dbname string) (*sql.DB, error) {
	URL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbname)
	db, err := sql.Open("postgres", URL)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres")
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}
	return db, nil
}
