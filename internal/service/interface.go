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
	FindUserRequestHistory(ctx context.Context, id uuid.UUID) ([]models.History, error)
	FindRequestStatus(ctx context.Context, userID, requestID uuid.UUID) (models.Status, error)
	UploadImage(ctx context.Context, img models.Image) (uuid.UUID, error)
	UploadResultedImage(ctx context.Context, img models.Image) error
	CreateRequest(ctx context.Context, user models.User, img models.Image, req models.Request) (uuid.UUID, error)
	FindResultedImage(ctx context.Context, id uuid.UUID) (models.Image, error)
	FindOriginalImage(ctx context.Context, id uuid.UUID) (models.Image, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.Status) error
	SetCompletedTime(ctx context.Context, id uuid.UUID) error
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
