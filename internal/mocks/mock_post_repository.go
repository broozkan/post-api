// Code generated by MockGen. DO NOT EDIT.
// Source: broozkan/postapi/internal/services (interfaces: RepositoryInterface)

// Package mocks is a generated GoMock package.
package mocks

import (
	models "broozkan/postapi/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRepositoryInterface is a mock of RepositoryInterface interface.
type MockRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryInterfaceMockRecorder
}

// MockRepositoryInterfaceMockRecorder is the mock recorder for MockRepositoryInterface.
type MockRepositoryInterfaceMockRecorder struct {
	mock *MockRepositoryInterface
}

// NewMockRepositoryInterface creates a new mock instance.
func NewMockRepositoryInterface(ctrl *gomock.Controller) *MockRepositoryInterface {
	mock := &MockRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepositoryInterface) EXPECT() *MockRepositoryInterfaceMockRecorder {
	return m.recorder
}

// CreatePost mocks base method.
func (m *MockRepositoryInterface) CreatePost(arg0 *models.Post) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePost", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreatePost indicates an expected call of CreatePost.
func (mr *MockRepositoryInterfaceMockRecorder) CreatePost(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePost", reflect.TypeOf((*MockRepositoryInterface)(nil).CreatePost), arg0)
}
