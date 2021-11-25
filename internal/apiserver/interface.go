package apiserver

import (
	"context"
	"image"
	"io"
	"os"

	"github.com/alisavch/image-service/internal/models"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// AMQP contains methods for working with message broker.
type AMQP interface {
	Publish(exchange, key string, message models.QueuedMessage) error
	DeclareQueue(name string) (amqp.Queue, error)
	QosQueue() error
}

// DisplayLog contains methods for log display.
type DisplayLog interface {
	Info(args ...interface{})
	Printf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

// Authorization contains methods for authorizing users.
type Authorization interface {
	CreateUser(ctx context.Context, user models.User) (id uuid.UUID, err error)
	GenerateToken(ctx context.Context, username, password string) (string, error)
	ParseToken(token string) (uuid.UUID, error)
}

// Image contains methods for working with images.
type Image interface {
	CompressImage(width int, format, resultedName string, img image.Image, newImg *os.File, storage string) (models.Image, error)
	UploadResultedImage(ctx context.Context, img models.Image) error
	ChangeFormat(filename string) (string, error)
	ConvertToType(format, resultedName string, img image.Image, newImg *os.File, storage string) (models.Image, error)
	CheckStatus(ctx context.Context, id uuid.UUID) error
	UploadImage(ctx context.Context, img models.Image) (uuid.UUID, error)
	CreateRequest(ctx context.Context, user models.User, img models.Image, req models.Request) (uuid.UUID, error)
	FindResultedImage(ctx context.Context, id uuid.UUID) (models.Image, error)
	FindOriginalImage(ctx context.Context, id uuid.UUID) (models.Image, error)
	FindUserRequestHistory(ctx context.Context, id uuid.UUID) ([]models.History, error)
	SaveImage(filename, location, storage string) (*models.SavedImage, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.Status) error
	FillInTheResultingImageForAWS(resultedName string) (models.Image, error)
	FillInTheResultingImage(storage, resultedName string, newImg *os.File) (models.Image, error)
	SetCompletedTime(ctx context.Context, id uuid.UUID) error
}

// S3Bucket contains the basic functions for interacting with the bucket.
type S3Bucket interface {
	UploadToS3Bucket(file io.Reader, filename string) (string, error)
	DownloadFromS3Bucket(filename string) (*os.File, error)
}

// ServiceOperations combines the basic service operations.
type ServiceOperations interface {
	Authorization
	Image
}
