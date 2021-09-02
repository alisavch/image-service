package models

import "time"

// Request contains information for logs.
type Request struct {
	ID          int       `json:"id,omitempty"`
	UserImageID int       `json:"user_image_id,omitempty"`
	TimeStart   time.Time `json:"time_start,omitempty"`
	EndOfTime   time.Time `json:"end_of_time,omitempty"`
}
