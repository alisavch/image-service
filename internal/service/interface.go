package service

import (
	"context"
	"io"
	"os"

	"github.com/alisavch/image-service/internal/models"
	"github.com/google/uuid"
)

// AuthorizationRepo consists of authorization methods.
type AuthorizationRepo interface {
	CreateUser(ctx context.Context, user models.User) (id uuid.UUID, err error)
	GetUser(ctx context.Context, username string) (models.User, error)
}

// ImageRepo consists of methods for working with images.
type ImageRepo interface {
	FindUserHistoryByID(ctx context.Context, id uuid.UUID) ([]models.History, error)
	UploadImage(ctx context.Context, img models.UploadedImage) (uuid.UUID, error)
	CreateRequest(ctx context.Context, user models.User, uplImg models.UploadedImage, resImg models.ResultedImage, uI models.UserImage, r models.Request) (uuid.UUID, error)
	FindTheResultingImage(ctx context.Context, id uuid.UUID, service models.Service) (models.ResultedImage, error)
	FindOriginalImage(ctx context.Context, id uuid.UUID) (models.UploadedImage, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.Status) error
}

// S3Bucket contains the basic functions for interacting with the bucket.
type S3Bucket interface {
	UploadToS3Bucket(file io.Reader, filename string) (string, error)
	DownloadFromS3Bucket(filename string) (*os.File, error)
}

// FormattingOutput contains methods for formatting log output.
type FormattingOutput interface {
	Fatalf(format string, args ...interface{})
}
