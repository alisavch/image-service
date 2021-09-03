package service

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
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

// EncodeOption sets an optional parameter for the Encode and Save functions.
type EncodeOption func(config *EncodeConfig)

// Encode encodes the image from jpeg to png and vice versa.
func Encode(w io.Writer, img image.Image, format string, opts ...EncodeOption) error {
	cfg := defaultEncodeConfig
	for _, option := range opts {
		option(&cfg)
	}
	switch format {
	case "jpeg":
		return png.Encode(w, img)
	case "png":
		var rgba *image.RGBA
		if nrgba, ok := img.(*image.NRGBA); ok {
			if nrgba.Opaque() {
				rgba = &image.RGBA{
					Pix:    nrgba.Pix,
					Stride: nrgba.Stride,
					Rect:   nrgba.Rect,
				}
			}
		}
		if rgba != nil {
			return jpeg.Encode(w, rgba, &jpeg.Options{Quality: cfg.jpegQuality})
		} else {
			return jpeg.Encode(w, img, &jpeg.Options{Quality: cfg.jpegQuality})
		}
	default:
		return utils.ErrUnsupportedFormat
	}
}

// SaveToDownloads saves the image to download folder on your computer.
func SaveToDownloads(name string) (*os.File, error) {
	userprofile := os.Getenv("HOME")
	return os.Create(userprofile + "/Downloads/" + name)
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
