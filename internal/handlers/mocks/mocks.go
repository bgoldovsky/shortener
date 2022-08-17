// Code generated by MockGen. DO NOT EDIT.
// Source: handlers.go

// Package mock_handlers is a generated GoMock package.
package mock_handlers

import (
	context "context"
	reflect "reflect"

	models "github.com/bgoldovsky/shortener/internal/app/models"
	gomock "github.com/golang/mock/gomock"
)

// MockurlsService is a mock of urlsService interface.
type MockurlsService struct {
	ctrl     *gomock.Controller
	recorder *MockurlsServiceMockRecorder
}

// MockurlsServiceMockRecorder is the mock recorder for MockurlsService.
type MockurlsServiceMockRecorder struct {
	mock *MockurlsService
}

// NewMockurlsService creates a new mock instance.
func NewMockurlsService(ctrl *gomock.Controller) *MockurlsService {
	mock := &MockurlsService{ctrl: ctrl}
	mock.recorder = &MockurlsServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockurlsService) EXPECT() *MockurlsServiceMockRecorder {
	return m.recorder
}

// Expand mocks base method.
func (m *MockurlsService) Expand(ctx context.Context, id string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Expand", ctx, id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Expand indicates an expected call of Expand.
func (mr *MockurlsServiceMockRecorder) Expand(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Expand", reflect.TypeOf((*MockurlsService)(nil).Expand), ctx, id)
}

// GetUrls mocks base method.
func (m *MockurlsService) GetUrls(ctx context.Context, userID string) ([]models.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUrls", ctx, userID)
	ret0, _ := ret[0].([]models.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUrls indicates an expected call of GetUrls.
func (mr *MockurlsServiceMockRecorder) GetUrls(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUrls", reflect.TypeOf((*MockurlsService)(nil).GetUrls), ctx, userID)
}

// Shorten mocks base method.
func (m *MockurlsService) Shorten(ctx context.Context, url, userID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Shorten", ctx, url, userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Shorten indicates an expected call of Shorten.
func (mr *MockurlsServiceMockRecorder) Shorten(ctx, url, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shorten", reflect.TypeOf((*MockurlsService)(nil).Shorten), ctx, url, userID)
}

// ShortenBatch mocks base method.
func (m *MockurlsService) ShortenBatch(ctx context.Context, originalURLs []models.OriginalURL, userID string) ([]models.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShortenBatch", ctx, originalURLs, userID)
	ret0, _ := ret[0].([]models.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ShortenBatch indicates an expected call of ShortenBatch.
func (mr *MockurlsServiceMockRecorder) ShortenBatch(ctx, originalURLs, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShortenBatch", reflect.TypeOf((*MockurlsService)(nil).ShortenBatch), ctx, originalURLs, userID)
}

// Mockauth is a mock of auth interface.
type Mockauth struct {
	ctrl     *gomock.Controller
	recorder *MockauthMockRecorder
}

// MockauthMockRecorder is the mock recorder for Mockauth.
type MockauthMockRecorder struct {
	mock *Mockauth
}

// NewMockauth creates a new mock instance.
func NewMockauth(ctrl *gomock.Controller) *Mockauth {
	mock := &Mockauth{ctrl: ctrl}
	mock.recorder = &MockauthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockauth) EXPECT() *MockauthMockRecorder {
	return m.recorder
}

// UserID mocks base method.
func (m *Mockauth) UserID(ctx context.Context) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserID", ctx)
	ret0, _ := ret[0].(string)
	return ret0
}

// UserID indicates an expected call of UserID.
func (mr *MockauthMockRecorder) UserID(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserID", reflect.TypeOf((*Mockauth)(nil).UserID), ctx)
}

// Mockinfra is a mock of infra interface.
type Mockinfra struct {
	ctrl     *gomock.Controller
	recorder *MockinfraMockRecorder
}

// MockinfraMockRecorder is the mock recorder for Mockinfra.
type MockinfraMockRecorder struct {
	mock *Mockinfra
}

// NewMockinfra creates a new mock instance.
func NewMockinfra(ctrl *gomock.Controller) *Mockinfra {
	mock := &Mockinfra{ctrl: ctrl}
	mock.recorder = &MockinfraMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockinfra) EXPECT() *MockinfraMockRecorder {
	return m.recorder
}

// Ping mocks base method.
func (m *Mockinfra) Ping(ctx context.Context) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockinfraMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*Mockinfra)(nil).Ping), ctx)
}

// Mockcleaner is a mock of cleaner interface.
type Mockcleaner struct {
	ctrl     *gomock.Controller
	recorder *MockcleanerMockRecorder
}

// MockcleanerMockRecorder is the mock recorder for Mockcleaner.
type MockcleanerMockRecorder struct {
	mock *Mockcleaner
}

// NewMockcleaner creates a new mock instance.
func NewMockcleaner(ctrl *gomock.Controller) *Mockcleaner {
	mock := &Mockcleaner{ctrl: ctrl}
	mock.recorder = &MockcleanerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockcleaner) EXPECT() *MockcleanerMockRecorder {
	return m.recorder
}

// Queue mocks base method.
func (m *Mockcleaner) Queue(urls models.UserCollection) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Queue", urls)
}

// Queue indicates an expected call of Queue.
func (mr *MockcleanerMockRecorder) Queue(urls interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Queue", reflect.TypeOf((*Mockcleaner)(nil).Queue), urls)
}
