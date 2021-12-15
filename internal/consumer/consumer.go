package consumer

import (
	"database/sql"

	"github.com/alisavch/image-service/internal/bucket"
	"github.com/alisavch/image-service/internal/service"

	"github.com/alisavch/image-service/internal/broker"
	"github.com/alisavch/image-service/internal/repository"
	"github.com/alisavch/image-service/internal/utils"
)

// ConversionService contains interfaces.
type ConversionService struct {
	AMQP
}

// NewConversionService configures conversion service.
func NewConversionService(amqp AMQP) *ConversionService {
	return &ConversionService{
		AMQP: amqp,
	}
}

// Consume starts the message consumer.
func Consume() {
	logger := NewLogger()
	conf := utils.NewConfig()

	db, err := repository.NewDB(conf.DBConfig)
	if err != nil {
		logger.Fatalf("%s: %s", "Failed to initialize database", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Fatalf("%s: %s", "Failed to close database", err)
		}
	}(db)
	repos := repository.NewRepository(db)
	aws := bucket.NewAWS()
	services := service.NewService(repos, aws)
	rabbit := broker.NewAMQPBrokerConsumer(services, aws)

	currentService := NewConversionService(rabbit)

	err = currentService.Connect()
	if err != nil {
		logger.Fatalf("%s: %s", "Failed to open a channel", err)
	}

	q, err := currentService.DeclareQueue("publisher")
	if err != nil {
		logger.Fatalf("%s: %s", "Failed to declare a queue", err)
	}

	err = currentService.QosQueue()
	if err != nil {
		logger.Fatalf("%s: %s", "Failed to set qos parameters", err)
	}

	errorChan := make(chan error)

	err = currentService.ConsumeQueue(q.Name, errorChan)
	if err != nil {
		logger.Fatalf("%s: %s", "Failed to consume a queue", err)
	}

	close(errorChan)
}
