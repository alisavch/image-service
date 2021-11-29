package models

import "github.com/google/uuid"

// RequestStatus contains information about the status of the request output.
type RequestStatus struct {
	RequestID uuid.UUID `json:"request_id"`
	Status    Status    `json:"status"`
}
