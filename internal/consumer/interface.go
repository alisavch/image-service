package consumer

import (
	"context"
	"image"
	"io"
	"os"

	"github.com/alisavch/image-service/internal/models"
	"github.com/streadway/amqp"
)

// AMQP contains methods for working with message broker.
type AMQP interface {
	Connect() error
	DeclareQueue(name string) (amqp.Queue, error)
	ConsumeQueue(queue string) error
}

// DisplayLog contains methods for log display.
type DisplayLog interface {
	Info(args ...interface{})
	Fatalf(format string, args ...interface{})
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
}
