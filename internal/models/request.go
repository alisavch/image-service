package models

import (
	"time"

	"github.com/google/uuid"
)

// Request contains information for logs.
type Request struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id,omitempty"`
	UserImageID uuid.UUID `json:"user_image_id,omitempty"`
	TimeStart   time.Time `json:"time_start,omitempty"`
	EndOfTime   time.Time `json:"end_of_time,omitempty"`
}
