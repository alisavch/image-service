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

// CompressImage provides a mock function with given fields: width, format, resultedName, img, newImg, isRemoteStorage
func (_m *Image) CompressImage(width int, format string, resultedName string, img image.Image, newImg *os.File, isRemoteStorage bool) (models.ResultedImage, error) {
	ret := _m.Called(width, format, resultedName, img, newImg, isRemoteStorage)

	var r0 models.ResultedImage
	if rf, ok := ret.Get(0).(func(int, string, string, image.Image, *os.File, bool) models.ResultedImage); ok {
		r0 = rf(width, format, resultedName, img, newImg, isRemoteStorage)
	} else {
		r0 = ret.Get(0).(models.ResultedImage)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, string, string, image.Image, *os.File, bool) error); ok {
		r1 = rf(width, format, resultedName, img, newImg, isRemoteStorage)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ConvertToType provides a mock function with given fields: format, newImageName, img, newImg, isRemoteStorage
func (_m *Image) ConvertToType(format string, newImageName string, img image.Image, newImg *os.File, isRemoteStorage bool) (models.ResultedImage, error) {
	ret := _m.Called(format, newImageName, img, newImg, isRemoteStorage)

	var r0 models.ResultedImage
	if rf, ok := ret.Get(0).(func(string, string, image.Image, *os.File, bool) models.ResultedImage); ok {
		r0 = rf(format, newImageName, img, newImg, isRemoteStorage)
	} else {
		r0 = ret.Get(0).(models.ResultedImage)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, image.Image, *os.File, bool) error); ok {
		r1 = rf(format, newImageName, img, newImg, isRemoteStorage)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateRequest provides a mock function with given fields: ctx, user, uplImg, resImg, uI, r
func (_m *Image) CreateRequest(ctx context.Context, user models.User, uplImg models.UploadedImage, resImg models.ResultedImage, uI models.UserImage, r models.Request) (uuid.UUID, error) {
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

// FindOriginalImage provides a mock function with given fields: ctx, id
func (_m *Image) FindOriginalImage(ctx context.Context, id uuid.UUID) (models.UploadedImage, error) {
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

// FindTheResultingImage provides a mock function with given fields: ctx, id, _a2
func (_m *Image) FindTheResultingImage(ctx context.Context, id uuid.UUID, _a2 models.Service) (models.ResultedImage, error) {
	ret := _m.Called(ctx, id, _a2)

	var r0 models.ResultedImage
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, models.Service) models.ResultedImage); ok {
		r0 = rf(ctx, id, _a2)
	} else {
		r0 = ret.Get(0).(models.ResultedImage)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, models.Service) error); ok {
		r1 = rf(ctx, id, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindUserHistoryByID provides a mock function with given fields: ctx, id
func (_m *Image) FindUserHistoryByID(ctx context.Context, id uuid.UUID) ([]models.History, error) {
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

// SaveImage provides a mock function with given fields: filename, location, isRemoteStorage
func (_m *Image) SaveImage(filename string, location string, isRemoteStorage bool) (*models.Image, error) {
	ret := _m.Called(filename, location, isRemoteStorage)

	var r0 *models.Image
	if rf, ok := ret.Get(0).(func(string, string, bool) *models.Image); ok {
		r0 = rf(filename, location, isRemoteStorage)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Image)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, bool) error); ok {
		r1 = rf(filename, location, isRemoteStorage)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
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
func (_m *Image) UploadImage(ctx context.Context, img models.UploadedImage) (uuid.UUID, error) {
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
