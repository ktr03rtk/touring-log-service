// Code generated by MockGen. DO NOT EDIT.
// Source: payload_subscribe_repository.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/ktr03rtk/touring-log-service/data-store/domain/model"
)

// MockPayloadSubscribeRepository is a mock of PayloadSubscribeRepository interface.
type MockPayloadSubscribeRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPayloadSubscribeRepositoryMockRecorder
}

// MockPayloadSubscribeRepositoryMockRecorder is the mock recorder for MockPayloadSubscribeRepository.
type MockPayloadSubscribeRepositoryMockRecorder struct {
	mock *MockPayloadSubscribeRepository
}

// NewMockPayloadSubscribeRepository creates a new mock instance.
func NewMockPayloadSubscribeRepository(ctrl *gomock.Controller) *MockPayloadSubscribeRepository {
	mock := &MockPayloadSubscribeRepository{ctrl: ctrl}
	mock.recorder = &MockPayloadSubscribeRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPayloadSubscribeRepository) EXPECT() *MockPayloadSubscribeRepositoryMockRecorder {
	return m.recorder
}

// Subscribe mocks base method.
func (m *MockPayloadSubscribeRepository) Subscribe(arg0 context.Context, arg1 <-chan *model.Payload) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockPayloadSubscribeRepositoryMockRecorder) Subscribe(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockPayloadSubscribeRepository)(nil).Subscribe), arg0, arg1)
}
