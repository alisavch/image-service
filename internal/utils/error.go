package utils

import "errors"

var (
	// ErrFailedConvert checks the correctness of the conversion.
	ErrFailedConvert = errors.New("failed convert to int userID")
	// ErrRequest checks the correctness of the request.
	ErrRequest = errors.New("incorrect request")
	// ErrEmptyHeader checks empty header.
	ErrEmptyHeader = errors.New("auth header is empty")
	// ErrInvalidAuthHeader checks invalid header.
	ErrInvalidAuthHeader = errors.New("auth header is invalid")
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
	// ErrUnsupportedFormat checks supported formats.
	ErrUnsupportedFormat = errors.New("unsupported file format")
	// ErrFindImage the correctness of finding the image is checked.
	ErrFindImage = errors.New("cannot find image")
	// ErrSaveImage the correctness of saving the image is checked.
	ErrSaveImage = errors.New("cannot save image")
	// ErrConvert the correctness of converting the image is checked.
	ErrConvert = errors.New("cannot convert")
	// ErrMultipartForm the correctness of parse multipart form.
	ErrMultipartForm = errors.New("cannot parse multipart form")
	// ErrEmptyUsername checks username.
	ErrEmptyUsername = errors.New("username must not be empty")
	// ErrEmptyPassword checks password.
	ErrEmptyPassword = errors.New("password must not be empty")
	// ErrIncorrectRatio checks input ratio.
	ErrIncorrectRatio = errors.New("it is only possible to compress the image. enter correct size")
	// ErrFindVariable finds variable in .env file.
	ErrFindVariable = errors.New("could not find variable")
	// ErrMissingParams checks id in params.
	ErrMissingParams = errors.New("id is missing in parameters")
	// ErrPrivacy checks for equality of user ids from context and from parameters.
	ErrPrivacy = errors.New("you can only view your data")
	// ErrAtoi checks to convert to type int.
	ErrAtoi = errors.New("int conversion error")
)
