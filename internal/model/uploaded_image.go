package model

// UploadedImage contains information about uploaded image.
type UploadedImage struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"uploaded_name,omitempty"`
	Location string `json:"uploaded_location,omitempty"`
}
