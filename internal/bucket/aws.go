package bucket

// AWS contains amazon services.
type AWS struct {
	*S3Session
}

// NewAWS configures AWS.
func NewAWS() *AWS {
	return &AWS{NewS3Session()}
}
