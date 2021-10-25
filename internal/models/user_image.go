package models

import "github.com/google/uuid"

// UserImage contains details of a single request.
type UserImage struct {
	ID              uuid.UUID `json:"id"`
	UserAccountID   uuid.UUID `json:"user_account_id,omitempty"`
	UploadedImageID uuid.UUID `json:"uploaded_image_id,omitempty"`
	ResultedImageID uuid.UUID `json:"resulted_image_id,omitempty"`
	Status          Status    `json:"status,omitempty"`
}

type (
	// Status is the information about the execution of the request.
	Status string
)

const (
	// Queued is the status of the request.
	Queued Status = "queued"
	// Processing is the status of the request.
	Processing Status = "processing"
	// Done is the status of the request.
	Done Status = "done"
)
