package domain

type (
	// Service is the name of the command being executed.
	Service string
)

const (
	// Conversion is a command with an image.
	Conversion  Service = "conversion"
	// Compression is a command with an image.
	Compression Service = "compression"
)

// ResultedImage contains information about resulted image.
type ResultedImage struct {
	ID       int64   `json:"id,omitempty"`
	Name     string  `json:"name,omitempty"`
	Location string  `json:"location,omitempty"`
	Service  Service `json:"service,omitempty"`
}
