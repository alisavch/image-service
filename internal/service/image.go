package service

import (
	"context"
	"fmt"
	"image"
	"os"

	"github.com/google/uuid"

	"github.com/alisavch/image-service/internal/utils"

	"github.com/alisavch/image-service/internal/models"
	"github.com/alisavch/image-service/internal/repository"
)

var (
	convertedType = map[string]string{
		"jpeg": "png",
		"jpg":  "png",
		"png":  "jpeg",
	}
	remoteStorage = NewAWS()
	logger        = NewLogger()
)

// ImageService provides access to repository.
type ImageService struct {
	repo repository.Image
}

// NewImageService is the ImageService constructor.
func NewImageService(repo repository.Image) *ImageService {
	return &ImageService{repo: repo}
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
func (s *ImageService) CompressImage(width int, format, resultedName string, img image.Image, newImg *os.File, isRemoteStorage bool) (models.ResultedImage, error) {
	var result models.ResultedImage

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

	if isRemoteStorage {
		fileReader, err := os.Open(resultedName)
		if err != nil {
			return models.ResultedImage{}, fmt.Errorf("%s:%s", "unable to open file", err)
		}
		defer func(fileReader *os.File) {
			err := fileReader.Close()
			if err != nil {
				logger.Fatalf("%s:%s", "failed fileReader.Close", err)
			}
		}(fileReader)

		imageLocation, err := remoteStorage.UploadToS3Bucket(fileReader, resultedName)
		if err != nil {
			return models.ResultedImage{}, fmt.Errorf("%s:%s", "failed to upload to aws", err)
		}

		result.Name = resultedName
		result.Location = imageLocation

		return result, nil
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return models.ResultedImage{}, utils.ErrGetDir
	}

	err = newImg.Close()
	if err != nil {
		return models.ResultedImage{}, err
	}

	result.Name = resultedName
	result.Location = currentDir + "/results/"

	return result, nil
}

// ConvertToType converts from png to jpeg and vice versa.
func (s *ImageService) ConvertToType(format, resultedName string, img image.Image, newImg *os.File, isRemoteStorage bool) (models.ResultedImage, error) {
	var result models.ResultedImage
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

	if isRemoteStorage {
		fileReader, err := os.Open(resultedName)
		if err != nil {
			return models.ResultedImage{}, fmt.Errorf("%s:%s", "unable to open file", err)
		}
		defer func(fileReader *os.File) {
			err := fileReader.Close()
			if err != nil {
				logger.Fatalf("%s:%s", "failed fileReader.Close", err)
			}
		}(fileReader)

		imageLocation, err := remoteStorage.UploadToS3Bucket(fileReader, resultedName)
		if err != nil {
			return models.ResultedImage{}, fmt.Errorf("%s:%s", "failed to upload to aws", err)
		}

		err = newImg.Close()
		if err != nil {
			return models.ResultedImage{}, err
		}

		result.Name = resultedName
		result.Location = imageLocation

		return result, nil
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return models.ResultedImage{}, utils.ErrGetDir
	}

	err = newImg.Close()
	if err != nil {
		return models.ResultedImage{}, err
	}

	result.Name = resultedName
	result.Location = currentDir + "/results/"

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
func (s *ImageService) SaveImage(filename, location string, isRemoteStorage bool) (*models.Image, error) {
	var file *os.File
	var err error
	img := models.Image{Filename: filename}

	switch isRemoteStorage {
	case true:
		f, err := remoteStorage.DownloadFromS3Bucket(filename)
		if err != nil {
			return nil, fmt.Errorf("%s:%s", "cannot download from s3 bucket", err)
		}
		file = f

	default:
		f, err := os.Open(location + filename)
		if err != nil {
			return &models.Image{}, utils.ErrOpen
		}
		file = f
	}

	img.File = file

	img.ContentType, err = GetFileContentType(file)
	if err != nil {
		return &models.Image{}, fmt.Errorf("%s:%s", "cannot get file content type", err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return &models.Image{}, utils.ErrFileStat
	}

	img.Filesize = fileInfo.Size()

	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	return &img, nil
}

// UpdateStatus updates the status of image processing.
func (s *ImageService) UpdateStatus(ctx context.Context, id uuid.UUID, status models.Status) error {
	return s.repo.UpdateStatus(ctx, id, status)
}
