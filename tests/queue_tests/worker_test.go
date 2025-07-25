package queue_tests

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"moneyTransfer/internal/domain/model"
	"moneyTransfer/internal/queue"
	"moneyTransfer/tests"
	"testing"
)

func TestProcessJob_InvalidAmount(t *testing.T) {
	ctx := context.Background()
	userRepo := new(tests.MockUserRepo)
	transferRepo := new(tests.MockTransferRepo)

	job := queue.TransferJob{
		Amount:        0,
		SenderId:      uuid.MustParse("d489b057-aa2e-4d34-9020-d2b42294dc42"),
		ReceiverId:    uuid.MustParse("ed9c2b61-3908-413b-b355-a6c36d1a0cb3"),
		TransactionId: uuid.MustParse("f5c184f5-38f1-46d0-b9c4-47da6ad55552"),
	}

	transferRepo.On("UpdateTransactionStatus", ctx, job.TransactionId.String(), model.StatusFailed).Return(nil)

	err := queue.ProcessJob(ctx, job, userRepo, transferRepo)

	require.NoError(t, err)
	transferRepo.AssertCalled(t, "UpdateTransactionStatus", ctx, job.TransactionId.String(), model.StatusFailed)
}

func TestProcessJob_FailedToGetSenderBalance(t *testing.T) {
	ctx := context.Background()
	userRepo := new(tests.MockUserRepo)
	transferRepo := new(tests.MockTransferRepo)

	job := queue.TransferJob{
		Amount:        80,
		SenderId:      uuid.MustParse("eeb552ab-bea8-4183-8f62-9e4fe9281759"),
		ReceiverId:    uuid.MustParse("ed9c2b61-3908-413b-b355-a6c36d1a0cb3"),
		TransactionId: uuid.MustParse("f5c184f5-38f1-46d0-b9c4-47da6ad55552"),
	}

	userRepo.On("GetBalance", ctx, job.SenderId.String()).Return(0.0, errors.New("db error"))
	transferRepo.On("UpdateTransactionStatus", ctx, job.TransactionId.String(), model.StatusFailed).Return(nil)

	err := queue.ProcessJob(ctx, job, userRepo, transferRepo)

	require.NoError(t, err)
	userRepo.AssertExpectations(t)
	transferRepo.AssertExpectations(t)
}

func TestProcessJob_FailedToGetReceiverBalance(t *testing.T) {
	ctx := context.Background()
	userRepo := new(tests.MockUserRepo)
	transferRepo := new(tests.MockTransferRepo)

	job := queue.TransferJob{
		Amount:        80,
		SenderId:      uuid.MustParse("ed9c2b61-3908-413b-b355-a6c36d1a0cb3"),
		ReceiverId:    uuid.MustParse("eeb552ab-bea8-4183-8f62-9e4fe9281759"),
		TransactionId: uuid.MustParse("f5c184f5-38f1-46d0-b9c4-47da6ad55552"),
	}

	userRepo.On("GetBalance", ctx, job.SenderId.String()).Return(100.0, nil)
	userRepo.On("GetBalance", ctx, job.ReceiverId.String()).Return(0.0, errors.New("db error"))
	transferRepo.On("UpdateTransactionStatus", ctx, job.TransactionId.String(), model.StatusFailed).Return(nil)

	err := queue.ProcessJob(ctx, job, userRepo, transferRepo)

	require.NoError(t, err)
	userRepo.AssertExpectations(t)
	transferRepo.AssertExpectations(t)
}

func TestProcessJob_InsufficientFunds(t *testing.T) {
	ctx := context.Background()
	userRepo := new(tests.MockUserRepo)
	transferRepo := new(tests.MockTransferRepo)

	job := queue.TransferJob{
		Amount:        100,
		SenderId:      uuid.MustParse("d489b057-aa2e-4d34-9020-d2b42294dc42"),
		ReceiverId:    uuid.MustParse("ed9c2b61-3908-413b-b355-a6c36d1a0cb3"),
		TransactionId: uuid.MustParse("f5c184f5-38f1-46d0-b9c4-47da6ad55552"),
	}

	userRepo.On("GetBalance", ctx, job.SenderId.String()).Return(50.0, nil)
	transferRepo.On("UpdateTransactionStatus", ctx, job.TransactionId.String(), model.StatusFailed).Return(nil)

	err := queue.ProcessJob(ctx, job, userRepo, transferRepo)

	require.NoError(t, err)
	userRepo.AssertCalled(t, "GetBalance", ctx, job.SenderId.String())
	transferRepo.AssertCalled(t, "UpdateTransactionStatus", ctx, job.TransactionId.String(), model.StatusFailed)
}

