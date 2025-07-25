package service_tests

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"moneyTransfer/internal/domain/model"
	"moneyTransfer/internal/domain/service"
	"moneyTransfer/internal/queue"
	"moneyTransfer/tests"
	"testing"
	"time"
)

func TestTransferService_GetTransactionsByUserId_Success(t *testing.T) {
	ctx := context.Background()
	userRepo := new(tests.MockUserRepo)
	transferRepo := new(tests.MockTransferRepo)
	svc := service.NewTransferService(transferRepo, userRepo)

	expectedTransactions := []model.Transaction{
		{
			Id:         uuid.MustParse("a5bceab4-9dab-4d7a-8cd5-4ba832ebf899"),
			SenderId:   uuid.MustParse("7141b92f-a8c8-471e-83e5-7fc72da61cb9"),
			ReceiverId: uuid.MustParse("9d02adbc-27ca-4695-9d92-10cb35db67f4"),
			Amount:     100,
			Status:     model.StatusSuccess,
			CreatedAt:  time.Time{},
		},
	}

	transferRepo.On("GetTransactionsByUserId", ctx, "7141b92f-a8c8-471e-83e5-7fc72da61cb9").Return(expectedTransactions, nil)

	transactions, err := svc.GetTransactionsByUserId(ctx, "7141b92f-a8c8-471e-83e5-7fc72da61cb9")
	assert.NoError(t, err)
	assert.Equal(t, expectedTransactions, transactions)

	transferRepo.AssertExpectations(t)
}

func TestTransferService_GetTransactionsByUserId_Error(t *testing.T) {
	ctx := context.Background()
	userRepo := new(tests.MockUserRepo)
	transferRepo := new(tests.MockTransferRepo)
	svc := service.NewTransferService(transferRepo, userRepo)

	transferRepo.On("GetTransactionsByUserId", ctx, "7141b92f-a8c8-471e-83e5-7fc72da61cb9").
		Return(make([]model.Transaction, 0), errors.New("db error"))

	transactions, err := svc.GetTransactionsByUserId(ctx, "7141b92f-a8c8-471e-83e5-7fc72da61cb9")
	assert.Error(t, err)
	assert.Len(t, transactions, 0)

	transferRepo.AssertExpectations(t)
}

func TestTransferService_CreateTransfer_AmountLessOrEqualZero(t *testing.T) {
	ctx := context.Background()
	transferRepo := new(tests.MockTransferRepo)
	userRepo := new(tests.MockUserRepo)
	svc := service.NewTransferService(transferRepo, userRepo)

	fromID := uuid.New().String()
	toID := uuid.New().String()

	id, err := svc.CreateTransfer(ctx, fromID, toID, 0)
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
}

func TestTransferService_TransferService_CreateTransfer_RepoError(t *testing.T) {
	ctx := context.Background()
	transferRepo := new(tests.MockTransferRepo)
	userRepo := new(tests.MockUserRepo)
	svc := service.NewTransferService(transferRepo, userRepo)

	fromID := uuid.New().String()
	toID := uuid.New().String()
	amount := 100.0

	transferRepo.On("CreateTransfer", ctx, mock.Anything).Return(errors.New("db error")).Once()

	id, err := svc.CreateTransfer(ctx, fromID, toID, amount)
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)

	transferRepo.AssertExpectations(t)
}

func TestTransferService_TransferService_CreateTransfer_Success(t *testing.T) {
	ctx := context.Background()
	transferRepo := new(tests.MockTransferRepo)
	userRepo := new(tests.MockUserRepo)
	svc := service.NewTransferService(transferRepo, userRepo)

	fromID := uuid.New().String()
	toID := uuid.New().String()
	amount := 100.0

	transferRepo.On("CreateTransfer", ctx, mock.Anything).Return(nil).Once()

	select {
	case <-queue.JobsChan:
	default:
	}

	id, err := svc.CreateTransfer(ctx, fromID, toID, amount)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)

	select {
	case job := <-queue.JobsChan:
		assert.Equal(t, uuid.MustParse(fromID), job.SenderId)
		assert.Equal(t, uuid.MustParse(toID), job.ReceiverId)
		assert.Equal(t, amount, job.Amount)
		assert.Equal(t, id, job.TransactionId)
	case <-time.After(time.Second):
		t.Fatal("expected job in queue but got none")
	}

	transferRepo.AssertExpectations(t)
}
