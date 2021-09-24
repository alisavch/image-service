package utils

import "errors"

var (
	// ErrGetUserID checks the correctness of the conversion.
	ErrGetUserID = errors.New("cannot get userID from context")
	// ErrRequest checks the correctness of the request.
	ErrRequest = errors.New("invalid path in request")
	// ErrEmptyHeader checks empty header.
	ErrEmptyHeader = errors.New("auth header is empty")
	// ErrInvalidAuthHeader checks invalid header.
	ErrInvalidAuthHeader = errors.New("auth header is invalid")
	// ErrEmptyToken checks token.
	ErrEmptyToken = errors.New("token is empty")
	// ErrUpload checks file upload.
	ErrUpload = errors.New("cannot upload the file")
	// ErrCreateFile checks the correctness of the file.
	ErrCreateFile = errors.New("cannot create new file")
	// ErrCopyFile checks copy file.
	ErrCopyFile = errors.New("cannot copy file")
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
	ErrIncorrectRatio = errors.New("input ratio is incorrect. it should not exceed the image size")
	// ErrMissingParams checks id in params.
	ErrMissingParams = errors.New("id is missing in parameters")
	// ErrPrivacy checks for equality of user ids from context and from parameters.
	ErrPrivacy = errors.New("users IDs do not match")
	// ErrAtoi checks to convert to type int.
	ErrAtoi = errors.New("cannot convert string to int")
	// ErrGetDir finds current location.
	ErrGetDir = errors.New("cannot get current directory")
	// ErrOpen opens the image.
	ErrOpen = errors.New("cannot open image")
	// ErrDecode decodes the image.
	ErrDecode = errors.New("cannot decode image")
	// ErrEnsureDir checks base directory.
	ErrEnsureDir = errors.New("cannot ensure base directory")
	// ErrCompress checks to compress the image.
	ErrCompress = errors.New("cannot compress")
	// ErrFileStat checks to get information about the file.
	ErrFileStat = errors.New("cannot get file info")
	// ErrCreateRequest verifies the execution of the request.
	ErrCreateRequest = errors.New("cannot create request")
)
