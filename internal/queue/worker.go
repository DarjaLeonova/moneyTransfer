package queue

import (
	"context"
	"moneyTransfer/internal/domain/contracts"
	"moneyTransfer/pkg/logger"
)

func StartWorker(userRepo contracts.UserRepository, transferRepo contracts.TransferRepository) {
	go func() {
		for job := range JobsChan {
			ctx := context.Background()

			if job.Amount <= 0 {
				logger.Log.Error("amount must be greater than zero", "amount", job.Amount)
				_ = transferRepo.UpdateTransactionStatus(ctx, job.TransactionId.String(), "FAILED")
				continue
			}

			senderBalance, err := userRepo.GetBalance(ctx, job.SenderId.String())
			if err != nil {
				logger.Log.Error("failed to get sender balance", "err", err)
				_ = transferRepo.UpdateTransactionStatus(ctx, job.TransactionId.String(), "FAILED")
				continue
			}

			if senderBalance < job.Amount {
				logger.Log.Error("insufficient funds", "balance", senderBalance, "amount", job.Amount)
				_ = transferRepo.UpdateTransactionStatus(ctx, job.TransactionId.String(), "FAILED")
				continue
			}

			receiverBalance, err := userRepo.GetBalance(ctx, job.ReceiverId.String())
			if err != nil {
				logger.Log.Error("failed to get receiver balance", "err", err)
				_ = transferRepo.UpdateTransactionStatus(ctx, job.TransactionId.String(), "FAILED")
				continue
			}

			err = userRepo.UpdateBalance(ctx, job.SenderId.String(), senderBalance-job.Amount)
			if err != nil {
				logger.Log.Error("failed to update sender balance", "err", err)
				_ = transferRepo.UpdateTransactionStatus(ctx, job.TransactionId.String(), "FAILED")
				continue
			}

			err = userRepo.UpdateBalance(ctx, job.ReceiverId.String(), receiverBalance+job.Amount)
			if err != nil {
				logger.Log.Error("failed to update receiver balance", "err", err)
				_ = transferRepo.UpdateTransactionStatus(ctx, job.TransactionId.String(), "FAILED")
				continue
			}

			err = transferRepo.UpdateTransactionStatus(ctx, job.TransactionId.String(), "SUCCESS")
			if err != nil {
				logger.Log.Error("failed to update transaction status", "err", err)
				continue
			}

			logger.Log.Info("transfer completed", "transaction_id", job.TransactionId)
		}
	}()
}
