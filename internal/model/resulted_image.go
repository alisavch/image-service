package model

type (
	// Service is the name of the command being executed.
	Service string
)

const (
	// Conversion is a command with an image.
	Conversion Service = "conversion"
	// Compression is a command with an image.
	Compression Service = "compression"
)

// ResultedImage contains information about resulted image.
//
// A resulted image is the general summary.
//
// swagger:model
type ResultedImage struct {
	// the ID for this resulted image
	//
	// required: false
	ID int `json:"id,omitempty"`
	// the Name for this resulted image
	//
	// required: false
	Name string `json:"resulted_name,omitempty"`
	// the Location for this resulted image
	//
	// required: false
	Location string `json:"resulted_location,omitempty"`
	// the Service for this resulted image
	//
	// required: false
	Service Service `json:"service,omitempty"`
}