func TestProcessJob_FailedToUpdateSenderBalance(t *testing.T) {
	ctx := context.Background()
	userRepo := new(tests.MockUserRepo)
	transferRepo := new(tests.MockTransferRepo)

	job := queue.TransferJob{
		Amount:        80,
		SenderId:      uuid.MustParse("ed9c2b61-3908-413b-b355-a6c36d1a0cb3"),
		ReceiverId:    uuid.MustParse("eeb552ab-bea8-4183-8f62-9e4fe9281759"),
		TransactionId: uuid.MustParse("f5c184f5-38f1-46d0-b9c4-47da6ad55552"),
	}

	userRepo.On("GetBalance", ctx, job.SenderId.String()).Return(100.0, nil)
	userRepo.On("GetBalance", ctx, job.ReceiverId.String()).Return(50.0, nil)

	userRepo.On("UpdateBalance", ctx, job.SenderId.String(), 20.0).Return(errors.New("update failed"))

	transferRepo.On("UpdateTransactionStatus", ctx, job.TransactionId.String(), model.StatusFailed).Return(nil)

	err := queue.ProcessJob(ctx, job, userRepo, transferRepo)

	require.NoError(t, err)
	userRepo.AssertExpectations(t)
	transferRepo.AssertExpectations(t)
}

func TestProcessJob_FailedToUpdateReceiverBalance(t *testing.T) {
	ctx := context.Background()
	userRepo := new(tests.MockUserRepo)
	transferRepo := new(tests.MockTransferRepo)

	job := queue.TransferJob{
		Amount:        80,
		SenderId:      uuid.MustParse("ed9c2b61-3908-413b-b355-a6c36d1a0cb3"),
		ReceiverId:    uuid.MustParse("eeb552ab-bea8-4183-8f62-9e4fe9281759"),
		TransactionId: uuid.MustParse("f5c184f5-38f1-46d0-b9c4-47da6ad55552"),
	}

	userRepo.On("GetBalance", ctx, job.SenderId.String()).Return(100.0, nil)
	userRepo.On("GetBalance", ctx, job.ReceiverId.String()).Return(50.0, nil)

	userRepo.On("UpdateBalance", ctx, job.SenderId.String(), 20.0).Return(nil)
	userRepo.On("UpdateBalance", ctx, job.ReceiverId.String(), 130.0).Return(errors.New("update failed"))

	transferRepo.On("UpdateTransactionStatus", ctx, job.TransactionId.String(), model.StatusFailed).Return(nil)

	err := queue.ProcessJob(ctx, job, userRepo, transferRepo)

	require.NoError(t, err)
	userRepo.AssertExpectations(t)
	transferRepo.AssertExpectations(t)
}

func TestProcessJob_FailedToUpdateTransactionStatusSuccess(t *testing.T) {
	ctx := context.Background()
	userRepo := new(tests.MockUserRepo)
	transferRepo := new(tests.MockTransferRepo)

	job := queue.TransferJob{
		Amount:        80,
		SenderId:      uuid.MustParse("ed9c2b61-3908-413b-b355-a6c36d1a0cb3"),
		ReceiverId:    uuid.MustParse("eeb552ab-bea8-4183-8f62-9e4fe9281759"),
		TransactionId: uuid.MustParse("f5c184f5-38f1-46d0-b9c4-47da6ad55552"),
	}

	userRepo.On("GetBalance", ctx, job.SenderId.String()).Return(100.0, nil)
	userRepo.On("GetBalance", ctx, job.ReceiverId.String()).Return(50.0, nil)

	userRepo.On("UpdateBalance", ctx, job.SenderId.String(), 20.0).Return(nil)
	userRepo.On("UpdateBalance", ctx, job.ReceiverId.String(), 130.0).Return(nil)

	transferRepo.On("UpdateTransactionStatus", ctx, job.TransactionId.String(), model.StatusSuccess).
		Return(errors.New("update status failed"))

	err := queue.ProcessJob(ctx, job, userRepo, transferRepo)

	require.Error(t, err)
	userRepo.AssertExpectations(t)
	transferRepo.AssertExpectations(t)
}
