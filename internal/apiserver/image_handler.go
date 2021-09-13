package apiserver

import (
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/alisavch/image-service/internal/models"
	"github.com/alisavch/image-service/internal/utils"

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
	id, ok := r.Context().Value(userCtx).(int)
	if !ok {
		return utils.ErrFailedConvert
	}

	vars := mux.Vars(r)
	paramID, ok := vars["userID"]
	if !ok {
		return utils.ErrMissingParams
	}

	intParamID, err := strconv.Atoi(paramID)
	if err != nil {
		return utils.ErrAtoi
	}

	if id != intParamID {
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
			s.errorJSON(w, r, http.StatusUnauthorized, err)
			return
		}

		history, err := s.service.Image.FindUserHistoryByID(r.Context(), req.User.ID)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respondJSON(w, r, http.StatusOK, history)
	}
}

type compressImageRequest struct {
	models.UploadedImage
	User  models.User
	Width int
}

// Build builds a request to compress image.
func (req *compressImageRequest) Build(r *http.Request) error {
	id, ok := r.Context().Value(userCtx).(int)
	if !ok {
		return utils.ErrFailedConvert
	}

	vars := mux.Vars(r)
	paramID, ok := vars["userID"]
	if !ok {
		return utils.ErrMissingParams
	}

	intParamID, err := strconv.Atoi(paramID)
	if err != nil {
		return utils.ErrAtoi
	}

	if id != intParamID {
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
			s.errorJSON(w, r, http.StatusUnauthorized, err)
			return
		}

		startOfExecution := time.Now()

		newUploadedImage, err := s.uploadImage(r, req.UploadedImage)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, err)
			return
		}

		q, err := s.mq.DeclareQueue("publisher")
		if err != nil {
			logrus.Fatalf("%s: %s", "Failed to declare a queue", err)
		}

		err = s.mq.QosQueue()
		if err != nil {
			logrus.Fatalf("%s: %s", "Failed to controls messages", err)
		}

		err = s.mq.Publish("", q.Name, string(models.Queued))
		if err != nil {
			logrus.Fatalf("%s: %s", "Failed to publish a message", err)
		}

		err = s.service.Image.UpdateStatus(r.Context(), newUploadedImage.ID, models.Processing)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, err)
			return
		}

		err = s.mq.Publish("", q.Name, string(models.Processing))
		if err != nil {
			logrus.Fatalf("%s: %s", "Failed to publish a message", err)
		}

		resultedImage, err := s.service.Image.CompressImage(req.Width, newUploadedImage)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, err)
			return
		}

		err = s.service.Image.UpdateStatus(r.Context(), newUploadedImage.ID, models.Done)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, err)
			return
		}

		err = s.mq.Publish("", q.Name, string(models.Done))
		if err != nil {
			logrus.Fatalf("%s: %s", "Failed to publish a message", err)
		}

		endOfExecution := time.Now()
		resultedImage.Service = models.Compression

		requestID, err := s.service.CreateRequest(
			r.Context(),
			req.User,
			newUploadedImage,
			resultedImage,
			models.UserImage{
				UserAccountID:   req.User.ID,
				UploadedImageID: newUploadedImage.ID,
				Status:          models.Done},
			models.Request{
				TimeStart: startOfExecution,
				EndOfTime: endOfExecution,
			})
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respondFormData(w, r, http.StatusOK, requestID)
	}
}

type findCompressedImageRequest struct {
	models.ResultedImage
	User       models.User
	requestID  int
	isOriginal bool
}

