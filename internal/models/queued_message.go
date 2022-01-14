package models

import "github.com/google/uuid"

// QueuedMessage contains information about the request sent to RabbitMQ.
type QueuedMessage struct {
	Service
	Image
	Width     int
	RequestID uuid.UUID
}

// NewQueuedMessage configures QueuedMessage.
func NewQueuedMessage(width int, requestID uuid.UUID, service Service, image Image) QueuedMessage {
	return QueuedMessage{Service: service, Image: image, Width: width, RequestID: requestID}
}
