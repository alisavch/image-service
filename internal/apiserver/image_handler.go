package apiserver

import (
	"net/http"
	"strconv"

	"github.com/alisavch/image-service/internal/model"
	"github.com/alisavch/image-service/internal/utils"
	"github.com/gorilla/mux"
)

const (
	// DefaultRatio is default value for compress jpeg.
	DefaultRatio = "95"
	// DefaultOriginal is default value for downloads original image.
	DefaultOriginal = "false"
)

func (s *Server) findUserHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := s.getUserID(r)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, utils.ErrRequest)
			return
		}

		result, err := s.service.Image.FindUserHistoryByID(r.Context(), userID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respondJSON(w, r, http.StatusOK, result)
	}
}

type uploadImageRequest struct {
	model.UploadedImage
	model.User
}

// Build builds a request to compress image.
func (req uploadImageRequest) Build(r *http.Request) (int, error) {
	ratio := r.FormValue("ratio")
	if ratio == "" {
		ratio = DefaultRatio
	}
	intRatio, err := strconv.Atoi(ratio)
	if err != nil {
		return 0, err
	}
	return intRatio, nil
}

func (s *Server) compressImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req uploadImageRequest

		userID, err := s.getUserID(r)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, utils.ErrFailedConvert)
			return
		}
		req.User.ID = userID

		newUploadedImage, err := s.uploadImage(r, req.UploadedImage)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		intRatio, err := ParseImageRequest(r, req)

		resultedImage, err := s.service.Image.CompressImage(intRatio, newUploadedImage)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		resultedImage.Service = model.Compression

		requestID, err := s.service.CreateRequest(
			r.Context(),
			req.User,
			newUploadedImage,
			resultedImage,
			model.UserImage{
				UserAccountID:   userID,
				UploadedImageID: newUploadedImage.ID,
				Status:          model.Queued},
			model.Request{})
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respondJSON(w, r, http.StatusOK, requestID)
	}
}

type findCompressedImageRequest struct {
	model.ResultedImage
	model.User
}

// Build builds a request to find compressed image.
func (req findCompressedImageRequest) Build(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	compressedImageID, ok := vars["id"]
	if !ok {
		return 0, utils.ErrRequest
	}

	compressedID, err := strconv.Atoi(compressedImageID)
	if err != nil {
		return 0, err
	}
	return compressedID, nil
}

func (s *Server) findCompressedImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req findCompressedImageRequest

		userID, err := s.getUserID(r)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, utils.ErrFailedConvert)
			return
		}
		req.User.ID = userID

		compressedID, err := ParseImageRequest(r, req)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		resultedImage, err := s.service.Image.FindTheResultingImage(r.Context(), compressedID, model.Compression)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, utils.ErrFindImage)
			return
		}

		err = s.service.Image.SaveImage(resultedImage.Name, "\\results\\", resultedImage.Name)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, utils.ErrSaveImage)
			return
		}

		err = s.findOriginalImage(r, compressedID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respondJSON(w, r, http.StatusOK, "successfully saved")
	}
}

type convertImageRequest struct {
	model.UploadedImage
	model.User
}

func (s *Server) convertImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req convertImageRequest

		userID, err := s.getUserID(r)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, utils.ErrFailedConvert)
			return
		}
		req.User.ID = userID

		err = r.ParseMultipartForm(32 << 20)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, utils.ErrMultipartForm)
			return
		}
		newUploadedImage, err := s.uploadImage(r, req.UploadedImage)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		resultedImage, err := s.service.Image.ConvertToType(newUploadedImage)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, utils.ErrConvert)
			return
		}

		resultedImage.Service = model.Conversion

		requestID, err := s.service.CreateRequest(
			r.Context(),
			req.User,
			newUploadedImage,
			resultedImage,
			model.UserImage{
				UserAccountID:   userID,
				UploadedImageID: newUploadedImage.ID,
				Status:          model.Queued},
			model.Request{})
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, utils.ErrCreateRequest)
			return
		}

		s.respondJSON(w, r, http.StatusOK, requestID)
	}
}

type findConvertedImageRequest struct {
	model.ResultedImage
	model.User
}

// Build builds a request to find converted image.
func (req findConvertedImageRequest) Build(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	convertedImageID, ok := vars["id"]
	if !ok {
		return 0, utils.ErrRequest
	}

	convertedID, err := strconv.Atoi(convertedImageID)
	if err != nil {
		return 0, err
	}
	return convertedID, nil
}

func (s *Server) findConvertedImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req findConvertedImageRequest

		userID, err := s.getUserID(r)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, utils.ErrFailedConvert)
			return
		}
		req.User.ID = userID

		convertedID, err := ParseImageRequest(r, req)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		resultedImage, err := s.service.Image.FindTheResultingImage(r.Context(), convertedID, model.Conversion)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, utils.ErrFindImage)
			return
		}

		err = s.service.Image.SaveImage(resultedImage.Name, "\\results\\", resultedImage.Name)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, utils.ErrSaveImage)
			return
		}

		err = s.findOriginalImage(r, convertedID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respondJSON(w, r, http.StatusOK, "successfully saved")
	}
}