// Build builds a request to find compressed image.
func (req *findCompressedImageRequest) Build(r *http.Request) error {
	id, ok := r.Context().Value(userCtx).(int)
	if !ok {
		return utils.ErrFailedConvert
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

	intParamID, err := strconv.Atoi(paramID)
	if err != nil {
		return utils.ErrAtoi
	}

	if id != intParamID {
		return utils.ErrPrivacy
	}
	req.User.ID = id

	originalImage := r.FormValue("original")
	if originalImage == "" {
		originalImage = DefaultOriginal
	}

	convertedID, err := strconv.Atoi(compressedImageID)
	if err != nil {
		return utils.ErrAtoi
	}
	convertedBool, err := strconv.ParseBool(originalImage)
	if err != nil {
		return err
	}
	req.requestID = convertedID
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

		err := ParseRequest(r, &req)
		if err != nil {
			s.errorJSON(w, r, http.StatusUnauthorized, err)
			return
		}

		if req.isOriginal {
			uploaded, err := s.service.Image.FindOriginalImage(r.Context(), req.requestID)
			if err != nil {
				s.errorJSON(w, r, http.StatusInternalServerError, utils.ErrFindImage)
				return
			}

			file, err := s.service.Image.SaveImage(uploaded.Name, "/uploads/")
			if err != nil {
				s.errorJSON(w, r, http.StatusInternalServerError, utils.ErrSaveImage)
				return
			}

			s.respondImage(w, file)
			return
		}

		resultedImage, err := s.service.Image.FindTheResultingImage(r.Context(), req.requestID, models.Compression)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, utils.ErrFindImage)
			return
		}

		file, err := s.service.Image.SaveImage(resultedImage.Name, "/results/")
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, utils.ErrSaveImage)
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
	id, ok := r.Context().Value(userCtx).(int)
	if !ok {
		return utils.ErrFailedConvert
	}

	vars := mux.Vars(r)
	paramID, ok := vars["userID"]
	if !ok {
		return utils.ErrMissingParams
	}

	intParamID, err := strconv.Atoi(paramID)
	if err != nil {
		return utils.ErrAtoi
	}

	if id != intParamID {
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
			s.errorJSON(w, r, http.StatusUnauthorized, err)
			return
		}

		startOfExecution := time.Now()

		newUploadedImage, err := s.uploadImage(r, req.UploadedImage)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, err)
			return
		}

		q, err := s.mq.DeclareQueue("publisher")
		if err != nil {
			logrus.Fatalf("%s: %s", "Failed to declare a queue", err)
		}

		err = s.mq.QosQueue()
		if err != nil {
			logrus.Fatalf("%s: %s", "Failed to controls messages", err)
		}

		err = s.mq.Publish("", q.Name, string(models.Queued))
		if err != nil {
			logrus.Fatalf("%s: %s", "Failed to publish a message", err)
		}

		err = s.service.Image.UpdateStatus(r.Context(), newUploadedImage.ID, models.Processing)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, err)
			return
		}

		err = s.mq.Publish("", q.Name, string(models.Processing))
		if err != nil {
			logrus.Fatalf("%s: %s", "Failed to publish a message", err)
		}

		resultedImage, err := s.service.Image.ConvertToType(newUploadedImage)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, utils.ErrConvert)
			return
		}

		err = s.service.Image.UpdateStatus(r.Context(), newUploadedImage.ID, models.Done)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, err)
			return
		}

		err = s.mq.Publish("", q.Name, string(models.Done))
		if err != nil {
			logrus.Fatalf("%s: %s", "Failed to publish a message", err)
		}

		endOfExecution := time.Now()
		resultedImage.Service = models.Conversion

		requestID, err := s.service.CreateRequest(
			r.Context(),
			req.User,
			newUploadedImage,
			resultedImage,
			models.UserImage{
				UserAccountID:   req.User.ID,
				UploadedImageID: newUploadedImage.ID,
				Status:          models.Done},
			models.Request{
				TimeStart: startOfExecution,
				EndOfTime: endOfExecution,
			})
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respondFormData(w, r, http.StatusOK, requestID)
	}
}

type findConvertedImageRequest struct {
	models.ResultedImage
	User       models.User
	requestID  int
	isOriginal bool
}

// Build builds a request to find converted image.
func (req *findConvertedImageRequest) Build(r *http.Request) error {
	id, ok := r.Context().Value(userCtx).(int)
	if !ok {
		return utils.ErrFailedConvert
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

	intParamID, err := strconv.Atoi(paramID)
	if err != nil {
		return utils.ErrAtoi
	}

	if id != intParamID {
		return utils.ErrPrivacy
	}
	req.User.ID = id

	convertedID, err := strconv.Atoi(convertedImageID)
	if err != nil {
		return utils.ErrAtoi
	}
	req.requestID = convertedID

	originalImage := r.FormValue("original")
	if originalImage == "" {
		originalImage = DefaultOriginal
	}
	req.isOriginal, err = strconv.ParseBool(originalImage)
	if err != nil {
		return err
	}

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
			s.errorJSON(w, r, http.StatusUnauthorized, err)
			return
		}

		if req.isOriginal {
			uploaded, err := s.service.Image.FindOriginalImage(r.Context(), req.requestID)
			if err != nil {
				s.errorJSON(w, r, http.StatusInternalServerError, utils.ErrFindImage)
				return
			}

			file, err := s.service.Image.SaveImage(uploaded.Name, "/uploads/")
			if err != nil {
				s.errorJSON(w, r, http.StatusInternalServerError, utils.ErrSaveImage)
				return
			}
			s.respondImage(w, file)
			return
		}

		resultedImage, err := s.service.Image.FindTheResultingImage(r.Context(), req.requestID, models.Conversion)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, utils.ErrFindImage)
			return
		}

		file, err := s.service.Image.SaveImage(resultedImage.Name, "/results/")
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, utils.ErrSaveImage)
			return
		}

		s.respondImage(w, file)
	}
}
