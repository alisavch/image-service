package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	// Service is the name of the command being executed.
	Service string
	// Status is the information about the execution of the request.
	Status string
)

const (
	// Conversion is a command with an image.
	Conversion Service = "conversion"
	// Compression is a command with an image.
	Compression Service = "compression"
	// Queued is the status of the request.
	Queued Status = "queued"
	// Processing is the status of the request.
	Processing Status = "processing"
	// Done is the status of the request.
	Done Status = "done"
)

// Request contains information for logs.
type Request struct {
	ID            uuid.UUID `json:"id,omitempty"`
	UserAccountID uuid.UUID `json:"user_account_id,omitempty"`
	ImageID       uuid.UUID `json:"image_id,omitempty"`
	ServiceName   Service   `json:"service_name,omitempty"`
	Status        Status    `json:"status,omitempty"`
	TimeStarted   time.Time `json:"time_started,omitempty"`
	TimeCompleted time.Time `json:"time_completed,omitempty"`
}
