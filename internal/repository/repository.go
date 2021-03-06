package repository

import (
	"database/sql"
	"fmt"

	"github.com/alisavch/image-service/internal/utils"
)

// Repository unites interfaces.
type Repository struct {
	*AuthRepository
	*ImageRepository
}

// NewRepository configures Repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		AuthRepository:  NewAuthRepository(db),
		ImageRepository: NewImageRepository(db),
	}
}

// NewDB configures database.
func NewDB(config utils.DBConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres")
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}
	return db, nil
}
