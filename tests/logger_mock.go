package tests

import "github.com/stretchr/testify/mock"

type MockLogger struct {
	mock.Mock
}

func (_m *MockLogger) Debug(msg string, args ...any) {
	args = append([]any{msg}, args...)
	_m.Called(args...)
}

func (_m *MockLogger) Info(msg string, args ...any) {
	args = append([]any{msg}, args...)
	_m.Called(args...)
}

func (_m *MockLogger) Warn(msg string, args ...any) {
	args = append([]any{msg}, args...)
	_m.Called(args...)
}

func (_m *MockLogger) Error(msg string, args ...any) {
	args = append([]any{msg}, args...)
	_m.Called(args...)
}
