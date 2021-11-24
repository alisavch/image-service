// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"
	image "image"

	mock "github.com/stretchr/testify/mock"

	models "github.com/alisavch/image-service/internal/models"

	os "os"

	uuid "github.com/google/uuid"
)

// Image is an autogenerated mock type for the Image type
type Image struct {
	mock.Mock
}

// ChangeFormat provides a mock function with given fields: filename
func (_m *Image) ChangeFormat(filename string) (string, error) {
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
func (_m *Image) CompressImage(width int, format string, resultedName string, img image.Image, newImg *os.File, storage string) (models.Image, error) {
	ret := _m.Called(width, format, resultedName, img, newImg, storage)

	var r0 models.Image
	if rf, ok := ret.Get(0).(func(int, string, string, image.Image, *os.File, string) models.Image); ok {
		r0 = rf(width, format, resultedName, img, newImg, storage)
	} else {
		r0 = ret.Get(0).(models.Image)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, string, string, image.Image, *os.File, string) error); ok {
		r1 = rf(width, format, resultedName, img, newImg, storage)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ConvertToType provides a mock function with given fields: format, resultedName, img, newImg, storage
func (_m *Image) ConvertToType(format string, resultedName string, img image.Image, newImg *os.File, storage string) (models.Image, error) {
	ret := _m.Called(format, resultedName, img, newImg, storage)

	var r0 models.Image
	if rf, ok := ret.Get(0).(func(string, string, image.Image, *os.File, string) models.Image); ok {
		r0 = rf(format, resultedName, img, newImg, storage)
	} else {
		r0 = ret.Get(0).(models.Image)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, image.Image, *os.File, string) error); ok {
		r1 = rf(format, resultedName, img, newImg, storage)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateRequest provides a mock function with given fields: ctx, user, img, req
func (_m *Image) CreateRequest(ctx context.Context, user models.User, img models.Image, req models.Request) (uuid.UUID, error) {
	ret := _m.Called(ctx, user, img, req)

	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func(context.Context, models.User, models.Image, models.Request) uuid.UUID); ok {
		r0 = rf(ctx, user, img, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.User, models.Image, models.Request) error); ok {
		r1 = rf(ctx, user, img, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FillInTheResultingImage provides a mock function with given fields: storage, resultedName, newImg
func (_m *Image) FillInTheResultingImage(storage string, resultedName string, newImg *os.File) (models.Image, error) {
	ret := _m.Called(storage, resultedName, newImg)

	var r0 models.Image
	if rf, ok := ret.Get(0).(func(string, string, *os.File) models.Image); ok {
		r0 = rf(storage, resultedName, newImg)
	} else {
		r0 = ret.Get(0).(models.Image)
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
func (_m *Image) FillInTheResultingImageForAWS(resultedName string) (models.Image, error) {
	ret := _m.Called(resultedName)

	var r0 models.Image
	if rf, ok := ret.Get(0).(func(string) models.Image); ok {
		r0 = rf(resultedName)
	} else {
		r0 = ret.Get(0).(models.Image)
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
func (_m *Image) FindOriginalImage(ctx context.Context, id uuid.UUID) (models.Image, error) {
	ret := _m.Called(ctx, id)

	var r0 models.Image
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) models.Image); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(models.Image)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindResultedImage provides a mock function with given fields: ctx, id
func (_m *Image) FindResultedImage(ctx context.Context, id uuid.UUID) (models.Image, error) {
	ret := _m.Called(ctx, id)

	var r0 models.Image
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) models.Image); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(models.Image)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindUserRequestHistory provides a mock function with given fields: ctx, id
func (_m *Image) FindUserRequestHistory(ctx context.Context, id uuid.UUID) ([]models.History, error) {
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

// SaveImage provides a mock function with given fields: filename, location, storage
func (_m *Image) SaveImage(filename string, location string, storage string) (*models.SavedImage, error) {
	ret := _m.Called(filename, location, storage)

	var r0 *models.SavedImage
	if rf, ok := ret.Get(0).(func(string, string, string) *models.SavedImage); ok {
		r0 = rf(filename, location, storage)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SavedImage)
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

// SetCompletedTime provides a mock function with given fields: ctx, id
func (_m *Image) SetCompletedTime(ctx context.Context, id uuid.UUID) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateStatus provides a mock function with given fields: ctx, id, status
func (_m *Image) UpdateStatus(ctx context.Context, id uuid.UUID, status models.Status) error {
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
func (_m *Image) UploadImage(ctx context.Context, img models.Image) (uuid.UUID, error) {
	ret := _m.Called(ctx, img)

	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func(context.Context, models.Image) uuid.UUID); ok {
		r0 = rf(ctx, img)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.Image) error); ok {
		r1 = rf(ctx, img)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UploadResultedImage provides a mock function with given fields: ctx, img
func (_m *Image) UploadResultedImage(ctx context.Context, img models.Image) error {
	ret := _m.Called(ctx, img)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Image) error); ok {
		r0 = rf(ctx, img)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
