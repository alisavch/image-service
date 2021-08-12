package apiserver

import (
	"database/sql"
	"fmt"
	"github.com/alisavch/image-service/internal/broker"
	"github.com/alisavch/image-service/internal/model"
	"net/http"

	"github.com/alisavch/image-service/internal/utils"
	"github.com/sirupsen/logrus"

	"github.com/alisavch/image-service/internal/repository"
	"github.com/alisavch/image-service/internal/service"
	_ "github.com/lib/pq" // Registers database.
)

// Start starts the server.
func Start() error {
	conf := utils.NewConfig(".env")
	user, pass, host, port, dbname, err := utils.GetDBEnvironments(conf)
	if err != nil {
		logrus.Fatalf("error find variables :%s", err.Error())
	}
	db, err := newDB(user, pass, host, port, dbname)
	if err != nil {
		logrus.Fatalf("error initialize database: %s", err.Error())
	}
	defer db.Close()
	repos := repository.NewRepository(db)
	services := service.NewService(repos)

	rabbit := new(broker.RabbitMQ)
	if err = rabbit.Connect(); err != nil {
		logrus.Fatalf("rabbit connection: %s", err)
	}
	defer rabbit.Close()

	if err = rabbit.DeclareQueue(model.Queued); err != nil {
		logrus.Fatalf("declare queue: %s", err)
	}

	srv := newServer(services, rabbit)



	return http.ListenAndServe(
		":8080",
		srv,
	)
}

func newDB(user, pass, host, port, dbname string) (*sql.DB, error) {
	URL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbname)
	db, err := sql.Open("postgres", URL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}
	return db, nil
}
