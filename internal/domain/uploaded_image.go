package domain

// UploadedImage contains information about uploaded image.
type UploadedImage struct {
	ID       int64  `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Location string `json:"location,omitempty"`
}
