package broker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alisavch/image-service/internal/models"
	"github.com/alisavch/image-service/internal/utils"

	"github.com/streadway/amqp"
)

// RabbitMQ is client with RabbitMQ extensions.
type RabbitMQ struct {
	conn     *amqp.Connection
	ch       *amqp.Channel
	stopChan chan bool
}

// NewRabbitMQ configures RabbitMQ.
func NewRabbitMQ() *RabbitMQ {
	return &RabbitMQ{}
}

// ProcessMessage configures and processes messages.
type ProcessMessage struct {
	*ImageService
	*RabbitMQ
	repeater Repeater
	logger   *Logger
}

// NewProcessMessageConsumer configures ProcessMessage for consumer.
func NewProcessMessageConsumer(service *ImageService) *ProcessMessage {
	return &ProcessMessage{ImageService: service, logger: NewLogger(), repeater: NewRepeater(NewBackoff(100*time.Millisecond, 10*time.Second, 2, nil), nil), RabbitMQ: NewRabbitMQ()}
}

// NewProcessMessageAPI configures ProcessMessage for API.
func NewProcessMessageAPI() *ProcessMessage {
	return &ProcessMessage{logger: NewLogger(), RabbitMQ: NewRabbitMQ()}
}

// Connect instantiates the RabbitMQ instances using configuration defined in environment variables.
func (process *ProcessMessage) Connect() error {
	conf := utils.NewConfig()
	var err error
	process.RabbitMQ.conn, err = amqp.Dial(conf.Rabbitmq.RabbitmqURL)
	if err != nil {
		process.logger.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	process.RabbitMQ.ch, err = process.RabbitMQ.conn.Channel()
	if err != nil {
		process.logger.Fatalf("%s: %s", "Failed to open a channel", err)
	}

	process.RabbitMQ.stopChan = make(chan bool)

	return nil
}

// Publish sends data to the queue.
func (process *ProcessMessage) Publish(exchange, key string, message models.QueuedMessage) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("can't marshal queue message: %w", err)
	}

	err = process.RabbitMQ.ch.Publish(exchange, key, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		process.logger.Fatalf("%s:%s", "Failed to publish a message", err)
	}
	return nil
}

// DeclareQueue declares a queue.
func (process *ProcessMessage) DeclareQueue(name string) (amqp.Queue, error) {
	q, err := process.RabbitMQ.ch.QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		process.logger.Fatalf("%s: %s", "Failed to declare a queue", err)
	}
	return q, nil
}

// QosQueue controls messages.
func (process *ProcessMessage) QosQueue() error {
	err := process.RabbitMQ.ch.Qos(
		1,
		0,
		false,
	)
	if err != nil {
		process.logger.Fatalf("%s: %s", "Failed qos", err)
	}
	return nil
}

// ConsumeQueue starts delivering queued messages.
func (process *ProcessMessage) ConsumeQueue(queue string, errorChan chan error) error {
	var message models.QueuedMessage
	deliveries, err := process.RabbitMQ.ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("%s: %s", "Failed to register consumer", err)
	}

	for {
		select {
		case d := <-deliveries:
			go func() {
				process.ConsumeOne(d, message, errorChan)
			}()
		case err := <-errorChan:
			process.logger.Errorf("%s:%s", "An error occurred while consuming", err)
		case <-process.stopChan:
			return nil
		}
	}
}

// ConsumeOne consumes one message
func (process *ProcessMessage) ConsumeOne(d amqp.Delivery, message models.QueuedMessage, errorsChan chan error) {
	if len(d.Body) == 0 {
		err := d.Nack(false, false)
		if err != nil {
			process.logger.Errorf("%s: %s", "Could not nack message", err)
		}
		errorsChan <- utils.ErrReceivedEmpty
		return
	}

	err := json.NewDecoder(bytes.NewReader(d.Body)).Decode(&message)
	if err != nil {
		process.logger.Errorf("%s: %s", "Failed to decode json", err)
		err := d.Nack(false, false)
		if err != nil {
			process.logger.Errorf("%s: %s", "Could not nack message", err)
		}
		errorsChan <- err
		return
	}

	err = process.RunRepeater(context.Background(), message)
	if err != nil {
		process.logger.Errorf("%s: %s", "Failed to process message", err)
		errorsChan <- err

		err := process.ImageService.UpdateStatus(context.Background(), message.RequestID, models.Failed)
		if err != nil {
			process.logger.Errorf("%s: %s", "Failed to update message status", err)
			err := d.Nack(false, false)
			if err != nil {
				process.logger.Errorf("%s: %s", "Could not nack message", err)
			}
		}

		err = d.Nack(false, false)
		if err != nil {
			process.logger.Errorf("%s: %s", "Could not nack message", err)
		}

		return
	}

	process.logger.Printf("%s: %s", "Successfully processed message", message.ID)
	err = d.Ack(false)
	if err != nil {
		process.logger.Printf("%s: %s", "Could not ack message", err)
		return
	}
}

// RunRepeater starts a processing function with limiting access attempts.
func (process *ProcessMessage) RunRepeater(ctx context.Context, message models.QueuedMessage) error {
	defer process.repeater.backoff.Reset()
	for {
		err := process.Process(message)

		switch process.repeater.retryPolicy(err) {
		case Succeed, Fail:
			return err
		case Retry:
			var delay time.Duration
			if delay = process.repeater.backoff.Next(); delay == Stop {
				return err
			}
			timeout := time.After(delay)
			if err := process.repeater.sleep(ctx, timeout); err != nil {
				return err
			}
		}
	}
}
