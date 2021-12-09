package broker

import (
	"context"
	"fmt"
	"image"
	"os"

	"github.com/alisavch/image-service/internal/models"
	"github.com/alisavch/image-service/internal/service"
	"github.com/alisavch/image-service/internal/utils"
)

const (
	aws   = "AWS"
	local = "local"
)

// Process service processes.
func (r *RabbitMQ) Process(message models.QueuedMessage) error {
	conf := utils.NewConfig()
	ctx := context.Background()

	switch message.Service {
	case models.Compression:
		compressedImage, err := r.Compress(message, conf.Storage)
		if err != nil {
			r.logger.Printf("%s:%s", "Failed to compress image", err)
			return err
		}
		message.Image.ResultedName = compressedImage.ResultedName
		message.Image.ResultedLocation = compressedImage.ResultedLocation

	case models.Conversion:
		convertedImage, err := r.Convert(message, conf.Storage)
		if err != nil {
			r.logger.Printf("%s:%s", "Failed to convert image", err)
			return err
		}
		message.Image.ResultedName = convertedImage.ResultedName
		message.Image.ResultedLocation = convertedImage.ResultedLocation
	}

	err := r.UploadResultedImage(ctx, message.Image)
	if err != nil {
		return err
	}
	r.logger.Printf("%s:%s", "Resulted image uploaded", message.Image.ResultedName)

	err = r.CompleteRequest(ctx, message.RequestID, models.Done)
	if err != nil {
		return err
	}
	r.logger.Printf("%s:%s", "Request completed", message.RequestID)

	return nil
}

// Compress is the compression service.
func (r *RabbitMQ) Compress(message models.QueuedMessage, storage string) (models.Image, error) {
	r.logger.Printf("%s:%s", "Process started", message.Service)
	resultedName := newImgName("cmp-" + message.UploadedName)
	r.logger.Printf("%s:%s", "Image renamed", resultedName)

	img, format, file, err := r.prepareImage(message.Image, message.Image.UploadedName, resultedName)
	if err != nil {
		return models.Image{}, err
	}

	compressedImage, err := r.CompressImage(message.Width, format, resultedName, img, file, storage)
	if err != nil {
		return models.Image{}, err
	}
	r.logger.Printf("%s:%s", "Process finished", message.Service)

	return compressedImage, nil
}

// Convert is the conversion service.
func (r *RabbitMQ) Convert(message models.QueuedMessage, storage string) (models.Image, error) {
	r.logger.Printf("%s:%s", "Process started", message.Service)
	convertedName, err := r.ChangeFormat(message.UploadedName)
	if err != nil {
		return models.Image{}, err
	}
	r.logger.Printf("%s:%s", "Image format changed", convertedName)

	resultedName := newImgName("cnv-" + convertedName)
	r.logger.Printf("%s:%s", "Image renamed", resultedName)

	img, format, file, err := r.prepareImage(message.Image, message.Image.UploadedName, resultedName)
	if err != nil {
		return models.Image{}, err
	}

	convertedImage, err := r.ConvertToType(format, resultedName, img, file, storage)
	if err != nil {
		return models.Image{}, err
	}
	r.logger.Printf("%s:%s", "Process finished", message.Service)

	return convertedImage, nil
}

func (r *RabbitMQ) prepareImage(uploadedImage models.Image, originalImageName, resultedImageName string) (image.Image, string, *os.File, error) {
	conf := utils.NewConfig()

	switch conf.Storage {
	case aws:
		img, format, resultedFile, err := r.downloadOriginalImageFormAWS(originalImageName, resultedImageName)
		if err != nil {
			return nil, "", nil, err
		}
		return img, format, resultedFile, nil

	case local:
		img, format, resultedFile, err := getOriginalImageLocally(uploadedImage, resultedImageName)
		if err != nil {
			return nil, "", nil, err
		}
		return img, format, resultedFile, nil
	}

	return nil, "", nil, nil
}

func (r *RabbitMQ) downloadOriginalImageFormAWS(originalImageName, resultedImageName string) (image.Image, string, *os.File, error) {
	file, err := r.DownloadFromS3Bucket(originalImageName)
	if err != nil {
		return nil, "", nil, err
	}

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", nil, utils.ErrDecode
	}

	resultedFile, err := os.Create(resultedImageName)
	if err != nil {
		return nil, "", nil, utils.ErrCreateFile
	}

	return img, format, resultedFile, nil
}

func getOriginalImageLocally(uploadedImage models.Image, resultedImageName string) (image.Image, string, *os.File, error) {
	file, err := os.Open(uploadedImage.UploadedLocation + uploadedImage.UploadedName)
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

	resultedFile, err := os.Create(fmt.Sprintf("./results/%s", resultedImageName))
	if err != nil {
		return nil, "", nil, utils.ErrCreateFile
	}

	return img, format, resultedFile, nil
}

func newImgName(str string) string {
	return str
}
