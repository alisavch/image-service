package model

import "time"

// History contains information about user's history.
//
// A history is general information about requests.
//
// swagger:model
type History struct {
	// the UploadedName for this request
	//
	// required: false
	UploadedName string `json:"uploaded_name,omitempty"`
	// the ResultedName for this request
	//
	// required: false
	ResultedName string `json:"resulted_name,omitempty"`
	// the Service for this request
	//
	// required: false
	Service Service `json:"service,omitempty"`
	// the TimeStart for this request
	//
	// required: false
	TimeStart time.Time `json:"time_start,omitempty"`
	// the EndOfTime for this request
	//
	// required: false
	EndOfTime time.Time `json:"end_of_time,omitempty"`
	// the Status for this request
	//
	// required: false
	Status Status `json:"status,omitempty"`
}
