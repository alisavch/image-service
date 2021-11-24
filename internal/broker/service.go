package broker

// Service unites interfaces.
type Service struct {
	Image
	S3Bucket
}

// NewService configures Service.
func NewService(image Image, bucket S3Bucket) *Service {
	return &Service{
		Image:    image,
		S3Bucket: bucket,
	}
}
