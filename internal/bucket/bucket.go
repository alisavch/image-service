package bucket

import (
	"fmt"
	"io"
	"os"

	"github.com/alisavch/image-service/internal/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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
	displayProgress  bool
}

// NewS3Session configures S3Session.
func NewS3Session() *S3Session {
	return &S3Session{
		bucketName:       conf.Bucket.BucketName,
		region:           conf.Bucket.AWSRegion,
		credentialKey:    conf.Bucket.AWSAccessKeyID,
		credentialSecret: conf.Bucket.AWSSecretAccessKey,
		logger:           NewLogger(),
		displayProgress:  true,
		//request: models.RequestStatus{}
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
	uploader := s3manager.NewUploader(sess, func(d *s3manager.Uploader) {
		d.PartSize = 64 * 1024 * 1024 // 64MB per part
		d.Concurrency = 6
	})

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
	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("%s:%s", utils.ErrCreateFile, err)
	}

	s3ObjectSize := s3sess.GetS3ObjectSize(filename)

	downloader := s3manager.NewDownloader(sess, func(d *s3manager.Downloader) {
		d.PartSize = 64 * 1024 * 1024 // 64MB per part
		d.Concurrency = 6
	})

	pw := &progressWriter{writer: file, size: s3ObjectSize}
	pw.display = s3sess.displayProgress
	pw.init(s3ObjectSize)

	numBytes, err := downloader.Download(pw, &s3.GetObjectInput{
		Bucket: aws.String(s3sess.bucketName),
		Key:    aws.String(filename),
	})
	if err != nil {
		return nil, fmt.Errorf("%s:%s", "failed to get object", err)
	}

	pw.finish()

	s3sess.logger.Printf("%s:%s", "Download status", pw.bar.String())
	s3sess.logger.Printf("%s:%s, %d %s", "Successfully downloaded", file.Name(), numBytes, "bytes")
	return file, nil
}

// DeleteItem deletes an item from a bucket
func (s3sess *S3Session) DeleteItem(item string) error {
	svc := s3.New(sess)
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s3sess.bucketName),
		Key:    aws.String(item),
	})
	if err != nil {
		return fmt.Errorf("%s:%s", "failed to delete an object", err)
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(s3sess.bucketName),
		Key:    aws.String(item),
	})
	if err != nil {
		return fmt.Errorf("%s:%s", "failed to wait unlit object not exists", err)
	}

	return nil
}

// GetS3ObjectSize get the size of the file.
func (s3sess *S3Session) GetS3ObjectSize(item string) int64 {
	svc := s3.New(sess)
	input := &s3.HeadObjectInput{
		Bucket: aws.String(s3sess.bucketName),
		Key:    aws.String(item),
	}

	result, err := svc.HeadObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			s3sess.logger.Fatalf("%s:%s", "Error getting size of file", aerr)
			return 0
		}
		s3sess.logger.Fatalf("%s:%s", "Error getting size of file", err)
	}
	return *result.ContentLength
}
