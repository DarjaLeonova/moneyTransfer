package queue

import (
	"context"
	"moneyTransfer/internal/domain/contracts"
	"moneyTransfer/internal/domain/model"
	"moneyTransfer/pkg/logger"
	"time"
)

func ProcessJob(ctx context.Context, job TransferJob, userRepo contracts.UserRepository, transferRepo contracts.TransferRepository, log logger.Logger) error {
	if job.Amount <= 0 {
		log.Error("amount must be greater than zero", "amount", job.Amount)
		return transferRepo.UpdateTransactionStatus(ctx, job.TransactionId.String(), model.StatusFailed)
	}

	senderBalance, err := userRepo.GetBalance(ctx, job.SenderId.String())
	if err != nil {
		log.Error("failed to get sender balance", "error", err)
		return transferRepo.UpdateTransactionStatus(ctx, job.TransactionId.String(), model.StatusFailed)
	}

	if senderBalance < job.Amount {
		log.Error("insufficient funds", "balance", senderBalance, "amount", job.Amount)
		return transferRepo.UpdateTransactionStatus(ctx, job.TransactionId.String(), model.StatusFailed)
	}

	receiverBalance, err := userRepo.GetBalance(ctx, job.ReceiverId.String())
	if err != nil {
		log.Error("failed to get receiver balance", "error", err)
		return transferRepo.UpdateTransactionStatus(ctx, job.TransactionId.String(), model.StatusFailed)
	}

	err = userRepo.UpdateBalance(ctx, job.SenderId.String(), senderBalance-job.Amount)
	if err != nil {
		log.Error("failed to update sender balance", "error", err)
		return transferRepo.UpdateTransactionStatus(ctx, job.TransactionId.String(), model.StatusFailed)
	}

	err = userRepo.UpdateBalance(ctx, job.ReceiverId.String(), receiverBalance+job.Amount)
	if err != nil {
		log.Error("failed to update receiver balance", "error", err)
		return transferRepo.UpdateTransactionStatus(ctx, job.TransactionId.String(), model.StatusFailed)
	}

	err = transferRepo.UpdateTransactionStatus(ctx, job.TransactionId.String(), model.StatusSuccess)
	if err != nil {
		log.Error("failed to update transaction status", "error", err)
		return err
	}

	log.Info("transfer completed", "transaction_id", job.TransactionId)
	return nil
}

func StartWorker(userRepo contracts.UserRepository, transferRepo contracts.TransferRepository, log logger.Logger) {
	go func() {
		for job := range JobsChan {
			ctx := context.Background()
			err := ProcessJob(ctx, job, userRepo, transferRepo, log)
			if err != nil {
				log.Error("failed to process job", "error", err)
			}
			time.Sleep(500 * time.Millisecond) //emulate long processing
		}
	}()
}
