package apiserver

import (
	"context"
	"fmt"
	_ "image/jpeg" // It allows using jpeg
	_ "image/png"  // It allows using png
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/alisavch/image-service/internal/models"
	"github.com/alisavch/image-service/internal/service"
	"github.com/alisavch/image-service/internal/utils"

	"github.com/google/uuid"
)

type key string

const (
	authorizationHeader key = "Authorization"
	userCtx             key = "userId"
	aws                     = "AWS"
	local                   = "local"
)

// Request is an interface which must be implemented by request models.
type Request interface {
	Build(*http.Request) error
	Validate() error
}

// ParseRequest parses request from http Request, stores it in the value pointed to by req and validates it.
func ParseRequest(r *http.Request, req Request) error {
	err := req.Build(r)
	if err != nil {
		return err
	}
	return req.Validate()
}

func (s *Server) healthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.respondJSON(w, http.StatusOK, "Healthy")
			s.logger.Printf("%s:%s", w, "OK")
		default:
			s.respondJSON(w, http.StatusMethodNotAllowed, "Unhealthy")
		}
	}
}

type authorization struct {
	header      string
	headerParts []string
	token       string
}

// Build builds a request to authorize.
func (req *authorization) Build(r *http.Request) error {
	req.header = r.Header.Get(string(authorizationHeader))
	if req.header == "" {
		return utils.ErrEmptyHeader
	}
	req.headerParts = strings.Split(req.header, " ")
	req.token = req.headerParts[1]
	return nil
}

// Validate validates request to authorize.
func (req authorization) Validate() error {
	if len(req.headerParts) != 2 || req.headerParts[0] != "Bearer" {
		return utils.ErrInvalidAuthHeader
	}
	if req.token == "" {
		return utils.ErrEmptyToken
	}
	return nil
}

func (s *Server) authorize(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req authorization

		err := ParseRequest(r, &req)
		if err != nil {
			s.errorJSON(w, http.StatusUnauthorized, err)
		}

		userID, err := s.service.ParseToken(req.token)
		if err != nil {
			s.errorJSON(w, http.StatusUnauthorized, err)
			return
		}
		ctx := context.WithValue(r.Context(), userCtx, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

type uploaded struct {
	file    multipart.File
	handler *multipart.FileHeader
}

// Build builds a request to load an image.
func (req *uploaded) Build(r *http.Request) error {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return err
	}

	req.file, req.handler, err = r.FormFile("uploadFile")
	if err != nil {
		return err
	}

	return nil
}

// Validate builds a request to validate the upload of an image.
func (req uploaded) Validate() error {
	if !(req.handler.Header["Content-Type"][0] == "image/jpeg" || req.handler.Header["Content-Type"][0] == "image/png") {
		return utils.ErrAllowedFormat
	}
	return nil
}

func (s *Server) uploadImage(r *http.Request, uploadedImage models.Image) (models.Image, error) {
	var req uploaded
	err := ParseRequest(r, &req)
	if err != nil {
		return models.Image{}, err
	}

	req.handler.Filename = strings.ReplaceAll(uuid.New().String(), "-", "") + req.handler.Filename

	conf := utils.NewConfig()
	switch conf.Storage {
	case aws:
		uploadedImage, err := s.uploadImageToAWS(r, req)
		if err != nil {
			return models.Image{}, err
		}

		return uploadedImage, nil

	case local:
		uploadedImage, err := s.uploadImageLocally(r, req)
		if err != nil {
			return models.Image{}, err
		}

		return uploadedImage, nil
	}

	return uploadedImage, nil
}

func (s *Server) uploadImageLocally(r *http.Request, req uploaded) (models.Image, error) {
	err := s.prepareImageForLocalLoading(req)
	if err != nil {
		return models.Image{}, err
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return models.Image{}, err
	}

	uploadedImage := fillInTheUploadedImageNameAndLocation(req.handler.Filename, currentDir+"/uploads/")

	uploadedID, err := s.service.ServiceOperations.UploadImage(r.Context(), uploadedImage)
	if err != nil {
		return models.Image{}, fmt.Errorf("%s:%s", utils.ErrUpload, err)
	}
	uploadedImage.ID = uploadedID

	return uploadedImage, nil
}

func (s *Server) uploadImageToAWS(r *http.Request, req uploaded) (models.Image, error) {
	imageLocation, err := s.service.UploadToS3Bucket(req.file, req.handler.Filename)
	if err != nil {
		return models.Image{}, err
	}

	err = req.file.Close()
	if err != nil {
		return models.Image{}, err
	}

	uploadedImage := fillInTheUploadedImageNameAndLocation(req.handler.Filename, imageLocation)

	uploadedID, err := s.service.ServiceOperations.UploadImage(r.Context(), uploadedImage)
	if err != nil {
		return models.Image{}, fmt.Errorf("%s:%s", utils.ErrUpload, err)
	}
	uploadedImage.ID = uploadedID

	return uploadedImage, nil
}

func (s *Server) prepareImageForLocalLoading(req uploaded) error {
	err := service.EnsureBaseDir("./uploads/")
	if err != nil {
		return nil
	}

	out, err := os.Create(fmt.Sprintf("./uploads/%s", req.handler.Filename))
	if err != nil {
		return utils.ErrCreateFile
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			s.logger.Fatalf("%s:%s", "failed fileReader.Close", err)
		}
	}(out)

	_, err = io.Copy(out, req.file)
	if err != nil {
		return utils.ErrCopyFile
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			s.logger.Fatalf("%s:%s", "failed fileReader.Close", err)
		}
	}(req.file)

	return nil
}

func fillInTheUploadedImageNameAndLocation(name, location string) models.Image {
	var uploadedImage models.Image
	uploadedImage.UploadedName = name
	uploadedImage.UploadedLocation = location
	return uploadedImage
}
