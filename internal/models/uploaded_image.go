package models

// UploadedImage contains information about uploaded image.
//
// Uploaded image is the general information.
//
// swagger:model UploadedImage
type UploadedImage struct {
	// the ID for this uploaded image
	//
	// required: false
	ID int `json:"id"`

	// the Name for this uploaded image
	//
	// required: false
	Name string `json:"uploaded_name,omitempty"`

	// the Location for this uploaded image
	//
	// required: false
	Location string `json:"uploaded_location,omitempty"`
}
