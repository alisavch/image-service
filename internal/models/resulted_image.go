package models

import "github.com/google/uuid"

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
type ResultedImage struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Name     string    `json:"resulted_name,omitempty"`
	Location string    `json:"resulted_location,omitempty"`
	Service  Service   `json:"service,omitempty"`
}
