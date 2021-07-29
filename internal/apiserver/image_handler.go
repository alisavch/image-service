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
		}

		vars := mux.Vars(r)
		paramID, ok := vars["userID"]
		if !ok {
			s.error(w, r, http.StatusBadRequest, utils.ErrRequest)
		}

		intIDParam, _ := strconv.Atoi(paramID)
		if userID != intIDParam {
			s.error(w, r, http.StatusNotFound, utils.ErrPrivileges)
			return
		}
		s.respond(w, r, http.StatusOK, userID)
	}
}

func (s *Server) compressImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User
		var uploadedImage model.UploadedImage

		userID, err := s.getUserID(r)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, utils.ErrFailedConvert)
			return
		}
		user.ID = userID

		err = r.ParseMultipartForm(32 << 20)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, utils.ErrMultipartForm)
			return
		}

		newUploadedImage, err := s.uploadImage(r, uploadedImage)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}

		ratio := r.FormValue("ratio")
		if ratio == "" {
			ratio = DefaultRatio
		}
		intRatio, _ := strconv.Atoi(ratio)

		resultedImage, err := s.service.Image.CompressImage(intRatio, newUploadedImage)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		resultedImage.Service = model.Compression

		requestID, err := s.service.CreateRequest(
			user,
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

		s.respond(w, r, http.StatusOK, requestID)
	}
}

func (s *Server) findCompressedImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		userID, err := s.getUserID(r)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, utils.ErrFailedConvert)
			return
		}
		user.ID = userID

		vars := mux.Vars(r)
		compressedImageID, ok := vars["id"]
		if !ok {
			s.error(w, r, http.StatusBadRequest, utils.ErrRequest)
		}

		compressedID, _ := strconv.Atoi(compressedImageID)

		resultedImage, err := s.service.Image.FindTheResultingImage(compressedID, model.Compression)
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

		s.respond(w, r, http.StatusOK, "successfully saved")
	}
}

func (s *Server) convertImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User
		var uploadedImage model.UploadedImage

		userID, err := s.getUserID(r)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, utils.ErrFailedConvert)
			return
		}
		user.ID = userID

		err = r.ParseMultipartForm(32 << 20)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, utils.ErrMultipartForm)
			return
		}
		newUploadedImage, err := s.uploadImage(r, uploadedImage)
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
			user,
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

		s.respond(w, r, http.StatusOK, requestID)
	}
}

func (s *Server) findConvertedImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		userID, err := s.getUserID(r)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, utils.ErrFailedConvert)
			return
		}
		user.ID = userID

		vars := mux.Vars(r)
		convertedImageID, ok := vars["id"]
		if !ok {
			s.error(w, r, http.StatusBadRequest, utils.ErrRequest)
		}

		convertedID, _ := strconv.Atoi(convertedImageID)

		resultedImage, err := s.service.Image.FindTheResultingImage(convertedID, model.Conversion)
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

		s.respond(w, r, http.StatusOK, "successfully saved")
	}
}
