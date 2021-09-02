package apiserver

import (
	"encoding/json"
	"net/http"

	"github.com/alisavch/image-service/internal/models"

	"github.com/alisavch/image-service/internal/utils"
)

type signUpRequest struct {
	models.User
}

// Build builds a request to sign up.
func (req *signUpRequest) Build(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&req.User)
	if err != nil {
		return err
	}
	return nil
}

// Validate validates request to sign up.
func (req *signUpRequest) Validate() error {
	if req.Username == "" {
		return utils.ErrEmptyUsername
	}
	if req.Password == "" {
		return utils.ErrEmptyPassword
	}
	return nil
}

func (s *Server) signUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req signUpRequest

		err := ParseRequest(r, &req)
		if err != nil {
			s.errorJSON(w, r, http.StatusBadRequest, err)
			return
		}

		id, err := s.service.Authorization.CreateUser(r.Context(), req.User)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respondJSON(w, r, http.StatusOK, id)
	}
}

type signInRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// Build builds a request to sign in.
func (req *signInRequest) Build(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	return nil
}

// Validate validates request to sign in.
func (req *signInRequest) Validate() error {
	if req.Username == "" {
		return utils.ErrEmptyUsername
	}
	if req.Password == "" {
		return utils.ErrEmptyPassword
	}
	return nil
}

func (s *Server) signIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req signInRequest
		err := ParseRequest(r, &req)
		if err != nil {
			s.errorJSON(w, r, http.StatusBadRequest, err)
			return
		}

		token, err := s.service.Authorization.GenerateToken(r.Context(), req.Username, req.Password)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respondJSON(w, r, http.StatusOK, token)
	}
}
