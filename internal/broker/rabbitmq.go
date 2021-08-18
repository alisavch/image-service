package broker

import (
	"fmt"

	"github.com/sirupsen/logrus"

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
	url, err := utils.GetRabbitMQURL(utils.NewConfig(".env"))
	if err != nil {
		return err
	}

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
		return fmt.Errorf("declare queue error: %w", err)
	}
	return nil
}

// BindQueue binds an exchange to queue.
func (r *RabbitMQ) BindQueue(name string, key string) (err error) {
	err = r.ch.QueueBind(name, key, "image", false, nil)
	if err != nil {
		return fmt.Errorf("bind queue error: %w", err)
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
func (r *RabbitMQ) ConsumeQueue(queue model.Status, message chan []byte) (err error) {
	deliveries, err := r.ch.Consume(string(queue), "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume queue error: %w", err)
	}
	go func(deliveries <-chan amqp.Delivery, done chan error, message chan []byte) {
		for d := range deliveries {
			logrus.Printf("Received a message: %s", d.Body)
			message <- d.Body
			logrus.Printf("Done")

			err = d.Ack(false)
			if err != nil {
				logrus.Printf("Error ack: %s", err)
			}
		}
		done <- nil
	}(deliveries, r.done, message)
	return nil
}

// QosQueue controls messages.
func (r *RabbitMQ) QosQueue() error {
	err := r.ch.Qos(
		1,
		0,
		false,
	)
	if err != nil {
		return fmt.Errorf("qos :%s", err)
	}
	return nil
}

// Close closes requests.
func (r *RabbitMQ) Close() (err error) {
	err = r.conn.Close()
	if err != nil {
		return fmt.Errorf("close error: %w", err)
	}
	<-r.done
	return nil
}
