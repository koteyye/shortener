// Code generated by MockGen. DO NOT EDIT.
// Source: storage.go

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	models "github.com/koteyye/shortener/internal/app/models"
)

// MockURLStorage is a mock of URLStorage interface.
type MockURLStorage struct {
	ctrl     *gomock.Controller
	recorder *MockURLStorageMockRecorder
}

// MockURLStorageMockRecorder is the mock recorder for MockURLStorage.
type MockURLStorageMockRecorder struct {
	mock *MockURLStorage
}

// NewMockURLStorage creates a new mock instance.
func NewMockURLStorage(ctrl *gomock.Controller) *MockURLStorage {
	mock := &MockURLStorage{ctrl: ctrl}
	mock.recorder = &MockURLStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockURLStorage) EXPECT() *MockURLStorageMockRecorder {
	return m.recorder
}

// AddURL mocks base method.
func (m *MockURLStorage) AddURL(arg0 context.Context, arg1, arg2, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddURL", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddURL indicates an expected call of AddURL.
func (mr *MockURLStorageMockRecorder) AddURL(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddURL", reflect.TypeOf((*MockURLStorage)(nil).AddURL), arg0, arg1, arg2, arg3)
}

// DeleteURLByUser mocks base method.
func (m *MockURLStorage) DeleteURLByUser(arg0 context.Context, arg1 chan string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteURLByUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteURLByUser indicates an expected call of DeleteURLByUser.
func (mr *MockURLStorageMockRecorder) DeleteURLByUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteURLByUser", reflect.TypeOf((*MockURLStorage)(nil).DeleteURLByUser), arg0, arg1)
}

// GetDBPing mocks base method.
func (m *MockURLStorage) GetDBPing(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDBPing", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetDBPing indicates an expected call of GetDBPing.
func (mr *MockURLStorageMockRecorder) GetDBPing(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDBPing", reflect.TypeOf((*MockURLStorage)(nil).GetDBPing), ctx)
}

// GetShortURL mocks base method.
func (m *MockURLStorage) GetShortURL(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShortURL", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShortURL indicates an expected call of GetShortURL.
func (mr *MockURLStorageMockRecorder) GetShortURL(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShortURL", reflect.TypeOf((*MockURLStorage)(nil).GetShortURL), arg0, arg1)
}

// GetURL mocks base method.
func (m *MockURLStorage) GetURL(arg0 context.Context, arg1 string) (*models.SingleURL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURL", arg0, arg1)
	ret0, _ := ret[0].(*models.SingleURL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURL indicates an expected call of GetURL.
func (mr *MockURLStorageMockRecorder) GetURL(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURL", reflect.TypeOf((*MockURLStorage)(nil).GetURL), arg0, arg1)
}

// GetURLByUser mocks base method.
func (m *MockURLStorage) GetURLByUser(arg0 context.Context, arg1 string) ([]*models.URLList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURLByUser", arg0, arg1)
	ret0, _ := ret[0].([]*models.URLList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURLByUser indicates an expected call of GetURLByUser.
func (mr *MockURLStorageMockRecorder) GetURLByUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURLByUser", reflect.TypeOf((*MockURLStorage)(nil).GetURLByUser), arg0, arg1)
}