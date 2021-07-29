package model

// UserImage contains details of a single request.
type UserImage struct {
	ID              int    `json:"id,omitempty"`
	UserAccountID   int    `json:"user_account_id,omitempty"`
	UploadedImageID int    `json:"uploaded_image_id,omitempty"`
	ResultedImageID int    `json:"resulted_image_id,omitempty"`
	Status          Status `json:"status,omitempty"`
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
