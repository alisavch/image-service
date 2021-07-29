package apiserver

import (
	"net/http"

	"github.com/alisavch/image-service/internal/model"
)

func (s *Server) signUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input model.User
		s.renderJSON(w, r, &input)

		id, err := s.service.Authorization.CreateUser(input)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, id)
	}
}

type signInInput struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (s *Server) signIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input signInInput
		s.renderJSON(w, r, &input)
		token, err := s.service.Authorization.GenerateToken(input.Username, input.Password)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, token)
	}
}
