package apiserver

import (
	"net/http"

	"github.com/alisavch/image-service/internal/broker"
	"github.com/alisavch/image-service/internal/model"
	"github.com/alisavch/image-service/internal/repository"
	"github.com/alisavch/image-service/internal/service"
	"github.com/alisavch/image-service/internal/utils"
	_ "github.com/lib/pq" // Registers database.
	"github.com/sirupsen/logrus"
)

// StartClient starts the server.
func StartClient() error {
	conf := utils.NewConfig(".env")
	user, pass, host, port, dbname, err := utils.GetDBEnvironments(conf)
	if err != nil {
		logrus.Fatalf("error find variables :%s", err.Error())
	}
	db, err := newDB(user, pass, host, port, dbname)
	if err != nil {
		logrus.Fatalf("error initialize database: %s", err.Error())
	}
	defer db.Close()
	repos := repository.NewRepository(db)
	services := service.NewService(repos)

	rabbit := new(broker.RabbitMQ)
	if err = rabbit.Connect(); err != nil {
		logrus.Fatalf("rabbit connection: %s", err)
	}
	defer rabbit.Close()

	if err = rabbit.DeclareQueue(model.Queued); err != nil {
		logrus.Fatalf("declare queue: %s", err)
	}

	if err = rabbit.QosQueue(); err != nil {
		logrus.Fatalf("qos: %s", err)
	}

	proc := make(chan []byte)
	err = rabbit.ConsumeQueue(model.Processing, proc)
	if err != nil {
		logrus.Fatalf("consume: %s", err)
	}

	srv := NewServer(services, rabbit)

	return http.ListenAndServe(
		":8081",
		srv,
	)
}
