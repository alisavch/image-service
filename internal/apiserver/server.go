package apiserver

import (
	"encoding/json"
	"github.com/alisavch/image-service/internal/broker"
	"net/http"

	"github.com/alisavch/image-service/internal/service"
	"github.com/gorilla/mux"
)

// Server are complex of routers and services.
type Server struct {
	router  *mux.Router
	service *service.Service
	mq *broker.RabbitMQ
}

func newServer(service *service.Service, mq *broker.RabbitMQ) *Server {
	s := &Server{
		router:  mux.NewRouter(),
		service: service,
		mq: mq,
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
