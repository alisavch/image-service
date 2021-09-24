package models

import "github.com/google/uuid"

// UploadedImage contains information about uploaded image.
//
// Uploaded image is the general information.
//
// swagger:model UploadedImage
type UploadedImage struct {
	// the ID for this uploaded image
	//
	// required: false
	ID uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`

	// the Name for this uploaded image
	//
	// required: false
	Name string `json:"uploaded_name,omitempty"`

	// the Location for this uploaded image
	//
	// required: false
	Location string `json:"uploaded_location,omitempty"`
}
