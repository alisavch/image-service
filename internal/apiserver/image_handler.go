package apiserver

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alisavch/image-service/internal/models"
	"github.com/alisavch/image-service/internal/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	// DefaultWidth is default value for compress JPEG and PNG.
	DefaultWidth = "150"
	// DefaultOriginal is default value for downloads original image.
	DefaultOriginal = "false"
)

type findUserHistoryRequest struct {
	models.User
}

// Build builds a request to find history.
func (req *findUserHistoryRequest) Build(r *http.Request) error {
	id, ok := r.Context().Value(userCtx).(uuid.UUID)
	if !ok {
		return utils.ErrGetUserID
	}

	req.User.ID = id

	return nil
}

// Validate validates request to find history.
func (req findUserHistoryRequest) Validate() error {
	return nil
}

func (s *Server) findUserHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req findUserHistoryRequest

		err := ParseRequest(r, &req)
		if err != nil {
			s.errorJSON(w, http.StatusUnauthorized, err)
			return
		}

		history, err := s.service.ServiceOperations.FindUserRequestHistory(r.Context(), req.User.ID)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, err)
			return
		}

		s.respondJSON(w, http.StatusOK, history)
	}
}

type compressImageRequest struct {
	models.Image
	User         models.User
	ImageRequest models.Request
	Width        int
}

// Build builds a request to compress image.
func (req *compressImageRequest) Build(r *http.Request) error {
	id, ok := r.Context().Value(userCtx).(uuid.UUID)
	if !ok {
		return utils.ErrGetUserID
	}

	req.User.ID = id

	width := r.FormValue("width")
	if width == "" {
		width = DefaultWidth
	}

	convertedWidth, err := strconv.Atoi(width)
	if err != nil {
		return utils.ErrAtoi
	}
	req.Width = convertedWidth
	req.ImageRequest.Status = models.Queued
	req.ImageRequest.ServiceName = models.Compression

	return nil
}

// Validate validates request to compress image.
func (req compressImageRequest) Validate() error {
	return nil
}

func (s *Server) compressImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req compressImageRequest
		err := ParseRequest(r, &req)
		if err != nil {
			s.errorJSON(w, http.StatusUnauthorized, err)
			return
		}

		originalImage, err := s.uploadImage(r, req.Image)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, err)
			return
		}
		s.logger.Printf("%s:%s", "Original image uploaded", originalImage.ID)

		req.Image.ID = originalImage.ID
		req.Image.UploadedName = originalImage.UploadedName
		req.Image.UploadedLocation = originalImage.UploadedLocation
		requestID, err := s.service.ServiceOperations.CreateRequest(r.Context(), req.User, req.Image, req.ImageRequest)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, err)
			return
		}
		s.logger.Printf("%s:%s", "Request created", requestID)

		q, err := s.mq.DeclareQueue("publisher")
		if err != nil {
			s.logger.Fatalf("%s: %s", "Failed to declare a queue", err)
		}

		err = s.mq.QosQueue()
		if err != nil {
			s.logger.Fatalf("%s: %s", "Failed to set qos parameters", err)
		}

		err = s.service.ServiceOperations.UpdateStatus(r.Context(), requestID, models.Processing)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, err)
			return
		}
		s.logger.Printf("%s:%s", "Status updated", models.Processing)

		message := models.NewQueuedMessage(req.Width, requestID, req.ImageRequest.ServiceName, originalImage)

		err = s.mq.Publish("", q.Name, message)
		if err != nil {
			s.logger.Fatalf("%s: %s", "Failed to publish a message", err)
		}
		s.logger.Printf("%s:%s", "Message sent", message.Service)

		s.respondFormData(w, http.StatusAccepted, requestID)
	}
}

type findCompressedImageRequest struct {
	models.Image
	User       models.User
	requestID  uuid.UUID
	isOriginal bool
}

