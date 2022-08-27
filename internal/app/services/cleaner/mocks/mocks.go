// Code generated by MockGen. DO NOT EDIT.
// Source: cleaner.go

// Package mock_cleaner is a generated GoMock package.
package mock_cleaner

import (
	context "context"
	reflect "reflect"

	models "github.com/bgoldovsky/shortener/internal/app/models"
	gomock "github.com/golang/mock/gomock"
)

// MockurlsRepository is a mock of urlsRepository interface.
type MockurlsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockurlsRepositoryMockRecorder
}

// MockurlsRepositoryMockRecorder is the mock recorder for MockurlsRepository.
type MockurlsRepositoryMockRecorder struct {
	mock *MockurlsRepository
}

// NewMockurlsRepository creates a new mock instance.
func NewMockurlsRepository(ctrl *gomock.Controller) *MockurlsRepository {
	mock := &MockurlsRepository{ctrl: ctrl}
	mock.recorder = &MockurlsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockurlsRepository) EXPECT() *MockurlsRepositoryMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockurlsRepository) Delete(ctx context.Context, urlsBatch []models.UserCollection) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, urlsBatch)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockurlsRepositoryMockRecorder) Delete(ctx, urlsBatch interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockurlsRepository)(nil).Delete), ctx, urlsBatch)
}
