package service

import (
	"github.com/alisavch/image-service/internal/repository"
)

// Service contains interfaces.
type Service struct {
	*AuthService
	*ImageService
}

// NewService configures Service.
func NewService(repo *repository.Repository, bucket S3Bucket) *Service {
	return &Service{
		AuthService:  NewAuthService(repo.AuthRepository),
		ImageService: NewImageService(repo.ImageRepository, bucket),
	}
}
