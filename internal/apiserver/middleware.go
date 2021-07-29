package apiserver

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/alisavch/image-service/internal/model"
	"github.com/alisavch/image-service/internal/utils"
	"github.com/google/uuid"
)

type  key string

const (
	authorizationHeader key = "Authorization"
	userCtx key             = "userId"
)

func (s *Server) authorize(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(string(authorizationHeader))
		if header == "" {
			s.error(w, r, http.StatusUnauthorized, utils.ErrEmptyHeader)
			return
		}
		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			s.error(w, r, http.StatusUnauthorized, utils.ErrInvalidHeader)
			return
		}
		if len(headerParts[1]) == 0 {
			s.error(w, r, http.StatusUnauthorized, utils.ErrEmptyToken)
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

func (s *Server) getUserID(r *http.Request) (int, error) {
	id, ok := r.Context().Value(userCtx).(int)
	if !ok {
		return 0, utils.ErrFailedConvert
	}
	return id, nil
}

func (s *Server) uploadImage(r *http.Request, uploadedImage model.UploadedImage) (model.UploadedImage, error) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadFile")
	if err != nil {
		return model.UploadedImage{}, utils.ErrUpload
	}
	defer file.Close()
	handler.Filename = strings.Replace(uuid.New().String(), "-", "", -1) + handler.Filename
	out, err := os.Create(fmt.Sprintf("./uploads/%s", handler.Filename))
	if err != nil {
		return model.UploadedImage{}, utils.ErrCreateFile
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		return model.UploadedImage{}, utils.ErrCopyFile
	}
	if !(handler.Header["Content-Type"][0] == "image/jpeg" || handler.Header["Content-Type"][0] == "image/png") {
		return model.UploadedImage{}, utils.ErrAllowedFormat
	}

	uploadedImage.Name = handler.Filename
	currentDir, _ := os.Getwd()
	uploadedImage.Location = currentDir + "\\uploaded\\"
	uploadedID, err := s.service.Image.UploadImage(uploadedImage)
	if err != nil {
		return model.UploadedImage{}, utils.ErrUpload
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
	if isOriginal {
		uploaded, err := s.service.Image.FindOriginalImage(id)
		if err != nil {
			return utils.ErrFindImage
		}
		err = s.service.Image.SaveImage(uploaded.Name, "\\uploads\\", "orig"+uploaded.Name)
		if err != nil {
			return utils.ErrSaveImage
		}
	}
	return nil
}
