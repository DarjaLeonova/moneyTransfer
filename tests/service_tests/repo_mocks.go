package service_tests

import (
	"context"
	"github.com/stretchr/testify/mock"
	"moneyTransfer/internal/domain/model"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetBalance(ctx context.Context, userId string) (float64, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockUserRepo) GetById(ctx context.Context, userId string) (model.User, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepo) UpdateBalance(ctx context.Context, userID string, newBalance float64) error {
	args := m.Called(ctx, userID, newBalance)
	return args.Error(0)
}

type MockTransferRepo struct {
	mock.Mock
}

func (m *MockTransferRepo) GetTransactionsByUserId(ctx context.Context, userId string) ([]model.Transaction, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).([]model.Transaction), args.Error(1)
}

func (m *MockTransferRepo) CreateTransfer(ctx context.Context, tx model.Transaction) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *MockTransferRepo) UpdateTransactionStatus(ctx context.Context, txID, status string) error {
	args := m.Called(ctx, txID, status)
	return args.Error(0)
}
