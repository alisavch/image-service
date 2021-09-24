package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/alisavch/image-service/internal/models"
)

// Authorization consists of authorization methods.
type Authorization interface {
	CreateUser(ctx context.Context, user models.User) (id uuid.UUID, err error)
	GetUser(ctx context.Context, username string) (models.User, error)
}

// Image consists of methods for working with images.
type Image interface {
	FindUserHistoryByID(ctx context.Context, id uuid.UUID) ([]models.History, error)
	UploadImage(ctx context.Context, image models.UploadedImage) (uuid.UUID, error)
	CreateRequest(ctx context.Context, user models.User, uplImg models.UploadedImage, resImg models.ResultedImage, uI models.UserImage, r models.Request) (uuid.UUID, error)
	FindTheResultingImage(ctx context.Context, id uuid.UUID, service models.Service) (models.ResultedImage, error)
	FindOriginalImage(ctx context.Context, id uuid.UUID) (models.UploadedImage, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.Status) error
}

// Repository unites interfaces.
type Repository struct {
	Authorization
	Image
}

// NewRepository is the repository constructor.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthRepository(db),
		Image:         NewImageRepository(db),
	}
}
