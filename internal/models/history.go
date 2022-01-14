package models

import "time"

// History is general information about requests.
type History struct {
	UploadedName  string    `json:"uploaded_name"`
	ResultedName  string    `json:"resulted_name"`
	ServiceName   Service   `json:"service_name"`
	TimeStarted   time.Time `json:"time_started"`
	TimeCompleted time.Time `json:"time_completed"`
	Status        Status    `json:"status"`
}
