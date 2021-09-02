package service

import (
	"context"
	"fmt"
	"image"
	"io"
	"os"

	"github.com/alisavch/image-service/internal/utils"

	"github.com/alisavch/image-service/internal/models"
	"github.com/alisavch/image-service/internal/repository"
)

var convertedType = map[string]string{
	"jpeg": "png",
	"jpg":  "png",
	"png":  "jpeg",
}

// ImageService provides access to repository.
type ImageService struct {
	repo repository.Image
}

// NewImageService is ImageService constructor.
func NewImageService(repo repository.Image) *ImageService {
	return &ImageService{repo: repo}
}

// FindUserHistoryByID allows to get the history of interaction with the user's service.
func (s *ImageService) FindUserHistoryByID(ctx context.Context, id int) ([]models.History, error) {
	return s.repo.FindUserHistoryByID(ctx, id)
}

// UploadImage uploads image.
func (s *ImageService) UploadImage(ctx context.Context, image models.UploadedImage) (int, error) {
	return s.repo.UploadImage(ctx, image)
}

// CompressImage compress image.
func (s *ImageService) CompressImage(width int, uploadedImage models.UploadedImage) (models.ResultedImage, error) {
	var result models.ResultedImage
	currentDir, err := os.Getwd()
	if err != nil {
		return models.ResultedImage{}, err
	}

	imgFile, err := os.Open(uploadedImage.Location + uploadedImage.Name)
	if err != nil {
		return models.ResultedImage{}, err
	}
	defer imgFile.Close()

	imgSrc, format, err := image.Decode(imgFile)
	if err != nil {
		return models.ResultedImage{}, err
	}

	err = EnsureBaseDir("./results/")
	if err != nil {
		return models.ResultedImage{}, err
	}

	newImgFile, err := os.Create(fmt.Sprintf("./results/%s", uploadedImage.Name))
	if err != nil {
		return models.ResultedImage{}, err
	}
	defer newImgFile.Close()

	switch format {
	case "jpeg":
		if err := CompressJPEG(imgSrc, width, newImgFile); err != nil {
			return models.ResultedImage{}, err
		}
	case "png":
		if err := CompressPNG(imgSrc, width, newImgFile); err != nil {
			return models.ResultedImage{}, err
		}
	}
	result.Name = uploadedImage.Name
	result.Location = currentDir + "/results/"

	return result, nil
}

// ConvertToType converts from png to jpeg and vice versa.
func (s *ImageService) ConvertToType(uploadedImage models.UploadedImage) (models.ResultedImage, error) {
	var result models.ResultedImage
	currentDir, err := os.Getwd()
	if err != nil {
		return models.ResultedImage{}, err
	}

	file, err := os.Open(uploadedImage.Location + uploadedImage.Name)
	if err != nil {
		return models.ResultedImage{}, err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return models.ResultedImage{}, err
	}

	convertedName, err := ChangeFormat(uploadedImage.Name)
	if err != nil {
		return models.ResultedImage{}, err
	}

	err = EnsureBaseDir("./results/")
	if err != nil {
		return models.ResultedImage{}, err
	}

	newImg, err := os.Create(fmt.Sprintf("./results/%s", convertedName))
	defer file.Close()
	if err != nil {
		return models.ResultedImage{}, utils.ErrCreateFile
	}
	defer newImg.Close()

	err = Encode(newImg, img, format)
	if err != nil {
		return models.ResultedImage{}, err
	}

	result.Name = convertedName
	result.Location = currentDir + "/results/"
	return result, nil
}

// CreateRequest creates request.
func (s *ImageService) CreateRequest(ctx context.Context, user models.User, uplImg models.UploadedImage, resImg models.ResultedImage, uI models.UserImage, r models.Request) (int, error) {
	return s.repo.CreateRequest(ctx, user, uplImg, resImg, uI, r)
}

// FindTheResultingImage finds the resulting image by id.
func (s *ImageService) FindTheResultingImage(ctx context.Context, id int, service models.Service) (models.ResultedImage, error) {
	return s.repo.FindTheResultingImage(ctx, id, service)
}

// FindOriginalImage finds original image by id.
func (s *ImageService) FindOriginalImage(ctx context.Context, id int) (models.UploadedImage, error) {
	return s.repo.FindOriginalImage(ctx, id)
}

// SaveImage saves image to users machine.
func (s *ImageService) SaveImage(filename, folder, resultedFilename string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	file, err := os.Open(currentDir + folder + filename)
	if err != nil {
		return err
	}

	out, err := SaveToDownloads(resultedFilename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return err
	}
	return nil
}