// Build builds a request to find compressed image.
func (req *findCompressedImageRequest) Build(r *http.Request) error {
	id, ok := r.Context().Value(userCtx).(uuid.UUID)
	if !ok {
		return utils.ErrGetUserID
	}

	req.User.ID = id

	vars := mux.Vars(r)
	compressedImageID, ok := vars["compressedID"]
	if !ok {
		return utils.ErrMissingParams
	}

	originalImage := r.FormValue("original")
	if originalImage == "" {
		originalImage = DefaultOriginal
	}

	convertedBool, err := strconv.ParseBool(originalImage)
	if err != nil {
		return err
	}

	compressedID := uuid.MustParse(compressedImageID)

	req.requestID = compressedID
	req.isOriginal = convertedBool

	return nil
}

// Validate validates request to find compressed image.
func (req findCompressedImageRequest) Validate() error {
	return nil
}

func (s *Server) findCompressedImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req findCompressedImageRequest
		conf := utils.NewConfig()

		err := ParseRequest(r, &req)
		if err != nil {
			s.errorJSON(w, http.StatusUnauthorized, err)
			return
		}

		if req.isOriginal {
			uploadedImage, err := s.service.ServiceOperations.FindOriginalImage(r.Context(), req.requestID)
			if err != nil {
				s.errorJSON(w, http.StatusInternalServerError, fmt.Errorf("%s:%s", utils.ErrFindImage, err))
				return
			}
			s.logger.Printf("%s:%s", "Original image found", uploadedImage.UploadedName)

			file, err := s.service.ServiceOperations.SaveImage(uploadedImage.UploadedName, uploadedImage.UploadedLocation, conf.Storage)
			if err != nil {
				s.errorJSON(w, http.StatusInternalServerError, fmt.Errorf("%s:%s", utils.ErrSaveImage, err))
				return
			}
			s.logger.Printf("%s:%s", "Original image received", uploadedImage.UploadedName)

			s.respondImage(w, file)
			return
		}

		err = s.service.ServiceOperations.CheckStatus(r.Context(), req.requestID)
		if errors.Is(err, utils.ErrImageProcessing) {
			s.errorJSON(w, http.StatusNotFound, err)
			return
		}
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, err)
			return
		}

		resultedImage, err := s.service.ServiceOperations.FindResultedImage(r.Context(), req.requestID)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, fmt.Errorf("%s:%s", utils.ErrFindImage, err))
			return
		}
		s.logger.Printf("%s:%s", "Resulted image found", resultedImage.ResultedName)

		file, err := s.service.ServiceOperations.SaveImage(resultedImage.ResultedName, resultedImage.ResultedLocation, conf.Storage)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, fmt.Errorf("%s:%s", utils.ErrSaveImage, err))
			return
		}
		s.logger.Printf("%s:%s", "Resulted image received", resultedImage.ResultedName)

		s.respondImage(w, file)
	}
}

type convertImageRequest struct {
	models.Image
	User         models.User
	ImageRequest models.Request
}

// Build builds a request to convert image.
func (req *convertImageRequest) Build(r *http.Request) error {
	id, ok := r.Context().Value(userCtx).(uuid.UUID)
	if !ok {
		return utils.ErrGetUserID
	}

	req.User.ID = id
	req.ImageRequest.Status = models.Queued
	req.ImageRequest.ServiceName = models.Conversion

	return nil
}

// Validate validates request to convert image.
func (req convertImageRequest) Validate() error {
	return nil
}

