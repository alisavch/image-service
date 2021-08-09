package apiserver

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/alisavch/image-service/internal/model"
	"github.com/google/uuid"
)

type key string

const (
	authorizationHeader key = "Authorization"
	userCtx             key = "userId"
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

func (s *Server) authorize(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(string(authorizationHeader))
		if header == "" {
			s.error(w, r, http.StatusUnauthorized, fmt.Errorf("auth header is empty"))
			return
		}
		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			s.error(w, r, http.StatusUnauthorized, fmt.Errorf("auth header is invalid"))
			return
		}
		if len(headerParts[1]) == 0 {
			s.error(w, r, http.StatusUnauthorized, fmt.Errorf("token is empty"))
		}

		userID, err := s.service.Authorization.ParseToken(headerParts[1])
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}
		ctx := context.WithValue(r.Context(), userCtx, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

type request struct {
	file    multipart.File
	handler *multipart.FileHeader
}

func (req *request) Build(r *http.Request) error {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return fmt.Errorf("error parse multipart from")
	}
	req.file, req.handler, err = r.FormFile("uploadFile")
	if err != nil {
		return fmt.Errorf("error upload file from form")
	}
	defer req.file.Close()
	return nil
}

func (req request) Validate() error {
	return nil
}

func (s *Server) uploadImage(r *http.Request, uploadedImage model.UploadedImage) (model.UploadedImage, error) {
	var req request
	err := ParseRequest(r, &req)
	req.handler.Filename = strings.Replace(uuid.New().String(), "-", "", -1) + req.handler.Filename

	out, err := os.Create(fmt.Sprintf("./uploads/%s", req.handler.Filename))
	if err != nil {
		return model.UploadedImage{}, fmt.Errorf("error create new file")
	}
	defer out.Close()

	_, err = io.Copy(out, req.file)
	if err != nil {
		return model.UploadedImage{}, fmt.Errorf("error copy file")
	}

	if !(req.handler.Header["Content-Type"][0] == "image/jpeg" || req.handler.Header["Content-Type"][0] == "image/png") {
		return model.UploadedImage{}, fmt.Errorf("file format is not allowed. Please upload a JPEG or PNG")
	}

	uploadedImage.Name = req.handler.Filename
	currentDir, _ := os.Getwd()
	uploadedImage.Location = currentDir + "\\uploaded\\"
	uploadedID, err := s.service.Image.UploadImage(r.Context(), uploadedImage)
	if err != nil {
		return model.UploadedImage{}, fmt.Errorf("error upload file")
	}
	uploadedImage.ID = uploadedID

	return uploadedImage, nil
}

func (s *Server) findOriginalImage(r *http.Request, id int) error {
	originalImage := r.FormValue("original")
	if originalImage == "" {
		originalImage = DefaultOriginal
	}
	isOriginal, err := strconv.ParseBool(originalImage)

	if err != nil {
		return err
	}

	if !isOriginal {
		return nil
	}

	uploaded, err := s.service.Image.FindOriginalImage(r.Context(), id)
	if err != nil {
		return fmt.Errorf("cannot find image")
	}

	if err = s.service.Image.SaveImage(uploaded.Name, "\\uploads\\", "orig"+uploaded.Name); err != nil {
		return fmt.Errorf("cannot save image")
	}
	return nil
}
