package models

// QueuedMessage contains information about the request sent to RabbitMQ.
type QueuedMessage struct {
	Service
	Image
	Width int
}

// NewQueuedMessage configures QueuedMessage.
func NewQueuedMessage(width int, service Service, image Image) QueuedMessage {
	return QueuedMessage{Service: service, Image: image, Width: width}
}