func (s *Server) convertImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req convertImageRequest

		err := ParseRequest(r, &req)
		if err != nil {
			s.errorJSON(w, http.StatusUnauthorized, err)
			return
		}

		originalImage, err := s.uploadImage(r, req.Image)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, err)
			return
		}
		s.logger.Printf("%s:%s", "Original image uploaded", originalImage.ID)

		req.Image.ID = originalImage.ID
		req.Image.UploadedName = originalImage.UploadedName
		req.Image.UploadedLocation = originalImage.UploadedLocation
		requestID, err := s.service.ServiceOperations.CreateRequest(r.Context(), req.User, req.Image, req.ImageRequest)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, err)
			return
		}
		s.logger.Printf("%s:%s", "Request created", requestID)

		q, err := s.mq.DeclareQueue("publisher")
		if err != nil {
			s.logger.Fatalf("%s: %s", "Failed to declare a queue", err)
		}

		err = s.mq.QosQueue()
		if err != nil {
			s.logger.Fatalf("%s: %s", "Failed to controls messages", err)
		}

		err = s.service.ServiceOperations.UpdateStatus(r.Context(), requestID, models.Processing)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, err)
			return
		}
		s.logger.Printf("%s:%s", "Status updated", models.Processing)

		message := models.NewQueuedMessage(0, requestID, req.ImageRequest.ServiceName, originalImage)

		err = s.mq.Publish("", q.Name, message)
		if err != nil {
			s.logger.Fatalf("%s: %s", "Failed to publish a message", err)
		}
		s.logger.Printf("%s:%s", "Message sent", message.Service)

		s.respondFormData(w, http.StatusAccepted, requestID)
	}
}

type findConvertedImageRequest struct {
	models.Image
	User       models.User
	requestID  uuid.UUID
	isOriginal bool
}

// Build builds a request to find converted image.
func (req *findConvertedImageRequest) Build(r *http.Request) error {
	id, ok := r.Context().Value(userCtx).(uuid.UUID)
	if !ok {
		return utils.ErrGetUserID
	}

	req.User.ID = id

	vars := mux.Vars(r)
	convertedImageID, ok := vars["convertedID"]
	if !ok {
		return utils.ErrRequest
	}

	originalImage := r.FormValue("original")
	if originalImage == "" {
		originalImage = DefaultOriginal
	}

	convertedBool, err := strconv.ParseBool(originalImage)
	if err != nil {
		return err
	}

	convertedID := uuid.MustParse(convertedImageID)

	req.requestID = convertedID
	req.isOriginal = convertedBool

	return nil
}

// Validate validates request to find converted image.
func (req findConvertedImageRequest) Validate() error {
	return nil
}

func (s *Server) findConvertedImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req findConvertedImageRequest
		conf := utils.NewConfig()

		err := ParseRequest(r, &req)
		if err != nil {
			s.errorJSON(w, http.StatusUnauthorized, err)
			return
		}

		if req.isOriginal {
			uploadedImage, err := s.service.ServiceOperations.FindOriginalImage(r.Context(), req.requestID)
			if err != nil {
				s.errorJSON(w, http.StatusInternalServerError, fmt.Errorf("%s:%s", utils.ErrFindImage, err))
				return
			}
			s.logger.Printf("%s:%s", "Original image found", uploadedImage.UploadedName)

			file, err := s.service.ServiceOperations.SaveImage(uploadedImage.UploadedName, "/uploads/", conf.Storage)
			if err != nil {
				s.errorJSON(w, http.StatusInternalServerError, fmt.Errorf("%s:%s", utils.ErrSaveImage, err))
				return
			}
			s.logger.Printf("%s:%s", "Original image received", uploadedImage.UploadedName)

			s.respondImage(w, file)
			return
		}

		err = s.service.ServiceOperations.CheckStatus(r.Context(), req.requestID)
		if errors.Is(err, utils.ErrImageProcessing) {
			s.errorJSON(w, http.StatusNotFound, err)
			return
		}
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, err)
			return
		}

		resultedImage, err := s.service.ServiceOperations.FindResultedImage(r.Context(), req.requestID)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, fmt.Errorf("%s:%s", utils.ErrFindImage, err))
			return
		}
		s.logger.Printf("%s:%s", "Resulted image found", resultedImage.ResultedName)

		file, err := s.service.ServiceOperations.SaveImage(resultedImage.ResultedName, "/results/", conf.Storage)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, fmt.Errorf("%s:%s", utils.ErrSaveImage, err))
			return
		}
		s.logger.Printf("%s:%s", "Resulted image received", resultedImage.ResultedName)

		s.respondImage(w, file)
	}
}
