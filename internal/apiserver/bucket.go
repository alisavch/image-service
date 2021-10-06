package apiserver

import (
	"io"
	"os"

	"github.com/alisavch/image-service/internal/bucket"
)

// S3Bucket contains the basic functions for interacting with the bucket.
type S3Bucket interface {
	UploadToS3Bucket(file io.Reader, filename string) (string, error)
	DownloadFromS3Bucket(filename string) (*os.File, error)
}

// AWS contains amazon services.
type AWS struct {
	S3Bucket
}

// NewAWS is the AWS constructor.
func NewAWS() *AWS {
	return &AWS{bucket.NewS3Session()}
}
