package broker

import (
	"context"
	"image"
	"io"
	"os"

	"github.com/google/uuid"

	"github.com/alisavch/image-service/internal/models"
)

// FormattingOutput contains methods for formatting log output.
type FormattingOutput interface {
	Printf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// S3Bucket contains the basic functions for interacting with the bucket.
type S3Bucket interface {
	UploadToS3Bucket(file io.Reader, filename string) (string, error)
	DownloadFromS3Bucket(filename string) (*os.File, error)
}

// Image contains methods for working with images.
type Image interface {
	CompressImage(width int, format, resultedName string, img image.Image, newImg *os.File, storage string) (models.Image, error)
	UploadResultedImage(ctx context.Context, img models.Image) error
	ChangeFormat(filename string) (string, error)
	ConvertToType(format, resultedName string, img image.Image, newImg *os.File, storage string) (models.Image, error)
	CompleteRequest(ctx context.Context, id uuid.UUID, status models.Status) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.Status) error
}
