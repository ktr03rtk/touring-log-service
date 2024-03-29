// Code generated by MockGen. DO NOT EDIT.
// Source: photo_image_repository.go

// Package mock is a generated GoMock package.
package mock

import (
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPhotoImageRepository is a mock of PhotoImageRepository interface.
type MockPhotoImageRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPhotoImageRepositoryMockRecorder
}

// MockPhotoImageRepositoryMockRecorder is the mock recorder for MockPhotoImageRepository.
type MockPhotoImageRepositoryMockRecorder struct {
	mock *MockPhotoImageRepository
}

// NewMockPhotoImageRepository creates a new mock instance.
func NewMockPhotoImageRepository(ctrl *gomock.Controller) *MockPhotoImageRepository {
	mock := &MockPhotoImageRepository{ctrl: ctrl}
	mock.recorder = &MockPhotoImageRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPhotoImageRepository) EXPECT() *MockPhotoImageRepositoryMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockPhotoImageRepository) Get(key string) (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockPhotoImageRepositoryMockRecorder) Get(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockPhotoImageRepository)(nil).Get), key)
}

// Store mocks base method.
func (m *MockPhotoImageRepository) Store(arg0 io.Reader, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Store indicates an expected call of Store.
func (mr *MockPhotoImageRepositoryMockRecorder) Store(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockPhotoImageRepository)(nil).Store), arg0, arg1)
}
