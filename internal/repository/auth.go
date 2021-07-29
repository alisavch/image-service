package repository

import (
	"database/sql"

	"github.com/alisavch/image-service/internal/model"
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
func (r *AuthRepository) CreateUser(user model.User) (int, error) {
	var id int
	query := "INSERT INTO image_service.user_account(username, password) VALUES ($1, $2) RETURNING id"
	row := r.db.QueryRow(query, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// GetUser gets the user.
func (r *AuthRepository) GetUser(username string) (model.User, error) {
	var user model.User
	query := "SELECT id, password FROM image_service.user_account where username=$1"
	row := r.db.QueryRow(query, username)
	if err := row.Scan(&user.ID, &user.Password); err != nil {
		return model.User{}, err
	}
	return user, nil
}
