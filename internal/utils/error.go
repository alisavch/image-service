package utils

import "errors"

var (
	// ErrWrong checks expected result.
	ErrWrong = errors.New("something went wrong")
	// ErrPrivileges checks user privileges.
	ErrPrivileges = errors.New("invalid user privileges")
	// ErrFailedConvert checks the correctness of the conversion.
	ErrFailedConvert = errors.New("failed convert to int userID")
	// ErrRequest checks the correctness of the request.
	ErrRequest = errors.New("incorrect request")
	// ErrEmptyHeader checks empty header.
	ErrEmptyHeader = errors.New("auth header is empty")
	// ErrInvalidHeader checks invalid header.
	ErrInvalidHeader = errors.New("auth header is invalid")
	// ErrEmptyToken checks token.
	ErrEmptyToken = errors.New("token is empty")
	// ErrUpload checks file upload.
	ErrUpload = errors.New("error upload file")
	// ErrCreateFile checks the correctness of the file.
	ErrCreateFile = errors.New("error create new file")
	// ErrCopyFile checks copy file.
	ErrCopyFile = errors.New("error copy file")
	// ErrAllowedFormat checks allowed format of the file.
	ErrAllowedFormat = errors.New("file format is not allowed. Please upload a JPEG or PNG")
	// ErrSigningMethod checks signing method.
	ErrSigningMethod = errors.New("invalid signing method")
	// ErrInvalidToken checks token.
	ErrInvalidToken = errors.New("token claims is invalid")
	// ErrJPEG checks jpeg quality.
	ErrJPEG = errors.New("jpeg support only quality from 0 to 100")
	// ErrPNG checks png quality.
	ErrPNG = errors.New("png support only quality from 0 to 255")
	// ErrUnsupportedFormat checks supported formats.
	ErrUnsupportedFormat = errors.New("unsupported file format")
	// ErrFindImage the correctness of finding the image is checked.
	ErrFindImage = errors.New("cannot find image")
	// ErrSaveImage the correctness of saving the image is checked.
	ErrSaveImage = errors.New("cannot save image")
	// ErrConvert the correctness of converting the image is checked.
	ErrConvert = errors.New("cannot convert")
	// ErrCreateRequest the correctness of create request with image.
	ErrCreateRequest = errors.New("cannot create request with image")
)
