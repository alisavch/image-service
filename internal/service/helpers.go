package service

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"

	"github.com/alisavch/image-service/internal/utils"
)

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
	var err error
	switch format {
	case "jpeg":
		err = png.Encode(w, img)
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
			err = jpeg.Encode(w, rgba, &jpeg.Options{Quality: cfg.jpegQuality})
		} else {
			err = jpeg.Encode(w, img, &jpeg.Options{Quality: cfg.jpegQuality})
		}
	default:
		err = utils.ErrUnsupportedFormat
	}
	return err
}

// SaveToDownloads saves the image to download folder on your computer.
func SaveToDownloads(name string) (*os.File, error) {
	userprofile := os.Getenv("USERPROFILE")
	out, err := os.Create(userprofile + "\\Downloads\\" + name)
	if err != nil {
		return nil, err
	}
	return out, nil
}
