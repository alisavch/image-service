package apiserver

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/alisavch/image-service/internal/model"
	"github.com/gorilla/mux"
)

const (
	// DefaultRatio is default value for compress jpeg.
	DefaultRatio = "95"
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
		return fmt.Errorf("failed convert to int userID")
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
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		history, err := s.service.Image.FindUserHistoryByID(r.Context(), req.User.ID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respondJSON(w, r, http.StatusOK, history)
	}
}

type uploadImageRequest struct {
	model.UploadedImage
	User  model.User
	Ratio int
}

// Build builds a request to compress image.
func (req *uploadImageRequest) Build(r *http.Request) error {
	id, ok := r.Context().Value(userCtx).(int)
	if !ok {
		return fmt.Errorf("failed convert to int userID")
	}
	req.User.ID = id

	ratio := r.FormValue("ratio")
	if ratio == "" {
		ratio = DefaultRatio
	}
	req.Ratio, _ = strconv.Atoi(ratio)

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
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		startOfExecution := time.Now()

		//rabbit := new(broker.RabbitMQ)
		//if err = rabbit.Connect(); err != nil {
		//	s.error(w, r, http.StatusInternalServerError, err)
		//}
		//defer rabbit.Close()
		//
		//if err = rabbit.DeclareQueue(model.Queued); err != nil {
		//	s.error(w, r, http.StatusInternalServerError, err)
		//}

		newUploadedImage, err := s.uploadImage(r, req.UploadedImage)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		resultedImage, err := s.service.Image.CompressImage(req.Ratio, newUploadedImage)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
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
				Status:          model.Queued},
			model.Request{
				TimeStart: startOfExecution,
				EndOfTime: endOfExecution,
			})
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
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
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		resultedImage, err := s.service.Image.FindTheResultingImage(r.Context(), req.CompressedID, model.Compression)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, fmt.Errorf("cannot find image\""))
			return
		}

		err = s.service.Image.SaveImage(resultedImage.Name, "\\results\\", resultedImage.Name)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, fmt.Errorf("cannot save image"))
			return
		}

		err = s.findOriginalImage(r, req.CompressedID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
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
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		startOfExecution := time.Now()

		newUploadedImage, err := s.uploadImage(r, req.UploadedImage)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		resultedImage, err := s.service.Image.ConvertToType(newUploadedImage)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, fmt.Errorf("cannot convert image"))
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
				Status:          model.Queued},
			model.Request{
				TimeStart: startOfExecution,
				EndOfTime: endOfExecution,
			})
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
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
		return fmt.Errorf("failed convert to int userID")
	}
	req.User.ID = id

	vars := mux.Vars(r)
	convertedImageID, ok := vars["convertedID"]
	if !ok {
		return fmt.Errorf("incorrect request")
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
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		resultedImage, err := s.service.Image.FindTheResultingImage(r.Context(), req.requestID, model.Conversion)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, fmt.Errorf("cannot find image\""))
			return
		}

		err = s.service.Image.SaveImage(resultedImage.Name, "\\results\\", resultedImage.Name)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, fmt.Errorf("cannot save image"))
			return
		}

		err = s.findOriginalImage(r, req.requestID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respondJSON(w, r, http.StatusOK, "successfully saved")
	}
}
