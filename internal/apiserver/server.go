package apiserver

import (
	"encoding/json"
	"net/http"

	"github.com/alisavch/image-service/internal/service"
	"github.com/gorilla/mux"
)

// Server are complex of routers and services.
type Server struct {
	router  *mux.Router
	service *service.Service
}

func newServer(service *service.Service) *Server {
	s := &Server{
		router:  mux.NewRouter(),
		service: service,
	}
	s.ConfigureRouter()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *Server) renderJSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(data); err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}
}

func (s *Server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
