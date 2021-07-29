package repository

import (
	"database/sql"

	"github.com/alisavch/image-service/internal/model"
)

// Authorization consists of authorization methods.
type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUser(username string) (model.User, error)
}

// Image consists of methods for working with images.
type Image interface {
	UploadImage(image model.UploadedImage) (int, error)
	CreateRequest(user model.User, uplImg model.UploadedImage, resImg model.ResultedImage, uI model.UserImage, r model.Request) (int, error)
	FindTheResultingImage(id int, service model.Service) (model.ResultedImage, error)
	FindOriginalImage(id int) (model.UploadedImage, error)
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
