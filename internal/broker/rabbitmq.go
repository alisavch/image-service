package broker

import (
	"fmt"

	"github.com/alisavch/image-service/internal/model"
	"github.com/alisavch/image-service/internal/utils"
	"github.com/streadway/amqp"
)

// RabbitMQ Operate Wrapper
type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	done chan error
}

// Connect instantiates the RabbitMW instances using configuration defined in environment variables.
func (r *RabbitMQ) Connect() (err error) {
	url := utils.GetEnvWithKey("RABBITMQ_URL")

	r.conn, err = amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("connection.open: %w", err)
	}

	r.ch, err = r.conn.Channel()
	if err != nil {
		return fmt.Errorf("channel.open: %w", err)
	}

	r.done = make(chan error)

	return nil
}

// Publish sends data to the queue.
func (r *RabbitMQ) Publish(exchange, key string, deliveryMode, priority uint8, body string) (err error) {
	err = r.ch.Publish(exchange, key, false, false,
		amqp.Publishing{
			Headers:      amqp.Table{},
			ContentType:  "application/json",
			DeliveryMode: deliveryMode,
			Priority:     priority,
			Body:         []byte(body),
		})
	if err != nil {
		return fmt.Errorf("publish message errro: %w", err)
	}
	return nil
}

// DeclareQueue declares a queue.
func (r *RabbitMQ) DeclareQueue(name model.Status) (err error) {
	_, err = r.ch.QueueDeclare(string(name), true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declare queue errror: %w", err)
	}
	return nil
}

// DeleteQueue removes the queue from the server.
func (r *RabbitMQ) DeleteQueue(name string) (err error) {
	_, err = r.ch.QueueDelete(name, false, false, false)
	if err != nil {
		return fmt.Errorf("delete queue error: %w", err)
	}
	return nil
}

// ConsumeQueue starts delivering queued messages.
func (r *RabbitMQ) ConsumeQueue(queue string, message chan []byte) (err error) {
	deliveries, err := r.ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume queue error: %w", err)
	}
	go func(deliveries <-chan amqp.Delivery, done chan error, message chan []byte) {
		for d := range deliveries {
			message <- d.Body
		}
		done <- nil
	}(deliveries, r.done, message)
	return nil
}

// Close closes requests.
func (r *RabbitMQ) Close() (err error) {
	err = r.conn.Close()
	if err != nil {
		return fmt.Errorf("close error: %w", err)
	}
	return nil
}
