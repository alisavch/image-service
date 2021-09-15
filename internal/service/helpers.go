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
	"strings"

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

// ChangeFormat changes image format.
func ChangeFormat(filename string) (string, error) {
	imgNames := strings.Split(strings.ToLower(filename), ".")
	extension := imgNames[len(imgNames)-1]
	if format, ok := convertedType[extension]; ok {
		imgNames[len(imgNames)-1] = format
		convertedName := strings.Join(imgNames, ".")
		return convertedName, nil
	}
	return "", utils.ErrUnsupportedFormat
}

// EncodeOption sets an optional parameter for to Encode and Save functions.
type EncodeOption func(config *EncodeConfig)

// ConvertToPNG converts from JPEG to PNG.
func ConvertToPNG(w io.Writer, imgSrc image.Image) error {
	var enc png.Encoder
	if err := enc.Encode(w, imgSrc); err != nil {
		return fmt.Errorf("%s to png", utils.ErrConvert)
	}
	return nil
}

// ConvertToJPEG convert from PNG to JPEG.
func ConvertToJPEG(w io.Writer, imgSrc image.Image) error {
	cfg := defaultEncodeConfig
	if err := jpeg.Encode(w, imgSrc, &jpeg.Options{Quality: cfg.pngCompressionLevel}); err != nil {
		return fmt.Errorf("%s to jpeg", utils.ErrConvert)
	}
	return nil
}

// CompressJPEG allows you to compress the JPEG image in width while maintaining the aspect ratio.
func CompressJPEG(imgSrc image.Image, width int, newImgFile *os.File) error {
	if width < 0 || width > imgSrc.Bounds().Max.X {
		return utils.ErrIncorrectRatio
	}
	m := resize.Resize(uint(width), 0, imgSrc, resize.Lanczos3)
	return jpeg.Encode(newImgFile, m, nil)
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

// GetFileContentType checks the content type of the file.
func GetFileContentType(out *os.File) (string, error) {
	buffer := make([]byte, 512)
	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(buffer)
	return contentType, nil
}
