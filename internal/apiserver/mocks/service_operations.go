// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"
	image "image"

	io "io"

	mock "github.com/stretchr/testify/mock"

	models "github.com/alisavch/image-service/internal/models"

	os "os"

	uuid "github.com/google/uuid"
)

// ServiceOperations is an autogenerated mock type for the ServiceOperations type
type ServiceOperations struct {
	mock.Mock
}

// ChangeFormat provides a mock function with given fields: filename
func (_m *ServiceOperations) ChangeFormat(filename string) (string, error) {
	ret := _m.Called(filename)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(filename)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(filename)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CompressImage provides a mock function with given fields: width, format, resultedName, img, newImg, storage
func (_m *ServiceOperations) CompressImage(width int, format string, resultedName string, img image.Image, newImg *os.File, storage string) (models.ResultedImage, error) {
	ret := _m.Called(width, format, resultedName, img, newImg, storage)

	var r0 models.ResultedImage
	if rf, ok := ret.Get(0).(func(int, string, string, image.Image, *os.File, string) models.ResultedImage); ok {
		r0 = rf(width, format, resultedName, img, newImg, storage)
	} else {
		r0 = ret.Get(0).(models.ResultedImage)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, string, string, image.Image, *os.File, string) error); ok {
		r1 = rf(width, format, resultedName, img, newImg, storage)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ConvertToType provides a mock function with given fields: format, newImageName, img, newImg, storage
func (_m *ServiceOperations) ConvertToType(format string, newImageName string, img image.Image, newImg *os.File, storage string) (models.ResultedImage, error) {
	ret := _m.Called(format, newImageName, img, newImg, storage)

	var r0 models.ResultedImage
	if rf, ok := ret.Get(0).(func(string, string, image.Image, *os.File, string) models.ResultedImage); ok {
		r0 = rf(format, newImageName, img, newImg, storage)
	} else {
		r0 = ret.Get(0).(models.ResultedImage)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, image.Image, *os.File, string) error); ok {
		r1 = rf(format, newImageName, img, newImg, storage)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateRequest provides a mock function with given fields: ctx, user, uplImg, resImg, uI, r
func (_m *ServiceOperations) CreateRequest(ctx context.Context, user models.User, uplImg models.UploadedImage, resImg models.ResultedImage, uI models.UserImage, r models.Request) (uuid.UUID, error) {
	ret := _m.Called(ctx, user, uplImg, resImg, uI, r)

	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func(context.Context, models.User, models.UploadedImage, models.ResultedImage, models.UserImage, models.Request) uuid.UUID); ok {
		r0 = rf(ctx, user, uplImg, resImg, uI, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.User, models.UploadedImage, models.ResultedImage, models.UserImage, models.Request) error); ok {
		r1 = rf(ctx, user, uplImg, resImg, uI, r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateUser provides a mock function with given fields: ctx, user
func (_m *ServiceOperations) CreateUser(ctx context.Context, user models.User) (uuid.UUID, error) {
	ret := _m.Called(ctx, user)

	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func(context.Context, models.User) uuid.UUID); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DownloadFromS3Bucket provides a mock function with given fields: filename
func (_m *ServiceOperations) DownloadFromS3Bucket(filename string) (*os.File, error) {
	ret := _m.Called(filename)

	var r0 *os.File
	if rf, ok := ret.Get(0).(func(string) *os.File); ok {
		r0 = rf(filename)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*os.File)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(filename)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FillInTheResultingImage provides a mock function with given fields: storage, resultedName, newImg
func (_m *ServiceOperations) FillInTheResultingImage(storage string, resultedName string, newImg *os.File) (models.ResultedImage, error) {
	ret := _m.Called(storage, resultedName, newImg)

	var r0 models.ResultedImage
	if rf, ok := ret.Get(0).(func(string, string, *os.File) models.ResultedImage); ok {
		r0 = rf(storage, resultedName, newImg)
	} else {
		r0 = ret.Get(0).(models.ResultedImage)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, *os.File) error); ok {
		r1 = rf(storage, resultedName, newImg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FillInTheResultingImageForAWS provides a mock function with given fields: resultedName
func (_m *ServiceOperations) FillInTheResultingImageForAWS(resultedName string) (models.ResultedImage, error) {
	ret := _m.Called(resultedName)

	var r0 models.ResultedImage
	if rf, ok := ret.Get(0).(func(string) models.ResultedImage); ok {
		r0 = rf(resultedName)
	} else {
		r0 = ret.Get(0).(models.ResultedImage)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(resultedName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOriginalImage provides a mock function with given fields: ctx, id
func (_m *ServiceOperations) FindOriginalImage(ctx context.Context, id uuid.UUID) (models.UploadedImage, error) {
	ret := _m.Called(ctx, id)

	var r0 models.UploadedImage
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) models.UploadedImage); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(models.UploadedImage)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindTheResultingImage provides a mock function with given fields: ctx, id, service
func (_m *ServiceOperations) FindTheResultingImage(ctx context.Context, id uuid.UUID, service models.Service) (models.ResultedImage, error) {
	ret := _m.Called(ctx, id, service)

	var r0 models.ResultedImage
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.Service) models.ResultedImage); ok {
		r0 = rf(ctx, id, service)
	} else {
		r0 = ret.Get(0).(models.ResultedImage)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, models.Service) error); ok {
		r1 = rf(ctx, id, service)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindUserHistoryByID provides a mock function with given fields: ctx, id
func (_m *ServiceOperations) FindUserHistoryByID(ctx context.Context, id uuid.UUID) ([]models.History, error) {
	ret := _m.Called(ctx, id)

	var r0 []models.History
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []models.History); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.History)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateToken provides a mock function with given fields: ctx, username, password
func (_m *ServiceOperations) GenerateToken(ctx context.Context, username string, password string) (string, error) {
	ret := _m.Called(ctx, username, password)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, username, password)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, username, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ParseToken provides a mock function with given fields: token
func (_m *ServiceOperations) ParseToken(token string) (uuid.UUID, error) {
	ret := _m.Called(token)

	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func(string) uuid.UUID); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveImage provides a mock function with given fields: filename, location, storage
func (_m *ServiceOperations) SaveImage(filename string, location string, storage string) (*models.Image, error) {
	ret := _m.Called(filename, location, storage)

	var r0 *models.Image
	if rf, ok := ret.Get(0).(func(string, string, string) *models.Image); ok {
		r0 = rf(filename, location, storage)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Image)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(filename, location, storage)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateStatus provides a mock function with given fields: ctx, id, status
func (_m *ServiceOperations) UpdateStatus(ctx context.Context, id uuid.UUID, status models.Status) error {
	ret := _m.Called(ctx, id, status)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.Status) error); ok {
		r0 = rf(ctx, id, status)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UploadImage provides a mock function with given fields: ctx, img
func (_m *ServiceOperations) UploadImage(ctx context.Context, img models.UploadedImage) (uuid.UUID, error) {
	ret := _m.Called(ctx, img)

	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func(context.Context, models.UploadedImage) uuid.UUID); ok {
		r0 = rf(ctx, img)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.UploadedImage) error); ok {
		r1 = rf(ctx, img)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UploadToS3Bucket provides a mock function with given fields: file, filename
func (_m *ServiceOperations) UploadToS3Bucket(file io.Reader, filename string) (string, error) {
	ret := _m.Called(file, filename)

	var r0 string
	if rf, ok := ret.Get(0).(func(io.Reader, string) string); ok {
		r0 = rf(file, filename)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(io.Reader, string) error); ok {
		r1 = rf(file, filename)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
