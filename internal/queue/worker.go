package queue

import (
	"context"
	"moneyTransfer/internal/domain/contracts"
	"moneyTransfer/internal/domain/model"
	"moneyTransfer/pkg/logger"
)

func ProcessJob(ctx context.Context, job TransferJob, userRepo contracts.UserRepository, transferRepo contracts.TransferRepository) error {
	if job.Amount <= 0 {
		logger.Log.Error("amount must be greater than zero", "amount", job.Amount)
		return transferRepo.UpdateTransactionStatus(ctx, job.TransactionId.String(), model.StatusFailed)
	}

	senderBalance, err := userRepo.GetBalance(ctx, job.SenderId.String())
	if err != nil {
		logger.Log.Error("failed to get sender balance", "err", err)
		return transferRepo.UpdateTransactionStatus(ctx, job.TransactionId.String(), model.StatusFailed)
	}

	if senderBalance < job.Amount {
		logger.Log.Error("insufficient funds", "balance", senderBalance, "amount", job.Amount)
		return transferRepo.UpdateTransactionStatus(ctx, job.TransactionId.String(), model.StatusFailed)
	}

	receiverBalance, err := userRepo.GetBalance(ctx, job.ReceiverId.String())
	if err != nil {
		logger.Log.Error("failed to get receiver balance", "err", err)
		return transferRepo.UpdateTransactionStatus(ctx, job.TransactionId.String(), model.StatusFailed)
	}

	err = userRepo.UpdateBalance(ctx, job.SenderId.String(), senderBalance-job.Amount)
	if err != nil {
		logger.Log.Error("failed to update sender balance", "err", err)
		return transferRepo.UpdateTransactionStatus(ctx, job.TransactionId.String(), model.StatusFailed)
	}

	err = userRepo.UpdateBalance(ctx, job.ReceiverId.String(), receiverBalance+job.Amount)
	if err != nil {
		logger.Log.Error("failed to update receiver balance", "err", err)
		return transferRepo.UpdateTransactionStatus(ctx, job.TransactionId.String(), model.StatusFailed)
	}

	err = transferRepo.UpdateTransactionStatus(ctx, job.TransactionId.String(), model.StatusSuccess)
	if err != nil {
		logger.Log.Error("failed to update transaction status", "err", err)
		return err
	}

	logger.Log.Info("transfer completed", "transaction_id", job.TransactionId)
	return nil
}

func StartWorker(userRepo contracts.UserRepository, transferRepo contracts.TransferRepository) {
	go func() {
		for job := range JobsChan {
			ctx := context.Background()
			err := ProcessJob(ctx, job, userRepo, transferRepo)
			if err != nil {
				logger.Log.Error("failed to process job", "err", err)
			}
		}
	}()
}
