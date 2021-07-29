// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	model "github.com/alisavch/image-service/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// Authorization is an autogenerated mock type for the Authorization type
type Authorization struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: user
func (_m *Authorization) CreateUser(user model.User) (int, error) {
	ret := _m.Called(user)

	var r0 int
	if rf, ok := ret.Get(0).(func(model.User) int); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(model.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateToken provides a mock function with given fields: username, password
func (_m *Authorization) GenerateToken(username string, password string) (string, error) {
	ret := _m.Called(username, password)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(username, password)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(username, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ParseToken provides a mock function with given fields: token
func (_m *Authorization) ParseToken(token string) (int, error) {
	ret := _m.Called(token)

	var r0 int
	if rf, ok := ret.Get(0).(func(string) int); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
