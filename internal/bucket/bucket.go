package bucket

import (
	"fmt"
	"io"
	"os"

	"github.com/alisavch/image-service/internal/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	conf = utils.NewConfig()
	sess = connectAWS(NewS3Session())
)

// S3Session contains an environment for the AWS bucket.
type S3Session struct {
	bucketName       string
	region           string
	credentialKey    string
	credentialSecret string
	logger           *Logger
}

// NewS3Session configures S3Session.
func NewS3Session() *S3Session {
	return &S3Session{
		bucketName:       conf.Bucket.BucketName,
		region:           conf.Bucket.AWSRegion,
		credentialKey:    conf.Bucket.AWSAccessKeyID,
		credentialSecret: conf.Bucket.AWSSecretAccessKey,
		logger:           NewLogger(),
	}
}

func connectAWS(s3sess *S3Session) *session.Session {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(s3sess.region),
			Credentials: credentials.NewStaticCredentials(
				s3sess.credentialKey,
				s3sess.credentialSecret,
				""),
		})
	if err != nil {
		s3sess.logger.Fatalf("%s:%s", "Failed to create session", err)
	}

	return sess
}

// UploadToS3Bucket uploads an object to S3.
func (s3sess *S3Session) UploadToS3Bucket(file io.Reader, filename string) (string, error) {
	uploader := s3manager.NewUploader(sess)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   file,
		Bucket: aws.String(s3sess.bucketName),
		Key:    aws.String(filename),
	})
	if err != nil {
		return "", fmt.Errorf("%s:%s", utils.ErrS3Uploading, err)
	}

	s3sess.logger.Printf("%s:%s", "Successfully uploaded", result.Location)
	return result.Location, nil
}

// DownloadFromS3Bucket downloads objects from S3.
func (s3sess *S3Session) DownloadFromS3Bucket(filename string) (*os.File, error) {
	downloader := s3manager.NewDownloader(sess)

	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("%s:%s", utils.ErrCreateFile, err)
	}

	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(s3sess.bucketName),
		Key:    aws.String(filename),
	})
	if err != nil {
		return nil, fmt.Errorf("%s:%s", "failed get object", err)
	}

	s3sess.logger.Printf("%s:%s", "Successfully downloaded", file.Name())
	return file, nil
}
