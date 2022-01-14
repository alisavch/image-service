package broker

// ImageService unites interfaces.
type ImageService struct {
	Image
	S3Bucket
}

// NewService configures Service.
func NewService(image Image, bucket S3Bucket) *ImageService {
	return &ImageService{
		Image:    image,
		S3Bucket: bucket,
	}
}
