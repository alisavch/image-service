package repository

import (
	"context"
	"database/sql"

	"github.com/alisavch/image-service/internal/model"
)

// Authorization consists of authorization methods.
type Authorization interface {
	CreateUser(ctx context.Context, user model.User) (id int, err error)
	GetUser(ctx context.Context, username string) (model.User, error)
}

// Image consists of methods for working with images.
type Image interface {
	FindUserHistoryByID(ctx context.Context, id int) ([]model.History, error)
	UploadImage(ctx context.Context, image model.UploadedImage) (int, error)
	CreateRequest(ctx context.Context, user model.User, uplImg model.UploadedImage, resImg model.ResultedImage, uI model.UserImage, r model.Request) (int, error)
	FindTheResultingImage(ctx context.Context, id int, service model.Service) (model.ResultedImage, error)
	FindOriginalImage(ctx context.Context, id int) (model.UploadedImage, error)
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
