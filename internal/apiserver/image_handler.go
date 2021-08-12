package apiserver

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"

	"github.com/alisavch/image-service/internal/utils"

	"github.com/alisavch/image-service/internal/model"
	"github.com/gorilla/mux"
)

const (
	// DefaultWidth is default value for compress JPEG and PNG.
	DefaultWidth = "150"
	// DefaultOriginal is default value for downloads original image.
	DefaultOriginal = "false"
)

type findUserHistoryRequest struct {
	model.User
}

// Build builds a request to find history.
func (req *findUserHistoryRequest) Build(r *http.Request) error {
	id, ok := r.Context().Value(userCtx).(int)
	if !ok {
		return utils.ErrFailedConvert
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
			s.errorJSON(w, r, http.StatusBadRequest, err)
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

type uploadImageRequest struct {
	model.UploadedImage
	User  model.User
	Width int
}

// Build builds a request to compress image.
func (req *uploadImageRequest) Build(r *http.Request) error {
	id, ok := r.Context().Value(userCtx).(int)
	if !ok {
		return utils.ErrFailedConvert
	}
	req.User.ID = id

	width := r.FormValue("width")
	if width == "" {
		width = DefaultWidth
	}
	req.Width, _ = strconv.Atoi(width)

	return nil
}

// Validate validates request to compress image.
func (req uploadImageRequest) Validate() error {
	return nil
}

func (s *Server) compressImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req uploadImageRequest

		err := ParseRequest(r, &req)
		if err != nil {
			s.errorJSON(w, r, http.StatusBadRequest, err)
			return
		}

		startOfExecution := time.Now()

		newUploadedImage, err := s.uploadImage(r, req.UploadedImage)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, err)
			return
		}

		err = s.mq.Connect()
		if err != nil {
			logrus.Fatalf("open channel: %s", err)
		}

		err = s.mq.DeclareQueue(model.Processing)
		if err != nil {
			logrus.Fatalf("queue declare: %s", err)
		}

		proc := make(chan []byte)
		err = s.mq.ConsumeQueue(model.Processing, proc)
		if err != nil {
			logrus.Fatalf("consume: %s", err)
		}

		resultedImage, err := s.service.Image.CompressImage(req.Width, newUploadedImage)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, err)
			return
		}

		err = s.mq.DeclareQueue(model.Done)
		if err != nil {
			logrus.Fatalf("queue declare: %s", err)
			return
		}

		done := make(chan []byte)
		err = s.mq.ConsumeQueue(model.Done, done)
		if err != nil {
			logrus.Fatalf("consume: %s", err)
			return
		}

		err = s.mq.Close()
		if err != nil {
			logrus.Fatalf("close: %s", err)
			return
		}

		endOfExecution := time.Now()
		resultedImage.Service = model.Compression

		requestID, err := s.service.CreateRequest(
			r.Context(),
			req.User,
			newUploadedImage,
			resultedImage,
			model.UserImage{
				UserAccountID:   req.User.ID,
				UploadedImageID: newUploadedImage.ID,
				Status:          model.Done},
			model.Request{
				TimeStart: startOfExecution,
				EndOfTime: endOfExecution,
			})
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respondJSON(w, r, http.StatusOK, requestID)
	}
}

type findCompressedImageRequest struct {
	model.ResultedImage
	User         model.User
	CompressedID int
}

// Build builds a request to find compressed image.
func (req *findCompressedImageRequest) Build(r *http.Request) error {
	id, ok := r.Context().Value(userCtx).(int)
	if !ok {
		return fmt.Errorf("failed convert to int userID")
	}
	req.User.ID = id

	vars := mux.Vars(r)
	compressedImageID, ok := vars["compressedID"]
	if !ok {
		return fmt.Errorf("incorrect request")
	}

	req.CompressedID, _ = strconv.Atoi(compressedImageID)
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
			s.errorJSON(w, r, http.StatusBadRequest, err)
			return
		}

		resultedImage, err := s.service.Image.FindTheResultingImage(r.Context(), req.CompressedID, model.Compression)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, utils.ErrFindImage)
			return
		}

		err = s.service.Image.SaveImage(resultedImage.Name, "\\results\\", resultedImage.Name)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, utils.ErrSaveImage)
			return
		}

		err = s.findOriginalImage(r, req.CompressedID)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respondJSON(w, r, http.StatusOK, "successfully saved")
	}
}

type convertImageRequest struct {
	model.UploadedImage
	User model.User
}

// Build builds a request to convert image.
func (req *convertImageRequest) Build(r *http.Request) error {
	id, ok := r.Context().Value(userCtx).(int)
	if !ok {
		return fmt.Errorf("failed convert to int userID")
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
			s.errorJSON(w, r, http.StatusBadRequest, err)
			return
		}

		startOfExecution := time.Now()

		newUploadedImage, err := s.uploadImage(r, req.UploadedImage)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, err)
			return
		}

		err = s.mq.Connect()
		if err != nil {
			logrus.Fatalf("open channel: %s", err)
		}

		err = s.mq.DeclareQueue(model.Processing)
		if err != nil {
			logrus.Fatalf("queue declare: %s", err)
		}

		proc := make(chan []byte)
		err = s.mq.ConsumeQueue(model.Processing, proc)
		if err != nil {
			logrus.Fatalf("consume: %s", err)
		}

		resultedImage, err := s.service.Image.ConvertToType(newUploadedImage)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, utils.ErrConvert)
			return
		}

		err = s.mq.DeclareQueue(model.Done)
		if err != nil {
			logrus.Fatalf("queue declare: %s", err)
			return
		}

		done := make(chan []byte)
		err = s.mq.ConsumeQueue(model.Done, done)
		if err != nil {
			logrus.Fatalf("consume: %s", err)
			return
		}

		err = s.mq.Close()
		if err != nil {
			logrus.Fatalf("close: %s", err)
			return
		}

		endOfExecution := time.Now()
		resultedImage.Service = model.Conversion

		requestID, err := s.service.CreateRequest(
			r.Context(),
			req.User,
			newUploadedImage,
			resultedImage,
			model.UserImage{
				UserAccountID:   req.User.ID,
				UploadedImageID: newUploadedImage.ID,
				Status:          model.Done},
			model.Request{
				TimeStart: startOfExecution,
				EndOfTime: endOfExecution,
			})
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respondJSON(w, r, http.StatusOK, requestID)
	}
}

type findConvertedImageRequest struct {
	model.ResultedImage
	User      model.User
	requestID int
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

	convertedID, err := strconv.Atoi(convertedImageID)
	if err != nil {
		return err
	}
	req.requestID = convertedID

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
			s.errorJSON(w, r, http.StatusBadRequest, err)
			return
		}

		resultedImage, err := s.service.Image.FindTheResultingImage(r.Context(), req.requestID, model.Conversion)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, utils.ErrFindImage)
			return
		}

		err = s.service.Image.SaveImage(resultedImage.Name, "\\results\\", resultedImage.Name)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, utils.ErrSaveImage)
			return
		}

		err = s.findOriginalImage(r, req.requestID)
		if err != nil {
			s.errorJSON(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respondJSON(w, r, http.StatusOK, "successfully saved")
	}
}
