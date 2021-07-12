package domain

import "time"

// Request contains information for logs.
type Request struct {
	ID          int64     `json:"id,omitempty"`
	UserImageID int64     `json:"user_image_id,omitempty"`
	TimeStart   time.Time `json:"time_start,omitempty"`
	EndOfTime   time.Time `json:"end_of_time,omitempty"`
}
