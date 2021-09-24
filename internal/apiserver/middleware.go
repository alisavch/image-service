package apiserver

import (
	"context"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/alisavch/image-service/internal/bucket"

	"github.com/alisavch/image-service/internal/service"

	"github.com/alisavch/image-service/internal/utils"

	"github.com/alisavch/image-service/internal/models"
	"github.com/google/uuid"

	_ "image/jpeg"
	_ "image/png"
)

type key string

var remoteStorage = bucket.NewAWS()

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

type authorization struct {
	header      string
	headerParts []string
	token       string
}

// Build builds a request to authorize.
func (req *authorization) Build(r *http.Request) error {
	req.header = r.Header.Get(string(authorizationHeader))
	if req.header == "" {
		return utils.ErrEmptyHeader
	}
	req.headerParts = strings.Split(req.header, " ")
	req.token = req.headerParts[1]
	return nil
}

// Validate validates request to authorize.
func (req authorization) Validate() error {
	if len(req.headerParts) != 2 || req.headerParts[0] != "Bearer" {
		return utils.ErrInvalidAuthHeader
	}
	if len(req.token) == 0 {
		return utils.ErrEmptyToken
	}
	return nil
}

func (s *Server) authorize(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req authorization

		err := ParseRequest(r, &req)
		if err != nil {
			s.errorJSON(w, r, http.StatusBadRequest, err)
		}

		userID, err := s.service.Authorization.ParseToken(req.token)
		if err != nil {
			s.errorJSON(w, r, http.StatusUnauthorized, err)
			return
		}
		ctx := context.WithValue(r.Context(), userCtx, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

type uploaded struct {
	file            multipart.File
	handler         *multipart.FileHeader
	isRemoteStorage bool
}

// Build builds a request to load an image.
func (req *uploaded) Build(r *http.Request) error {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return err
	}

	req.file, req.handler, err = r.FormFile("uploadFile")
	if err != nil {
		return err
	}

	req.isRemoteStorage = IsRemoteStorage

	return nil
}

// Validate builds a request to validate the upload of an image.
func (req uploaded) Validate() error {
	if !(req.handler.Header["Content-Type"][0] == "image/jpeg" || req.handler.Header["Content-Type"][0] == "image/png") {
		return utils.ErrAllowedFormat
	}
	return nil
}

func (s *Server) uploadImage(r *http.Request, uploadedImage models.UploadedImage) (models.UploadedImage, error) {
	var req uploaded
	err := ParseRequest(r, &req)
	if err != nil {
		return models.UploadedImage{}, err
	}
	req.handler.Filename = strings.Replace(uuid.New().String(), "-", "", -1) + req.handler.Filename

	if req.isRemoteStorage {
		imageLocation, err := remoteStorage.UploadToS3Bucket(req.file, req.handler.Filename)
		if err != nil {
			return models.UploadedImage{}, err
		}

		req.file.Close()
		uploadedImage.Name = req.handler.Filename
		uploadedImage.Location = imageLocation

		uploadedID, err := s.service.Image.UploadImage(r.Context(), uploadedImage)
		if err != nil {
			return models.UploadedImage{}, fmt.Errorf("%s:%s", utils.ErrUpload, err)
		}
		uploadedImage.ID = uploadedID

		return uploadedImage, nil
	}

	err = service.EnsureBaseDir("./uploads/")
	if err != nil {
		return models.UploadedImage{}, err
	}

	out, err := os.Create(fmt.Sprintf("./uploads/%s", req.handler.Filename))
	if err != nil {
		return models.UploadedImage{}, utils.ErrCreateFile
	}
	defer out.Close()

	_, err = io.Copy(out, req.file)
	defer req.file.Close()
	if err != nil {
		return models.UploadedImage{}, utils.ErrCopyFile
	}

	uploadedImage.Name = req.handler.Filename
	currentDir, _ := os.Getwd()
	uploadedImage.Location = currentDir + "/uploads/"
	uploadedID, err := s.service.Image.UploadImage(r.Context(), uploadedImage)
	if err != nil {
		return models.UploadedImage{}, utils.ErrUpload
	}
	uploadedImage.ID = uploadedID

	return uploadedImage, nil
}

type prepare struct {
	isRemoteStorage bool
}

// Build builds a request to load an image.
func (req *prepare) Build(r *http.Request) error {
	req.isRemoteStorage = IsRemoteStorage
	return nil
}

// Validate builds a request to validate the preparation of an image.
func (req prepare) Validate() error {
	return nil
}

func (s *Server) prepareImage(r *http.Request, uploadedImage models.UploadedImage, originalImageName string, resultedImageName string) (image.Image, string, string, *os.File, bool, error) {
	var req prepare
	err := ParseRequest(r, &req)
	if err != nil {
		return nil, "", "", nil, false, err
	}

	if req.isRemoteStorage {
		file, err := remoteStorage.DownloadFromS3Bucket(originalImageName)
		if err != nil {
			return nil, "", "", nil, req.isRemoteStorage, err
		}

		img, format, err := image.Decode(file)
		if err != nil {
			logger.Fatalf("%s:%s", "Cannot decode file", err)
			return nil, "", "", nil, req.isRemoteStorage, utils.ErrDecode
		}

		resultedFile, err := os.Create(resultedImageName)
		if err != nil {
			logger.Fatalf("%s:%s", "Cannot create resulted file", err)
		}

		return img, format, resultedImageName, resultedFile, req.isRemoteStorage, nil
	}

	file, err := os.Open(uploadedImage.Location + uploadedImage.Name)
	if err != nil {
		return nil, "", "", nil, req.isRemoteStorage, utils.ErrOpen
	}

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", "", nil, req.isRemoteStorage, utils.ErrDecode
	}

	file.Close()

	err = service.EnsureBaseDir("./results/")
	if err != nil {
		return nil, "", "", nil, req.isRemoteStorage, utils.ErrEnsureDir
	}

	newImg, err := os.Create(fmt.Sprintf("./results/%s", resultedImageName))
	if err != nil {
		return nil, "", "", nil, req.isRemoteStorage, utils.ErrCreateFile
	}

	return img, format, resultedImageName, newImg, req.isRemoteStorage, nil
}
