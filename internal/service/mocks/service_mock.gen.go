// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/service.go
//
// Generated by this command:
//
//	mockgen -source=internal/service/service.go -destination=internal/service/mocks/service_mock.gen.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	models "shortener/internal/models"

	gomock "go.uber.org/mock/gomock"
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

// BatchSave mocks base method.
func (m *MockURLStorage) BatchSave(ctx context.Context, input models.BatchArray) (models.BatchArray, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchSave", ctx, input)
	ret0, _ := ret[0].(models.BatchArray)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BatchSave indicates an expected call of BatchSave.
func (mr *MockURLStorageMockRecorder) BatchSave(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchSave", reflect.TypeOf((*MockURLStorage)(nil).BatchSave), ctx, input)
}

// Cleanup mocks base method.
func (m *MockURLStorage) Cleanup(ctx context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cleanup", ctx)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Cleanup indicates an expected call of Cleanup.
func (mr *MockURLStorageMockRecorder) Cleanup(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cleanup", reflect.TypeOf((*MockURLStorage)(nil).Cleanup), ctx)
}

// Close mocks base method.
func (m *MockURLStorage) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockURLStorageMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockURLStorage)(nil).Close))
}

// DeleteURLs mocks base method.
func (m *MockURLStorage) DeleteURLs(ctx context.Context, input models.DeleteURLs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteURLs", ctx, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteURLs indicates an expected call of DeleteURLs.
func (mr *MockURLStorageMockRecorder) DeleteURLs(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteURLs", reflect.TypeOf((*MockURLStorage)(nil).DeleteURLs), ctx, input)
}

// Get mocks base method.
func (m *MockURLStorage) Get(ctx context.Context, shortLink string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, shortLink)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockURLStorageMockRecorder) Get(ctx, shortLink any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockURLStorage)(nil).Get), ctx, shortLink)
}

// GetByUserID mocks base method.
func (m *MockURLStorage) GetByUserID(ctx context.Context) ([]models.BaseRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUserID", ctx)
	ret0, _ := ret[0].([]models.BaseRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUserID indicates an expected call of GetByUserID.
func (mr *MockURLStorageMockRecorder) GetByUserID(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUserID", reflect.TypeOf((*MockURLStorage)(nil).GetByUserID), ctx)
}

// Ping mocks base method.
func (m *MockURLStorage) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockURLStorageMockRecorder) Ping(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockURLStorage)(nil).Ping), ctx)
}

// Save mocks base method.
func (m *MockURLStorage) Save(ctx context.Context, shortLink, longLink string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, shortLink, longLink)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockURLStorageMockRecorder) Save(ctx, shortLink, longLink any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockURLStorage)(nil).Save), ctx, shortLink, longLink)
}
