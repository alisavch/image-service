package utils

import "errors"

var (
	// ErrGetUserID checks the correctness of the conversion.
	ErrGetUserID = errors.New("cannot get userID from context")
	// ErrRequest checks the correctness of the request.
	ErrRequest = errors.New("invalid path in request")
	// ErrGetJWTCookie checks jwt cookie.
	ErrGetJWTCookie = errors.New("cannot get jwt cookie")
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
	// ErrEmptyUsername checks username.
	ErrEmptyUsername = errors.New("username must not be empty")
	// ErrEmptyPassword checks password.
	ErrEmptyPassword = errors.New("password must not be empty")
	// ErrIncorrectRatio checks input ratio.
	ErrIncorrectRatio = errors.New("input ratio is incorrect. it should not exceed the image size")
	// ErrMissingParams checks id in params.
	ErrMissingParams = errors.New("id is missing in parameters")
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
	// ErrS3Uploading checks if the file can be uploaded to aws s3 bucket.
	ErrS3Uploading = errors.New("failed to upload file to S3 bucket")
	// ErrUserAlreadyExists checks the ability to create a user.
	ErrUserAlreadyExists = errors.New("user already exists")
	// ErrFindUser checks the ability to find the user.
	ErrFindUser = errors.New("cannot find the user in database")
	// ErrCreateQuery checks the possibility of creating a request.
	ErrCreateQuery = errors.New("cannot create a query")
	// ErrGetHistory checks the ability to get user history.
	ErrGetHistory = errors.New("cannot get history lines")
	// ErrUploadImageToDB verifies that the image information can be loaded into the database.
	ErrUploadImageToDB = errors.New("unable to insert image into database")
	// ErrFindTheResultingImage checks if the resulting image can be found.
	ErrFindTheResultingImage = errors.New("no such resulting image")
	// ErrFindOriginalImage checks if the original image can be found.
	ErrFindOriginalImage = errors.New("no such original image")
	// ErrUpdateStatusRequest checks if the status of the request can be updated.
	ErrUpdateStatusRequest = errors.New("cannot update image status")
	// ErrGenerateHash checks the hash generation capability.
	ErrGenerateHash = errors.New("cannot generate password hash")
	// ErrOpenFile checks if the file can be opened.
	ErrOpenFile = errors.New("unable to open file")
	// ErrRemoteUpload checks if the file can be loaded in aws.
	ErrRemoteUpload = errors.New("failed to upload to aws")
	// ErrRemoteDownload checks the ability to boot from the s3 bucket.
	ErrRemoteDownload = errors.New("cannot download from s3 bucket")
	// ErrGetContentType checks the ability to get the content type of the file.
	ErrGetContentType = errors.New("cannot get file content type")
	// ErrRowsAffected checks affected rows.
	ErrRowsAffected = errors.New("cannot return the number of affected rows")
	// ErrExpectedAffected checks expected affected rows.
	ErrExpectedAffected = errors.New("expected to affect 1 row, affected")
	// ErrSetCompletedTime verifies the ability to set a processing completion time.
	ErrSetCompletedTime = errors.New("cannot set the processing completion time")
	// ErrGetStatus checks request status.
	ErrGetStatus = errors.New("cannot find status for this request")
	// ErrImageProcessing checks processing status.
	ErrImageProcessing = errors.New("the image is being processed at the moment")
)
