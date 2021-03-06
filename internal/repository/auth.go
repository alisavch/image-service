package repository

import (
	"context"
	"database/sql"

	"github.com/lib/pq"

	"github.com/alisavch/image-service/internal/models"
	"github.com/alisavch/image-service/internal/utils"

	"github.com/google/uuid"
)

const codeUniqueViolation = "23505"

// AuthRepository provides access to the database.
type AuthRepository struct {
	db *sql.DB
}

// NewAuthRepository configures AuthRepository.
func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// CreateUser provides adding new user.
func (r *AuthRepository) CreateUser(ctx context.Context, user models.User) (id uuid.UUID, err error) {
	query := "INSERT INTO image_service.user_account(username, password) VALUES ($1, $2) RETURNING id"
	err = r.db.QueryRowContext(ctx, query, user.Username, user.Password).Scan(&id)
	if err, ok := err.(*pq.Error); ok && err.Code == codeUniqueViolation {
		return [16]byte{}, utils.ErrUserAlreadyExists
	}
	return id, err
}

// GetUser gets the user.
func (r *AuthRepository) GetUser(ctx context.Context, username string) (models.User, error) {
	var user models.User
	query := "SELECT id, password FROM image_service.user_account where username=$1"
	row := r.db.QueryRowContext(ctx, query, username)
	if err := row.Scan(&user.ID, &user.Password); err != nil {
		return models.User{}, utils.ErrFindUser
	}
	return user, nil
}
