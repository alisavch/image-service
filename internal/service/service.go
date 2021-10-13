package service

import (
	"github.com/alisavch/image-service/internal/bucket"
	"github.com/alisavch/image-service/internal/repository"
)

// Service contains interfaces.
type Service struct {
	*AuthService
	*ImageService
	*bucket.AWS
}

// NewService configures Service.
func NewService(repo *repository.Repository) *Service {
	return &Service{
		AuthService:  NewAuthService(repo.AuthRepository),
		ImageService: NewImageService(repo.ImageRepository),
		AWS:          bucket.NewAWS(),
	}
}
