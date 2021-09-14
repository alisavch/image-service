package apiserver

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/streadway/amqp"

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

	amqpURL, err := utils.GetRabbitMQURL(utils.NewConfig(".env"))
	if err != nil {
		logrus.Fatalf("%s: %s", "Failed to find variables", err)
	}

	conn, ch, done, err := newAMQP(amqpURL)
	if err != nil {
		logrus.Fatalf("%s: %s", "Failed to connect rabbitmq", err)
	}

	rabbit := broker.NewAMQPBroker(conn, ch, done)

	srv := NewServer(services, rabbit)

	err = conn.Close()
	if err != nil {
		logrus.Fatalf("%s: %s", "Failed to close", err)
	}
	<-done

	return http.ListenAndServe(
		":8080",
		srv,
	)
}

func newAMQP(amqpURL string) (conn *amqp.Connection, ch *amqp.Channel, done chan error, err error) {
	conn, err = amqp.Dial(amqpURL)
	if err != nil {
		logrus.Fatalf("%s: %s", "Failed to connect RabbitMQ", err)
	}

	ch, err = conn.Channel()
	if err != nil {
		logrus.Fatalf("%s: %s", "Failed to open a channel", err)
	}

	done = make(chan error)
	return
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
