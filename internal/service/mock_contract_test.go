// Code generated by MockGen. DO NOT EDIT.
// Source: contract.go

// Package service is a generated GoMock package.
package service

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/s21platform/advert-service/internal/model"
)

// MockDBRepo is a mock of DBRepo interface.
type MockDBRepo struct {
	ctrl     *gomock.Controller
	recorder *MockDBRepoMockRecorder
}

// MockDBRepoMockRecorder is the mock recorder for MockDBRepo.
type MockDBRepoMockRecorder struct {
	mock *MockDBRepo
}

// NewMockDBRepo creates a new mock instance.
func NewMockDBRepo(ctrl *gomock.Controller) *MockDBRepo {
	mock := &MockDBRepo{ctrl: ctrl}
	mock.recorder = &MockDBRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDBRepo) EXPECT() *MockDBRepoMockRecorder {
	return m.recorder
}

// GetAdverts mocks base method.
func (m *MockDBRepo) GetAdverts(UUID string) (*model.AdvertInfoList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAdverts", UUID)
	ret0, _ := ret[0].(*model.AdvertInfoList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdverts indicates an expected call of GetAdverts.
func (mr *MockDBRepoMockRecorder) GetAdverts(UUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdverts", reflect.TypeOf((*MockDBRepo)(nil).GetAdverts), UUID)
}
