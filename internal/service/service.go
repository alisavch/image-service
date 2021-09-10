package service

import (
	"context"

	"github.com/alisavch/image-service/internal/models"
	"github.com/alisavch/image-service/internal/repository"
)

// Authorization contains methods for authorizing users.
type Authorization interface {
	CreateUser(ctx context.Context, user models.User) (id int, err error)
	GenerateToken(ctx context.Context, username, password string) (string, error)
	ParseToken(token string) (int, error)
}

// Image contains methods for working with images.
type Image interface {
	UploadImage(ctx context.Context, image models.UploadedImage) (int, error)
	ConvertToType(uploadedImage models.UploadedImage) (models.ResultedImage, error)
	CompressImage(quality int, uploadedImage models.UploadedImage) (models.ResultedImage, error)
	CreateRequest(ctx context.Context, user models.User, uplImg models.UploadedImage, resImg models.ResultedImage, uI models.UserImage, r models.Request) (int, error)
	FindTheResultingImage(ctx context.Context, id int, service models.Service) (models.ResultedImage, error)
	FindOriginalImage(ctx context.Context, id int) (models.UploadedImage, error)
	FindUserHistoryByID(ctx context.Context, id int) ([]models.History, error)
	SaveImage(filename, folder string) (*models.Image, error)
	UpdateStatus(ctx context.Context, id int, status models.Status) error
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
