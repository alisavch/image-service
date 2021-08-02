package service

import (
	"context"

	"github.com/alisavch/image-service/internal/model"
	"github.com/alisavch/image-service/internal/repository"
)

// Authorization contains methods for authorizing users.
type Authorization interface {
	CreateUser(ctx context.Context, user model.User) (int, error)
	GenerateToken(ctx context.Context, username, password string) (string, error)
	ParseToken(token string) (int, error)
}

// Image contains methods for working with images.
type Image interface {
	UploadImage(ctx context.Context, image model.UploadedImage) (int, error)
	ConvertToType(uploadedImage model.UploadedImage) (model.ResultedImage, error)
	CompressImage(quality int, uploadedImage model.UploadedImage) (model.ResultedImage, error)
	CreateRequest(ctx context.Context, user model.User, uplImg model.UploadedImage, resImg model.ResultedImage, uI model.UserImage, r model.Request) (int, error)
	FindTheResultingImage(ctx context.Context, id int, service model.Service) (model.ResultedImage, error)
	FindOriginalImage(ctx context.Context, id int) (model.UploadedImage, error)
	SaveImage(filename, folder, resultedFilename string) error
}

// Service contains interfaces.
type Service struct {
	Authorization
	Image
}

// NewService is Service constructor.
func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Image:         NewImageService(repos.Image),
	}
}
