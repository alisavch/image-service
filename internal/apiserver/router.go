package apiserver

import (
	"net/http"
)

// ConfigureRouter registers a couple of URL paths and handlers.
func (s *Server) ConfigureRouter() {
	apiRouter := s.router.PathPrefix("/api").Subrouter()
	userRouter := apiRouter.PathPrefix("/user").Subrouter()
	apiRouter.HandleFunc("/sign-up", s.signUp()).Methods(http.MethodPost)
	apiRouter.HandleFunc("/sign-in", s.signIn()).Methods(http.MethodPost)
	userRouter.HandleFunc("/{userID}/history", s.authorize(s.findUserHistory())).Methods(http.MethodGet)
	userRouter.HandleFunc("/{userID}/compress", s.authorize(s.compressImage())).Methods(http.MethodPost)
	userRouter.HandleFunc("/{userID}/compress/{compressedID}", s.authorize(s.findCompressedImage())).Methods(http.MethodGet)
	userRouter.HandleFunc("/{userID}/convert", s.authorize(s.convertImage())).Methods(http.MethodPost)
	userRouter.HandleFunc("/{userID}/convert/{convertedID}", s.authorize(s.findConvertedImage())).Methods(http.MethodGet)
}
