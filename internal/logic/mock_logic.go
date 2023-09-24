package logic

import "github.com/golang/mock/gomock"

// MockLogic is a mock of Service interface.
type MockLogic struct {
	ctrl     *gomock.Controller
	recorder *MockLogicMockRecorder
}

// MockLogicMockRecorder is the mock recorder for MockLogic.
type MockLogicMockRecorder struct {
	mock *MockLogic
}

// NewMockLogic creates a new mock instance.
func NewMockLogic(ctrl *gomock.Controller) *MockLogic {
	mock := &MockLogic{ctrl: ctrl}
	mock.recorder = &MockLogicMockRecorder{mock}
	return mock
}

// MockUserLogic is a mock of Service interface.
type MockUserLogic struct {
	ctrl     *gomock.Controller
	recorder *MockUserLogicMockRecorder
}

// MockUserLogicMockRecorder is the mock recorder for MockUserLogic.
type MockUserLogicMockRecorder struct {
	mock *MockUserLogic
}

// NewMockUserLogic creates a new mock instance.
func NewMockUserLogic(ctrl *gomock.Controller) *MockUserLogic {
	mock := &MockUserLogic{ctrl: ctrl}
	mock.recorder = &MockUserLogicMockRecorder{mock}
	return mock
}
