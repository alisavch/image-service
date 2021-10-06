package service

import (
	"context"
	"image"
	"os"

	"github.com/google/uuid"

	"github.com/alisavch/image-service/internal/models"
	"github.com/alisavch/image-service/internal/repository"
)

// Authorization contains methods for authorizing users.
type Authorization interface {
	CreateUser(ctx context.Context, user models.User) (id uuid.UUID, err error)
	GenerateToken(ctx context.Context, username, password string) (string, error)
	ParseToken(token string) (uuid.UUID, error)
}

// Image contains methods for working with images.
type Image interface {
	UploadImage(ctx context.Context, img models.UploadedImage) (uuid.UUID, error)
	ConvertToType(format, newImageName string, img image.Image, newImg *os.File, isRemoteStorage bool) (models.ResultedImage, error)
	CompressImage(width int, format, resultedName string, img image.Image, newImg *os.File, isRemoteStorage bool) (models.ResultedImage, error)
	CreateRequest(ctx context.Context, user models.User, uplImg models.UploadedImage, resImg models.ResultedImage, uI models.UserImage, r models.Request) (uuid.UUID, error)
	FindTheResultingImage(ctx context.Context, id uuid.UUID, service models.Service) (models.ResultedImage, error)
	FindOriginalImage(ctx context.Context, id uuid.UUID) (models.UploadedImage, error)
	FindUserHistoryByID(ctx context.Context, id uuid.UUID) ([]models.History, error)
	SaveImage(filename, location string, isRemoteStorage bool) (*models.Image, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.Status) error
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
