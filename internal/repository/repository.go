package repository

import "database/sql"

// Repository unites interfaces.
type Repository struct {
	*AuthRepository
	*ImageRepository
}

// NewRepository is the repository constructor.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		AuthRepository:  NewAuthRepository(db),
		ImageRepository: NewImageRepository(db),
	}
}
