package models

import "time"

type History struct {
	UploadedName string    `json:"uploaded_name"`
	ResultedName string    `json:"resulted_name"`
	Service      Service   `json:"service"`
	TimeStart    time.Time `json:"time_start"`
	EndOfTime    time.Time `json:"end_of_time"`
	Status       Status    `json:"status"`
}
