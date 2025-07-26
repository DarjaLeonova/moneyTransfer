package tests

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"moneyTransfer/internal/domain/model"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetBalance(ctx context.Context, userId string) (float64, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockUserService) GetById(ctx context.Context, userId string) (model.User, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(model.User), args.Error(1)
}

type MockTransferService struct {
	mock.Mock
}

func (m *MockTransferService) GetTransactionsByUserId(ctx context.Context, userId string) ([]model.Transaction, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).([]model.Transaction), args.Error(1)
}

func (m *MockTransferService) CreateTransfer(ctx context.Context, from, to string, amount float64) (uuid.UUID, error) {
	args := m.Called(ctx, from, to, amount)
	return args.Get(0).(uuid.UUID), args.Error(1)
}
