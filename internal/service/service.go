package service

import (
	"github.com/alisavch/image-service/internal/model"
	"github.com/alisavch/image-service/internal/repository"
)

// Authorization contains methods for authorizing users.
type Authorization interface {
	CreateUser(user model.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

// Image contains methods for working with images.
type Image interface {
	UploadImage(image model.UploadedImage) (int, error)
	ConvertToType(uploadedImage model.UploadedImage) (model.ResultedImage, error)
	CompressImage(quality int, uploadedImage model.UploadedImage) (model.ResultedImage, error)
	CreateRequest(model.User, model.UploadedImage, model.ResultedImage, model.UserImage, model.Request) (int, error)
	FindTheResultingImage(id int, service model.Service) (model.ResultedImage, error)
	FindOriginalImage(id int) (model.UploadedImage, error)
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
