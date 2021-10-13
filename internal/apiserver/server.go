package apiserver

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/alisavch/image-service/internal/log"

	"github.com/google/uuid"

	"github.com/alisavch/image-service/internal/models"

	"github.com/gorilla/mux"
)

// Logger contains methods to display logs.
type Logger struct {
	DisplayLog
}

// NewLogger configures Logger.
func NewLogger() *Logger {
	return &Logger{log.NewLogger()}
}

// Service combines the interfaces for interaction with the service.
type Service struct {
	ServiceOperations
}

// NewService configures Service.
func NewService(operations ServiceOperations) *Service {
	return &Service{
		ServiceOperations: operations,
	}
}

// Server combines the basic constructs to run a server.
type Server struct {
	router  *mux.Router
	mq      AMQP
	service *Service
	logger  *Logger
}

// NewServer configures Server.
func NewServer(mq AMQP, service *Service) *Server {
	s := &Server{
		router:  mux.NewRouter(),
		mq:      mq,
		service: service,
		logger:  NewLogger(),
	}
	s.ConfigureRouter()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) errorJSON(w http.ResponseWriter, code int, err error) {
	s.respondJSON(w, code, map[string]string{"error": err.Error()})
}

func (s *Server) respondJSON(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func (s *Server) respondFormData(w http.ResponseWriter, code int, id uuid.UUID) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer func(writer *multipart.Writer) {
		err := writer.Close()
		if err != nil {
			s.logger.Fatalf("%s:%s", "failed fileReader.Close", err)
		}
	}(writer)
	w.Header().Set("Content-Type", writer.FormDataContentType())
	s.respondJSON(w, code, map[string]uuid.UUID{"Image ID": id})
}

func (s *Server) respondImage(w http.ResponseWriter, image *models.Image) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Disposition", "attachment; filename="+image.Filename)
	w.Header().Set("Content-Type", image.ContentType)
	w.Header().Set("Content-Length", strconv.FormatInt(image.Filesize, 10))
	_, err := io.Copy(w, image.File)
	if err != nil {
		return
	}
}
