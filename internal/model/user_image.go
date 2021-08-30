package model

// UserImage contains details of a single request.
//
// User image is general details of a single user request.
//
// swagger:model
type UserImage struct {
	// the ID for the user image
	//
	// required: false
	ID int `json:"id,omitempty"`
	// the UserAccountID for the user image
	//
	// required: false
	UserAccountID int `json:"user_account_id,omitempty"`
	// the UploadedImageID for the user image
	//
	// required: false
	UploadedImageID int `json:"uploaded_image_id,omitempty"`
	// the ResultedImageID for the user image
	//
	// required: false
	ResultedImageID int `json:"resulted_image_id,omitempty"`
	// the Status for the user image
	//
	// required: false
	Status Status `json:"status,omitempty"`
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
