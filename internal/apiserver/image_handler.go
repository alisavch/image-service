package apiserver

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/alisavch/image-service/internal/service"

	"github.com/alisavch/image-service/internal/models"
	"github.com/alisavch/image-service/internal/utils"

	"github.com/gorilla/mux"
)

const (
	// DefaultWidth is default value for compress JPEG and PNG.
	DefaultWidth = "150"
	// DefaultOriginal is default value for downloads original image.
	DefaultOriginal = "false"
	// IsRemoteStorage is default value for storing images.
	IsRemoteStorage = true
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

	vars := mux.Vars(r)
	paramID, ok := vars["userID"]
	if !ok {
		return utils.ErrMissingParams
	}

	if id.String() != paramID {
		return utils.ErrPrivacy
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

		history, err := s.service.Image.FindUserHistoryByID(r.Context(), req.User.ID)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, err)
			return
		}

		s.respondJSON(w, http.StatusOK, history)
	}
}

type compressImageRequest struct {
	models.UploadedImage
	User  models.User
	Width int
}

// Build builds a request to compress image.
func (req *compressImageRequest) Build(r *http.Request) error {
	id, ok := r.Context().Value(userCtx).(uuid.UUID)
	if !ok {
		return utils.ErrGetUserID
	}

	vars := mux.Vars(r)
	paramID, ok := vars["userID"]
	if !ok {
		return utils.ErrMissingParams
	}

	if id.String() != paramID {
		return utils.ErrPrivacy
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

		startOfExecution := time.Now()

		originalImage, err := s.uploadImage(r, req.UploadedImage)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, err)
			return
		}

		q, err := s.mq.DeclareQueue("publisher")
		if err != nil {
			logger.Fatalf("%s: %s", "Failed to declare a queue", err)
		}

		err = s.mq.QosQueue()
		if err != nil {
			logger.Fatalf("%s: %s", "Failed to controls messages", err)
		}

		err = s.mq.Publish("", q.Name, string(models.Queued))
		if err != nil {
			logger.Fatalf("%s: %s", "Failed to publish a message", err)
		}

		err = s.service.Image.UpdateStatus(r.Context(), originalImage.ID, models.Processing)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, err)
			return
		}

		err = s.mq.Publish("", q.Name, string(models.Processing))
		if err != nil {
			logger.Fatalf("%s: %s", "Failed to publish a message", err)
		}

		resultedName := newImgName("cmp-" + originalImage.Name)

		img, format, file, err := prepareImage(originalImage, originalImage.Name, resultedName)

		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, err)
			return
		}

		resultedImage, err := s.service.Image.CompressImage(req.Width, format, resultedName, img, file, IsRemoteStorage)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, err)
			return
		}

		err = s.service.Image.UpdateStatus(r.Context(), originalImage.ID, models.Done)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, err)
			return
		}

		err = s.mq.Publish("", q.Name, string(models.Done))
		if err != nil {
			logger.Fatalf("%s: %s", "Failed to publish a message", err)
		}

		endOfExecution := time.Now()
		resultedImage.Service = models.Compression

		requestID, err := s.service.CreateRequest(
			r.Context(),
			req.User,
			originalImage,
			resultedImage,
			models.UserImage{
				UserAccountID:   req.User.ID,
				UploadedImageID: originalImage.ID,
				Status:          models.Done,
			},
			models.Request{
				TimeStart: startOfExecution,
				EndOfTime: endOfExecution,
			})
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, err)
			return
		}

		s.respondFormData(w, http.StatusOK, requestID)
	}
}

type findCompressedImageRequest struct {
	models.ResultedImage
	User            models.User
	requestID       uuid.UUID
	isOriginal      bool
	isRemoteStorage bool
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

	paramID, ok := vars["userID"]
	if !ok {
		return utils.ErrMissingParams
	}

	if id.String() != paramID {
		return utils.ErrPrivacy
	}
	req.User.ID = id

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
	req.isRemoteStorage = IsRemoteStorage

	return nil
}

// Validate validates request to find compressed image.
func (req findCompressedImageRequest) Validate() error {
	return nil
}

func (s *Server) findCompressedImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req findCompressedImageRequest

		err := ParseRequest(r, &req)
		if err != nil {
			s.errorJSON(w, http.StatusUnauthorized, err)
			return
		}

		if req.isOriginal {
			uploaded, err := s.service.Image.FindOriginalImage(r.Context(), req.requestID)
			if err != nil {
				s.errorJSON(w, http.StatusInternalServerError, fmt.Errorf("%s:%s", utils.ErrFindImage, err))
				return
			}

			file, err := s.service.Image.SaveImage(uploaded.Name, uploaded.Location, req.isRemoteStorage)
			if err != nil {
				s.errorJSON(w, http.StatusInternalServerError, fmt.Errorf("%s:%s", utils.ErrSaveImage, err))
				return
			}

			s.respondImage(w, file)
			return
		}

		resultedImage, err := s.service.Image.FindTheResultingImage(r.Context(), req.requestID, models.Compression)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, fmt.Errorf("%s:%s", utils.ErrFindImage, err))
			return
		}

		file, err := s.service.Image.SaveImage(resultedImage.Name, resultedImage.Location, req.isRemoteStorage)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, fmt.Errorf("%s:%s", utils.ErrSaveImage, err))
			return
		}

		s.respondImage(w, file)
	}
}

