package service

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"os"

	"github.com/alisavch/image-service/internal/model"
	"github.com/alisavch/image-service/internal/repository"
	"github.com/alisavch/image-service/internal/utils"
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

// UploadImage uploads image.
func (s *ImageService) UploadImage(ctx context.Context, image model.UploadedImage) (int, error) {
	return s.repo.UploadImage(ctx, image)
}

// CompressImage compress image.
func (s *ImageService) CompressImage(quality int, uploadedImage model.UploadedImage) (model.ResultedImage, error) {
	var result model.ResultedImage
	currentDir, _ := os.Getwd()
	// TODO: image location instead of --currentDir+ "\\uploads\\"--
	imgFile, err := os.Open(currentDir + "\\uploads\\" + uploadedImage.Name)
	if err != nil {
		return model.ResultedImage{}, err
	}
	defer imgFile.Close()
	imgSrc, format, err := image.Decode(imgFile)
	if err != nil {
		return model.ResultedImage{}, err
	}
	newImgFile, err := os.Create(fmt.Sprintf("./results/%s", uploadedImage.Name))
	if err != nil {
		return model.ResultedImage{}, err
	}
	defer newImgFile.Close()
	switch format {
	case "jpeg":
		if quality < 0 || quality > 100 {
			return model.ResultedImage{}, utils.ErrJPEG
		}
		err = jpeg.Encode(newImgFile, imgSrc, &jpeg.Options{Quality: quality})
		if err != nil {
			return model.ResultedImage{}, err
		}
	case "png":
		// TODO: correct compress png
		if quality < 0 || quality > 255 {
			return model.ResultedImage{}, utils.ErrPNG
		}
		size := imgSrc.Bounds().Size()
		rect := image.Rect(0, 0, size.X, size.Y)
		wImg := image.NewRGBA(rect)
		for y := imgSrc.Bounds().Min.Y; y < imgSrc.Bounds().Max.Y; y++ {
			for x := imgSrc.Bounds().Min.X; x < imgSrc.Bounds().Max.X; x++ {
				pixel := imgSrc.At(x, y)
				originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)
				r := uint32(originalColor.R >> quality)
				g := uint32(originalColor.G >> quality)
				b := uint32(originalColor.B >> quality)

				newColor := uint8((r + g + b) / 3)
				c := color.RGBA{
					R: newColor, G: newColor, B: newColor, A: originalColor.A,
				}
				wImg.Set(x, y, c)
			}
		}
		err = png.Encode(newImgFile, wImg)
		if err != nil {
			return model.ResultedImage{}, err
		}
	}
	result.Name = uploadedImage.Name
	result.Location = currentDir + "\\results\\"
	return result, nil
}

// ConvertToType converts from png to jpeg and vice versa.
func (s *ImageService) ConvertToType(uploadedImage model.UploadedImage) (model.ResultedImage, error) {
	var result model.ResultedImage
	currentDir, _ := os.Getwd()

	// TODO: image location instead of --currentDir+ "\\uploads\\"--
	file, err := os.Open(currentDir + "\\uploads\\" + uploadedImage.Name)
	if err != nil {
		return model.ResultedImage{}, err
	}
	defer file.Close()
	img, format, err := image.Decode(file)
	if err != nil {
		return model.ResultedImage{}, err
	}
	convertedName, err := ChangeFormat(uploadedImage.Name)
	if err != nil {
		return model.ResultedImage{}, err
	}
	newImg, err := os.Create(fmt.Sprintf("./results/%s", convertedName))
	defer file.Close()
	if err != nil {
		return model.ResultedImage{}, utils.ErrCreateFile
	}
	defer newImg.Close()

	err = Encode(newImg, img, format)
	if err != nil {
		return model.ResultedImage{}, err
	}

	result.Name = convertedName
	result.Location = currentDir + "\\results\\"
	return result, nil
}

// CreateRequest creates request.
func (s *ImageService) CreateRequest(ctx context.Context, user model.User, uplImg model.UploadedImage, resImg model.ResultedImage, uI model.UserImage, r model.Request) (int, error) {
	return s.repo.CreateRequest(ctx, user, uplImg, resImg, uI, r)
}

// FindTheResultingImage finds the resulting image by id.
func (s *ImageService) FindTheResultingImage(ctx context.Context, id int, service model.Service) (model.ResultedImage, error) {
	return s.repo.FindTheResultingImage(ctx, id, service)
}

// FindOriginalImage finds original image by id.
func (s *ImageService) FindOriginalImage(ctx context.Context, id int) (model.UploadedImage, error) {
	return s.repo.FindOriginalImage(ctx, id)
}

// SaveImage saves image to users machine.
func (s *ImageService) SaveImage(filename, folder, resultedFilename string) error {
	currentDir, _ := os.Getwd()
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
