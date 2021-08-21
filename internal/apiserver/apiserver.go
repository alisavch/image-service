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
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}
	return db, nil
}
