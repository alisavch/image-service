package domain

// UserImage contains details of a single request.
type UserImage struct {
	ID              int64  `json:"id,omitempty"`
	UserAccountID   int64  `json:"user_account_id,omitempty"`
	UploadedImageID int64  `json:"uploaded_image_id,omitempty"`
	Status          Status `json:"status,omitempty"`
}

type (
	// Status is the information about the execution of the request.
	Status string
)

const (
	// Queued is the status of the request.
	Queued     Status = "queued"
	// Processing is the status of the request.
	Processing Status = "processing"
	// Done is the status of the request.
	Done       Status = "done"
)
