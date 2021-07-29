package apiserver

import (
	"net/http"
)

// ConfigureRouter registers a couple of URL paths and handlers.
func (s *Server) ConfigureRouter() {
	s.router.HandleFunc("/api/sign-up", s.signUp()).Methods(http.MethodPost)
	s.router.HandleFunc("/api/sign-in", s.signIn()).Methods(http.MethodPost)
	s.router.HandleFunc("/api/{userID}/history", s.authorize(s.findUserHistory())).Methods(http.MethodGet)
	s.router.HandleFunc("/api/{userID}/compress", s.authorize(s.compressImage())).Methods(http.MethodPost)
	s.router.HandleFunc("/api/{userID}/compress/{id}", s.authorize(s.findCompressedImage())).Methods(http.MethodGet)
	s.router.HandleFunc("/api/{userID}/convert", s.authorize(s.convertImage())).Methods(http.MethodPost)
	s.router.HandleFunc("/api/{userID}/convert/{id}", s.authorize(s.findConvertedImage())).Methods(http.MethodGet)
}
