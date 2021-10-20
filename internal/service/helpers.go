package service

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/alisavch/image-service/internal/models"

	"github.com/alisavch/image-service/internal/utils"

	"github.com/nfnt/resize"
)

// EncodeConfig contains image constants.
type EncodeConfig struct {
	jpegQuality         int
	pngCompressionLevel int
}

var defaultEncodeConfig = EncodeConfig{
	jpegQuality:         95,
	pngCompressionLevel: 0,
}

// EncodeOption sets an optional parameter for to Encode and Save functions.
type EncodeOption func(config *EncodeConfig)

// ConvertToPNG converts from JPEG to PNG.
func ConvertToPNG(w io.Writer, imgSrc image.Image) error {
	var enc png.Encoder

	return enc.Encode(w, imgSrc)
}

// ConvertToJPEG convert from PNG to JPEG.
func ConvertToJPEG(w io.Writer, imgSrc image.Image) error {
	cfg := defaultEncodeConfig

	return jpeg.Encode(w, imgSrc, &jpeg.Options{Quality: cfg.pngCompressionLevel})
}

// CompressJPEG allows you to compress the JPEG image in width while maintaining the aspect ratio.
func CompressJPEG(imgSrc image.Image, width int, newImgFile *os.File) error {
	cfg := defaultEncodeConfig
	if width < 0 || width > imgSrc.Bounds().Max.X {
		return utils.ErrIncorrectRatio
	}

	m := resize.Resize(uint(width), 0, imgSrc, resize.Lanczos3)

	return jpeg.Encode(newImgFile, m, &jpeg.Options{Quality: cfg.jpegQuality})
}

// CompressPNG allows you to compress the PNG image in width while maintaining the aspect ratio.
func CompressPNG(imgSrc image.Image, width int, newImgFile *os.File) error {
	if width < 0 || width > imgSrc.Bounds().Max.X {
		return utils.ErrIncorrectRatio
	}

	m := resize.Resize(uint(width), 0, imgSrc, resize.Lanczos3)

	return png.Encode(newImgFile, m)
}

// EnsureBaseDir checks if a directory exists.
func EnsureBaseDir(filepath string) error {
	baseDir := path.Dir(filepath)

	info, err := os.Stat(baseDir)
	if err == nil && info.IsDir() {
		return nil
	}

	return os.MkdirAll(baseDir, 0755)
}

// GetFileContentType gets the content type of the file.
func GetFileContentType(out *os.File) (string, error) {
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

// GetFileSize gets the filesize of the file.
func GetFileSize(file *os.File) (int64, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return 0, utils.ErrFileStat
	}

	size := fileInfo.Size()

	_, err = file.Seek(0, 0)
	if err != nil {
		return 0, err
	}
	return size, nil
}

// FillInTheReceivedNameAndLocation fills name and location
func FillInTheReceivedNameAndLocation(name, location string) models.ResultedImage {
	var result models.ResultedImage

	result.Name = name
	result.Location = location

	return result
}

// FillInTheImage fills models.Image.
func FillInTheImage(img models.Image, file *os.File) (models.Image, error) {
	var err error
	img.File = file

	img.ContentType, err = GetFileContentType(file)
	if err != nil {
		return models.Image{}, fmt.Errorf("%s:%s", utils.ErrGetContentType, err)
	}

	img.Filesize, err = GetFileSize(file)
	if err != nil {
		return models.Image{}, err
	}

	return img, nil
}

// FillInTheResultingImageLocally fills the resultedImage with information for local storage.
func FillInTheResultingImageLocally(resultedName string, newImg *os.File) (models.ResultedImage, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return models.ResultedImage{}, utils.ErrGetDir
	}

	err = newImg.Close()
	if err != nil {
		return models.ResultedImage{}, err
	}

	result := FillInTheReceivedNameAndLocation(resultedName, currentDir+"/results/")

	return result, nil
}
