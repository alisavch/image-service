package models

import "github.com/google/uuid"

// Image contains details of a single request.
//
// An image is what is processed in the application.
//
// swagger:model Image
type Image struct {
	// the ID for this image
	//
	// required: false
	ID uuid.UUID `json:"id"`

	// the uploaded name for this image
	//
	// required: true
	UploadedName string `json:"uploaded_name,omitempty"`

	// the uploaded location for this image
	//
	// required: true
	UploadedLocation string `json:"uploaded_location,omitempty"`

	// the resulted name for this image
	//
	// required: false
	ResultedName string `json:"resulted_name,omitempty"`

	// the resulted location for this image
	//
	// required: false
	ResultedLocation string `json:"resulted_location,omitempty"`
}
