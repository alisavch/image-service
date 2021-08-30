package model

import "time"

// Request contains information for logs.
//
// A request is additional general information about the user's image.
//
// swagger:model
type Request struct {
	// the ID for this request
	//
	// required: false
	ID int `json:"id,omitempty"`
	// the UserImageID for this request
	//
	// required: false
	UserImageID int `json:"user_image_id,omitempty"`
	// the TimeStart for this request
	//
	// required: false
	TimeStart time.Time `json:"time_start,omitempty"`
	// the EndOfTime for this request
	//
	// required: false
	EndOfTime time.Time `json:"end_of_time,omitempty"`
}
