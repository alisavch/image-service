package apiserver

import (
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
			s.errorJSON(w, http.StatusNotFound, err)
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

type findImageRequest struct {
	models.Image
	User       models.User
	requestID  uuid.UUID
	isOriginal bool
}

// Build builds a request to find converted image.
func (req *findImageRequest) Build(r *http.Request) error {
	id, ok := r.Context().Value(userCtx).(uuid.UUID)
	if !ok {
		return utils.ErrGetUserID
	}

	req.User.ID = id

	vars := mux.Vars(r)
	convertedImageID, ok := vars["requestID"]
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
func (req findImageRequest) Validate() error {
	return nil
}

func (s *Server) findImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req findImageRequest
		conf := utils.NewConfig()

		err := ParseRequest(r, &req)
		if err != nil {
			s.errorJSON(w, http.StatusUnauthorized, err)
			return
		}

		err = s.service.ServiceOperations.IsAuthenticated(r.Context(), req.User.ID, req.requestID)
		if err != nil {
			s.errorJSON(w, http.StatusForbidden, err)
			return
		}

		if req.isOriginal {
			uploadedImage, err := s.service.ServiceOperations.FindOriginalImage(r.Context(), req.requestID)
			if err != nil {
				s.errorJSON(w, http.StatusNotFound, fmt.Errorf("%s:%s", utils.ErrFindImage, err))
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

		status, err := s.service.ServiceOperations.FindRequestStatus(r.Context(), req.User.ID, req.requestID)
		if status == models.Processing {
			s.errorJSON(w, http.StatusConflict, fmt.Errorf("%s:%s", "cannot get image", utils.ErrImageProcessing))
			return
		}
		if err != nil {
			s.errorJSON(w, http.StatusNotFound, err)
			return
		}

		resultedImage, err := s.service.ServiceOperations.FindResultedImage(r.Context(), req.requestID)
		if err != nil {
			s.errorJSON(w, http.StatusNotFound, fmt.Errorf("%s:%s", utils.ErrFindImage, err))
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

type getStatusRequest struct {
	models.User
	models.RequestStatus
}

// Build builds a request to get status request.
func (req *getStatusRequest) Build(r *http.Request) error {
	id, ok := r.Context().Value(userCtx).(uuid.UUID)
	if !ok {
		return utils.ErrGetUserID
	}

	req.User.ID = id

	vars := mux.Vars(r)
	imageID, ok := vars["requestID"]
	if !ok {
		return utils.ErrMissingParams
	}

	ID := uuid.MustParse(imageID)
	req.RequestID = ID

	return nil
}

// Validate validates request to get status request.
func (req getStatusRequest) Validate() error {
	return nil
}

func (s *Server) findStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req getStatusRequest

		err := ParseRequest(r, &req)
		if err != nil {
			s.errorJSON(w, http.StatusUnauthorized, err)
			return
		}

		err = s.service.ServiceOperations.IsAuthenticated(r.Context(), req.User.ID, req.RequestID)
		if err != nil {
			s.errorJSON(w, http.StatusForbidden, err)
			return
		}

		status, err := s.service.ServiceOperations.FindRequestStatus(r.Context(), req.User.ID, req.RequestID)
		if err != nil {
			s.errorJSON(w, http.StatusNotFound, err)
			return
		}

		req.Status = status

		s.respondJSON(w, http.StatusOK, req.RequestStatus)
	}
}
