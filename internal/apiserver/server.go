package apiserver

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/alisavch/image-service/internal/models"

	"github.com/alisavch/image-service/internal/broker"

	"github.com/alisavch/image-service/internal/service"
	"github.com/gorilla/mux"
)

// Server are complex of routers and services.
type Server struct {
	router  *mux.Router
	service *service.Service
	mq      *broker.AMQPBroker
}

// NewServer configures server.
func NewServer(service *service.Service, mq *broker.AMQPBroker) *Server {
	s := &Server{
		router:  mux.NewRouter(),
		service: service,
		mq:      mq,
	}
	s.ConfigureRouter()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) errorJSON(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respondJSON(w, r, code, map[string]string{"error": err.Error()})
}

func (s *Server) respondJSON(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func (s *Server) respondFormData(w http.ResponseWriter, r *http.Request, code int, id int) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()
	w.Header().Set("Content-Type", writer.FormDataContentType())
	s.respondJSON(w, r, code, map[string]int{"Image ID": id})
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
