package service

import (
	"context"
	"fmt"
	"image"
	"os"
	"strings"

	"github.com/alisavch/image-service/internal/log"
	"github.com/alisavch/image-service/internal/models"
	"github.com/alisavch/image-service/internal/utils"

	"github.com/google/uuid"
)

const (
	aws   = "AWS"
	local = "local"
)

var (
	convertedType = map[string]string{
		"jpeg": "png",
		"jpg":  "png",
		"png":  "jpeg",
	}
)

// ImageService provides access to repository.
type ImageService struct {
	repo   ImageRepo
	bucket S3Bucket
	logger FormattingOutput
}

// NewImageService configures ImageService.
func NewImageService(repo ImageRepo, bucket S3Bucket) *ImageService {
	return &ImageService{
		repo:   repo,
		bucket: bucket,
		logger: log.NewLogger(),
	}
}

// FindUserHistoryByID allows to get the history of interaction with the user's service.
func (s *ImageService) FindUserHistoryByID(ctx context.Context, id uuid.UUID) ([]models.History, error) {
	return s.repo.FindUserHistoryByID(ctx, id)
}

// UploadImage uploads image.
func (s *ImageService) UploadImage(ctx context.Context, img models.UploadedImage) (uuid.UUID, error) {
	return s.repo.UploadImage(ctx, img)
}

// CompressImage compress image.
func (s *ImageService) CompressImage(width int, format, resultedName string, img image.Image, newImg *os.File, storage string) (models.ResultedImage, error) {
	switch format {
	case "jpeg":
		if err := CompressJPEG(img, width, newImg); err != nil {
			return models.ResultedImage{}, fmt.Errorf("%s:%s", utils.ErrCompress, err)
		}
	case "png":
		if err := CompressPNG(img, width, newImg); err != nil {
			return models.ResultedImage{}, fmt.Errorf("%s:%s", utils.ErrCompress, err)
		}
	}

	result, err := s.FillInTheResultingImage(storage, resultedName, newImg)
	if err != nil {
		return models.ResultedImage{}, err
	}

	return result, nil
}

// ConvertToType converts from png to jpeg and vice versa.
func (s *ImageService) ConvertToType(format, resultedName string, img image.Image, newImg *os.File, storage string) (models.ResultedImage, error) {
	switch format {
	case "jpeg":
		if err := ConvertToPNG(newImg, img); err != nil {
			return models.ResultedImage{}, fmt.Errorf("%s:%s", utils.ErrCompress, err)
		}
	case "png":
		if err := ConvertToJPEG(newImg, img); err != nil {
			return models.ResultedImage{}, fmt.Errorf("%s:%s", utils.ErrCompress, err)
		}
	}

	result, err := s.FillInTheResultingImage(storage, resultedName, newImg)
	if err != nil {
		return models.ResultedImage{}, err
	}

	return result, nil
}

// CreateRequest creates request.
func (s *ImageService) CreateRequest(ctx context.Context, user models.User, uplImg models.UploadedImage, resImg models.ResultedImage, uI models.UserImage, r models.Request) (uuid.UUID, error) {
	return s.repo.CreateRequest(ctx, user, uplImg, resImg, uI, r)
}

// FindTheResultingImage finds the resulting image by id.
func (s *ImageService) FindTheResultingImage(ctx context.Context, id uuid.UUID, service models.Service) (models.ResultedImage, error) {
	return s.repo.FindTheResultingImage(ctx, id, service)
}

// FindOriginalImage finds original image by id.
func (s *ImageService) FindOriginalImage(ctx context.Context, id uuid.UUID) (models.UploadedImage, error) {
	return s.repo.FindOriginalImage(ctx, id)
}

// SaveImage saves image to users machine.
func (s *ImageService) SaveImage(filename, location, storage string) (*models.Image, error) {
	var file *os.File
	img := models.Image{Filename: filename}

	switch storage {
	case aws:
		f, err := s.bucket.DownloadFromS3Bucket(filename)
		if err != nil {
			return nil, fmt.Errorf("%s:%s", utils.ErrRemoteDownload, err)
		}
		file = f

	case local:
		f, err := os.Open(location + filename)
		if err != nil {
			return nil, utils.ErrOpen
		}
		file = f
	}

	img, err := FillInTheImage(img, file)
	if err != nil {
		return nil, err
	}

	return &img, nil
}

// UpdateStatus updates the status of image processing.
func (s *ImageService) UpdateStatus(ctx context.Context, id uuid.UUID, status models.Status) error {
	return s.repo.UpdateStatus(ctx, id, status)
}

// ChangeFormat changes image format.
func (s *ImageService) ChangeFormat(filename string) (string, error) {
	imgNames := strings.Split(strings.ToLower(filename), ".")
	extension := imgNames[len(imgNames)-1]

	if format, ok := convertedType[extension]; ok {
		imgNames[len(imgNames)-1] = format
		convertedName := strings.Join(imgNames, ".")
		return convertedName, nil
	}

	return "", utils.ErrUnsupportedFormat
}

// FillInTheResultingImageForAWS fills the resultedImage with information for the aws storage.
func (s *ImageService) FillInTheResultingImageForAWS(resultedName string) (models.ResultedImage, error) {
	fileReader, err := os.Open(resultedName)
	if err != nil {
		return models.ResultedImage{}, fmt.Errorf("%s:%s", utils.ErrOpenFile, err)
	}
	defer func(fileReader *os.File) {
		err := fileReader.Close()
		if err != nil {
			s.logger.Fatalf("%s:%s", "failed fileReader.Close", err)
		}
	}(fileReader)

	imageLocation, err := s.bucket.UploadToS3Bucket(fileReader, resultedName)
	if err != nil {
		return models.ResultedImage{}, fmt.Errorf("%s:%s", utils.ErrRemoteUpload, err)
	}

	result := FillInTheReceivedNameAndLocation(resultedName, imageLocation)

	return result, nil
}

// FillInTheResultingImage fills the resultedImage with information.
func (s *ImageService) FillInTheResultingImage(storage, resultedName string, newImg *os.File) (models.ResultedImage, error) {
	var result models.ResultedImage

	switch storage {
	case aws:
		resultedImage, err := s.FillInTheResultingImageForAWS(resultedName)
		if err != nil {
			return models.ResultedImage{}, err
		}
		result = resultedImage

	case local:
		resultedImage, err := FillInTheResultingImageLocally(resultedName, newImg)
		if err != nil {
			return models.ResultedImage{}, err
		}
		result = resultedImage
	}

	return result, nil
}
