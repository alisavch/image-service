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

	"github.com/google/uuid"

	"github.com/alisavch/image-service/internal/service"

	"github.com/alisavch/image-service/internal/utils"

	_ "image/jpeg" // It allows using jpeg
	_ "image/png"  // It allows using png

	"github.com/alisavch/image-service/internal/models"
)

type key string

var remoteStorage = NewAWS()

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
	if req.token == "" {
		return utils.ErrEmptyToken
	}
	return nil
}

func (s *Server) authorize(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req authorization

		err := ParseRequest(r, &req)
		if err != nil {
			s.errorJSON(w, http.StatusBadRequest, err)
		}

		userID, err := s.service.Authorization.ParseToken(req.token)
		if err != nil {
			s.errorJSON(w, http.StatusUnauthorized, err)
			return
		}
		ctx := context.WithValue(r.Context(), userCtx, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

type uploaded struct {
	file    multipart.File
	handler *multipart.FileHeader
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
	// TODO: delete
	// req.handler.Filename = strings.Replace(uuid.New().String(), "-", "", -1) + req.handler.Filename
	req.handler.Filename = strings.ReplaceAll(uuid.New().String(), "-", "") + req.handler.Filename

	if IsRemoteStorage {
		imageLocation, err := remoteStorage.UploadToS3Bucket(req.file, req.handler.Filename)
		if err != nil {
			return models.UploadedImage{}, err
		}

		err = req.file.Close()
		if err != nil {
			return models.UploadedImage{}, err
		}
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
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			logger.Fatalf("%s:%s", "failed fileReader.Close", err)
		}
	}(out)

	_, err = io.Copy(out, req.file)
	if err != nil {
		return models.UploadedImage{}, utils.ErrCopyFile
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			logger.Fatalf("%s:%s", "failed fileReader.Close", err)
		}
	}(req.file)

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

func prepareImage(uploadedImage models.UploadedImage, originalImageName, resultedImageName string) (image.Image, string, *os.File, error) {
	if IsRemoteStorage {
		file, err := remoteStorage.DownloadFromS3Bucket(originalImageName)
		if err != nil {
			return nil, "", nil, err
		}

		img, format, err := image.Decode(file)
		if err != nil {
			logger.Fatalf("%s:%s", "Cannot decode file", err)
			return nil, "", nil, utils.ErrDecode
		}

		resultedFile, err := os.Create(resultedImageName)
		if err != nil {
			logger.Fatalf("%s:%s", "Cannot create resulted file", err)
		}

		return img, format, resultedFile, nil
	}

	file, err := os.Open(uploadedImage.Location + uploadedImage.Name)
	if err != nil {
		return nil, "", nil, utils.ErrOpen
	}

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", nil, utils.ErrDecode
	}

	err = file.Close()
	if err != nil {
		return nil, "", nil, err
	}

	err = service.EnsureBaseDir("./results/")
	if err != nil {
		return nil, "", nil, utils.ErrEnsureDir
	}

	newImg, err := os.Create(fmt.Sprintf("./results/%s", resultedImageName))
	if err != nil {
		return nil, "", nil, utils.ErrCreateFile
	}

	return img, format, newImg, nil
}

func newImgName(str string) string {
	return str
}
