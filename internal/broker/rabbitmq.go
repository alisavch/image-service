package broker

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/alisavch/image-service/internal/models"
	"github.com/alisavch/image-service/internal/utils"

	"github.com/streadway/amqp"
)

// RabbitMQ Operate Wrapper.
type RabbitMQ struct {
	*Service
	conn   *amqp.Connection
	ch     *amqp.Channel
	done   chan error
	logger *Logger
	*models.RequestStatus
}

// NewRabbitMQ configures RabbitMQ.
func NewRabbitMQ(service *Service) *RabbitMQ {
	return &RabbitMQ{Service: service, logger: NewLogger()}
}

// Connect instantiates the RabbitMW instances using configuration defined in environment variables.
func (r *RabbitMQ) Connect() error {
	conf := utils.NewConfig()
	var err error
	r.conn, err = amqp.Dial(conf.Rabbitmq.RabbitmqURL)
	if err != nil {
		r.logger.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	r.ch, err = r.conn.Channel()
	if err != nil {
		r.logger.Fatalf("%s: %s", "Failed to open a channel", err)
	}

	r.done = make(chan error)

	return nil
}

// Publish sends data to the queue.
func (r *RabbitMQ) Publish(exchange, key string, message models.QueuedMessage) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("can't marshal queue message: %w", err)
	}

	err = r.ch.Publish(exchange, key, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		r.logger.Fatalf("%s:%s", "Failed to publish a message", err)
	}
	return nil
}

// DeclareQueue declares a queue.
func (r *RabbitMQ) DeclareQueue(name string) (amqp.Queue, error) {
	q, err := r.ch.QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		r.logger.Fatalf("%s: %s", "Failed to declare a queue", err)
	}
	return q, nil
}

// ConsumeQueue starts delivering queued messages.
func (r *RabbitMQ) ConsumeQueue(queue string) error {
	var message models.QueuedMessage
	deliveries, err := r.ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("%s: %s", "Failed to register consumer", err)
	}
	forever := make(chan bool)
	for d := range deliveries {
		go func() {
			err := json.NewDecoder(bytes.NewReader(d.Body)).Decode(&message)
			if err != nil {
				r.logger.Printf("%s: %s", "Failed to decode json", err)
			}

			err = r.Process(message)
			if err != nil {
				r.logger.Printf("%s: %s", "Failed to process image", err)
			}

			err = d.Ack(false)
			if err != nil {
				r.logger.Printf("%s: %s", "Failed confirmation message", err)
			}
			r.logger.Printf("%s: %s", "Acknowledged message", d.Body)
		}()
		if err != nil {
			return err
		}
	}
	<-forever
	return nil
}

// QosQueue controls messages.
func (r *RabbitMQ) QosQueue() error {
	err := r.ch.Qos(
		7,
		0,
		false,
	)
	if err != nil {
		r.logger.Fatalf("%s: %s", "Failed qos", err)
	}
	return nil
}