type convertImageRequest struct {
	models.UploadedImage
	User models.User
}

// Build builds a request to convert image.
func (req *convertImageRequest) Build(r *http.Request) error {
	id, ok := r.Context().Value(userCtx).(uuid.UUID)
	if !ok {
		return utils.ErrGetUserID
	}

	vars := mux.Vars(r)
	paramID, ok := vars["userID"]
	if !ok {
		return utils.ErrMissingParams
	}

	if id.String() != paramID {
		return utils.ErrPrivacy
	}

	req.User.ID = id

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

		startOfExecution := time.Now()

		originalImage, err := s.uploadImage(r, req.UploadedImage)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, err)
			return
		}

		q, err := s.mq.DeclareQueue("publisher")
		if err != nil {
			logger.Fatalf("%s: %s", "Failed to declare a queue", err)
		}

		err = s.mq.QosQueue()
		if err != nil {
			logger.Fatalf("%s: %s", "Failed to controls messages", err)
		}

		err = s.mq.Publish("", q.Name, string(models.Queued))
		if err != nil {
			logger.Fatalf("%s: %s", "Failed to publish a message", err)
		}

		err = s.service.Image.UpdateStatus(r.Context(), originalImage.ID, models.Processing)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, err)
			return
		}

		err = s.mq.Publish("", q.Name, string(models.Processing))
		if err != nil {
			logger.Fatalf("%s: %s", "Failed to publish a message", err)
		}

		convertedName, err := service.ChangeFormat(originalImage.Name)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, err)
			return
		}

		resultedName := newImgName("cnv-" + convertedName)

		img, format, file, err := prepareImage(originalImage, originalImage.Name, resultedName)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, err)
			return
		}

		resultedImage, err := s.service.Image.ConvertToType(format, resultedName, img, file, IsRemoteStorage)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, err)
			return
		}

		err = s.service.Image.UpdateStatus(r.Context(), originalImage.ID, models.Done)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, err)
			return
		}

		err = s.mq.Publish("", q.Name, string(models.Done))
		if err != nil {
			logger.Fatalf("%s: %s", "Failed to publish a message", err)
		}

		endOfExecution := time.Now()
		resultedImage.Service = models.Conversion

		requestID, err := s.service.CreateRequest(
			r.Context(),
			req.User,
			originalImage,
			resultedImage,
			models.UserImage{
				UserAccountID:   req.User.ID,
				UploadedImageID: originalImage.ID,
				Status:          models.Done,
			},
			models.Request{
				TimeStart: startOfExecution,
				EndOfTime: endOfExecution,
			})
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, fmt.Errorf("%s:%s", utils.ErrCreateRequest, err))
			return
		}

		s.respondFormData(w, http.StatusOK, requestID)
	}
}

type findConvertedImageRequest struct {
	models.ResultedImage
	User            models.User
	requestID       uuid.UUID
	isOriginal      bool
	isRemoteStorage bool
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

	paramID, ok := vars["userID"]
	if !ok {
		return utils.ErrMissingParams
	}

	if id.String() != paramID {
		return utils.ErrPrivacy
	}
	req.User.ID = id

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
	req.isRemoteStorage = IsRemoteStorage

	return nil
}

// Validate validates request to find converted image.
func (req findConvertedImageRequest) Validate() error {
	return nil
}

func (s *Server) findConvertedImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req findConvertedImageRequest

		err := ParseRequest(r, &req)
		if err != nil {
			s.errorJSON(w, http.StatusUnauthorized, err)
			return
		}

		if req.isOriginal {
			uploaded, err := s.service.Image.FindOriginalImage(r.Context(), req.requestID)
			if err != nil {
				s.errorJSON(w, http.StatusInternalServerError, fmt.Errorf("%s:%s", utils.ErrFindImage, err))
				return
			}

			file, err := s.service.Image.SaveImage(uploaded.Name, "/uploads/", req.isRemoteStorage)
			if err != nil {
				s.errorJSON(w, http.StatusInternalServerError, fmt.Errorf("%s:%s", utils.ErrSaveImage, err))
				return
			}
			s.respondImage(w, file)
			return
		}

		resultedImage, err := s.service.Image.FindTheResultingImage(r.Context(), req.requestID, models.Conversion)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, fmt.Errorf("%s:%s", utils.ErrFindImage, err))
			return
		}

		file, err := s.service.Image.SaveImage(resultedImage.Name, "/results/", req.isRemoteStorage)
		if err != nil {
			s.errorJSON(w, http.StatusInternalServerError, fmt.Errorf("%s:%s", utils.ErrSaveImage, err))
			return
		}

		s.respondImage(w, file)
	}
}
