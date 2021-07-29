package apiserver

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/alisavch/image-service/internal/repository"
	"github.com/alisavch/image-service/internal/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Registers database.
)

// Start starts the server.
func Start() error {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading environment variables: %s", err.Error())
	}
	db, err := newDB(os.Getenv("DB_SERVER_URL"))
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

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}
	return db, nil
}
