package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/alisavch/image-service/internal/models"
)

// AuthRepository provides access to the database.
type AuthRepository struct {
	db *sql.DB
}

// NewAuthRepository is constructor of the AuthRepository.
func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// CreateUser provides adding new user.
func (r *AuthRepository) CreateUser(ctx context.Context, user models.User) (id uuid.UUID, err error) {
	query := "INSERT INTO image_service.user_account(username, password) VALUES ($1, $2) RETURNING id"
	row := r.db.QueryRowContext(ctx, query, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return [16]byte{}, fmt.Errorf("cannot insert user into database")
	}
	return id, nil
}

// GetUser gets the user.
func (r *AuthRepository) GetUser(ctx context.Context, username string) (models.User, error) {
	var user models.User
	query := "SELECT id, password FROM image_service.user_account where username=$1"
	row := r.db.QueryRowContext(ctx, query, username)
	if err := row.Scan(&user.ID, &user.Password); err != nil {
		return models.User{}, fmt.Errorf("cannot find the user in database")
	}
	return user, nil
}
