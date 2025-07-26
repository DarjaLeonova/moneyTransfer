package service_tests

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"moneyTransfer/internal/domain/model"
	"moneyTransfer/internal/domain/service"
	"moneyTransfer/internal/queue"
	"moneyTransfer/tests"
	"testing"
	"time"
)

func TestTransferService_GetTransactionsByUserId_Success(t *testing.T) {
	ctx, transferRepo, svc, logger := inittransferService()

	userId := "7141b92f-a8c8-471e-83e5-7fc72da61cb9"

	expectedTransactions := []model.Transaction{
		{
			Id:         uuid.MustParse("a5bceab4-9dab-4d7a-8cd5-4ba832ebf899"),
			SenderId:   uuid.MustParse(userId),
			ReceiverId: uuid.MustParse("9d02adbc-27ca-4695-9d92-10cb35db67f4"),
			Amount:     100,
			Status:     model.StatusSuccess,
			CreatedAt:  time.Time{},
		},
	}

	transferRepo.On("GetTransactionsByUserId", ctx, userId).Return(expectedTransactions, nil)
	logger.On("Info", "transactions retrieved", "transactions", expectedTransactions).Return()

	transactions, err := svc.GetTransactionsByUserId(ctx, userId)
	require.NoError(t, err)
	assert.Equal(t, expectedTransactions, transactions)

	transferRepo.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestTransferService_GetTransactionsByUserId_Error(t *testing.T) {
	ctx, transferRepo, svc, logger := inittransferService()

	userId := "7141b92f-a8c8-471e-83e5-7fc72da61cb9"

	transferRepo.On("GetTransactionsByUserId", ctx, userId).
		Return(make([]model.Transaction, 0), errors.New("db error"))

	logger.On("Error", "failed to get transactions by user id", "userId", userId, "error", mock.Anything).Return()

	transactions, err := svc.GetTransactionsByUserId(ctx, userId)
	assert.Error(t, err)
	assert.Len(t, transactions, 0)

	transferRepo.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestTransferService_CreateTransfer_AmountLessOrEqualZero(t *testing.T) {
	ctx, _, svc, logger := inittransferService()

	fromId := "7141b92f-a8c8-471e-83e5-7fc72da61cb9"
	toId := "ed9c2b61-3908-413b-b355-a6c36d1a0cb3"

	logger.On("Warn", "invalid transfer amount", "amount", 0.0, "from", fromId, "to", toId).Return()

	id, err := svc.CreateTransfer(ctx, fromId, toId, 0)
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
	logger.AssertExpectations(t)
}

func TestTransferService_TransferService_CreateTransfer_RepoError(t *testing.T) {
	ctx, transferRepo, svc, logger := inittransferService()

	fromId := "7141b92f-a8c8-471e-83e5-7fc72da61cb9"
	toId := "ed9c2b61-3908-413b-b355-a6c36d1a0cb3"
	amount := 100.0

	transferRepo.On("CreateTransfer", ctx, mock.Anything).Return(errors.New("db error")).Once()
	logger.On("Error", "failed to create transfer", "tx", mock.Anything, "error", mock.Anything).Return()

	id, err := svc.CreateTransfer(ctx, fromId, toId, amount)
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)

	transferRepo.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestTransferService_TransferService_CreateTransfer_Success(t *testing.T) {
	ctx, transferRepo, svc, logger := inittransferService()

	fromId := "7141b92f-a8c8-471e-83e5-7fc72da61cb9"
	toId := "ed9c2b61-3908-413b-b355-a6c36d1a0cb3"
	amount := 100.0

	transferRepo.On("CreateTransfer", ctx, mock.Anything).Return(nil).Once()
	logger.On("Info", "transfer created", "tx", mock.Anything).Return()
	logger.On("Info", "enqueuing transfer job", "job", mock.Anything).Return()

	select {
	case <-queue.JobsChan:
	default:
	}

	id, err := svc.CreateTransfer(ctx, fromId, toId, amount)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, id)

	select {
	case job := <-queue.JobsChan:
		assert.Equal(t, uuid.MustParse(fromId), job.SenderId)
		assert.Equal(t, uuid.MustParse(toId), job.ReceiverId)
		assert.Equal(t, amount, job.Amount)
		assert.Equal(t, id, job.TransactionId)
	case <-time.After(time.Second):
		t.Fatal("expected job in queue but got none")
	}

	transferRepo.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func inittransferService() (context.Context, *tests.MockTransferRepo, service.TransferService, *tests.MockLogger) {
	ctx := context.Background()
	transferRepo := new(tests.MockTransferRepo)
	userRepo := new(tests.MockUserRepo)
	logger := new(tests.MockLogger)
	svc := service.NewTransferService(transferRepo, userRepo, logger)
	return ctx, transferRepo, svc, logger
}
