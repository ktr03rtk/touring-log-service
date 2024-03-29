// Code generated by MockGen. DO NOT EDIT.
// Source: athena_query_repository.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAthenaQueryRepository is a mock of AthenaQueryRepository interface.
type MockAthenaQueryRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAthenaQueryRepositoryMockRecorder
}

// MockAthenaQueryRepositoryMockRecorder is the mock recorder for MockAthenaQueryRepository.
type MockAthenaQueryRepositoryMockRecorder struct {
	mock *MockAthenaQueryRepository
}

// NewMockAthenaQueryRepository creates a new mock instance.
func NewMockAthenaQueryRepository(ctrl *gomock.Controller) *MockAthenaQueryRepository {
	mock := &MockAthenaQueryRepository{ctrl: ctrl}
	mock.recorder = &MockAthenaQueryRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAthenaQueryRepository) EXPECT() *MockAthenaQueryRepositoryMockRecorder {
	return m.recorder
}

// Fetch mocks base method.
func (m *MockAthenaQueryRepository) Fetch(ctx context.Context, rawQuery string, args []interface{}) ([][]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fetch", ctx, rawQuery, args)
	ret0, _ := ret[0].([][]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Fetch indicates an expected call of Fetch.
func (mr *MockAthenaQueryRepositoryMockRecorder) Fetch(ctx, rawQuery, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockAthenaQueryRepository)(nil).Fetch), ctx, rawQuery, args)
}
