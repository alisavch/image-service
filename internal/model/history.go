package model

import "time"

// History contains information about user's history.
type History struct {
	UploadedName string    `json:"uploaded_name,omitempty"`
	ResultedName string    `json:"resulted_name,omitempty"`
	Service      Service   `json:"service,omitempty"`
	TimeStart    time.Time `json:"time_start,omitempty"`
	EndOfTime    time.Time `json:"end_of_time,omitempty"`
	Status       Status    `json:"status,omitempty"`
}
