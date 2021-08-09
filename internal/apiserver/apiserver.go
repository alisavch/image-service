package apiserver

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/alisavch/image-service/internal/utils"
	"github.com/sirupsen/logrus"

	"github.com/alisavch/image-service/internal/repository"
	"github.com/alisavch/image-service/internal/service"
	_ "github.com/lib/pq" // Registers database.
)

// Start starts the server.
func Start() error {
	utils.LoadEnv()
	get := utils.GetEnvWithKey
	db, err := newDB(
		get("DB_USER"),
		get("DB_PASSWORD"),
		get("DB_HOST"),
		get("DB_PORT"),
		get("DB_NAME"))
	if err != nil {
		logrus.Fatalf("error initialize database: %s", err.Error())
	}
	defer db.Close()
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	srv := newServer(services)

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
