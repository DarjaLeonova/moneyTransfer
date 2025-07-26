package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"moneyTransfer/internal/domain/contracts"
	"moneyTransfer/internal/domain/model"
	"moneyTransfer/internal/queue"
	"moneyTransfer/pkg/logger"
	"time"
)

type TransferService interface {
	CreateTransfer(ctx context.Context, from, to string, amount float64) (uuid.UUID, error)
	GetTransactionsByUserId(ctx context.Context, userId string) ([]model.Transaction, error)
}

type transferService struct {
	transferRepo contracts.TransferRepository
	userRepo     contracts.UserRepository
	log          logger.Logger
}

func NewTransferService(transferRepo contracts.TransferRepository, userRepo contracts.UserRepository, logger logger.Logger) TransferService {
	return &transferService{transferRepo: transferRepo, userRepo: userRepo, log: logger}
}

func (t *transferService) GetTransactionsByUserId(ctx context.Context, userId string) ([]model.Transaction, error) {
	transactions, err := t.transferRepo.GetTransactionsByUserId(ctx, userId)
	if err != nil {
		t.log.Error("failed to get transactions by user id", "userId", userId, "error", err)
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	t.log.Info("transactions retrieved", "transactions", transactions)
	return transactions, nil
}

func (t *transferService) CreateTransfer(ctx context.Context, from, to string, amount float64) (uuid.UUID, error) {
	if amount <= 0 {
		t.log.Warn("invalid transfer amount", "amount", amount, "from", from, "to", to)
		return uuid.Nil, fmt.Errorf("amount must be greater than zero")
	}

	tx := model.Transaction{
		Id:         uuid.New(),
		SenderId:   uuid.MustParse(from),
		ReceiverId: uuid.MustParse(to),
		Amount:     amount,
		Status:     model.StatusPending,
		CreatedAt:  time.Now(),
	}

	err := t.transferRepo.CreateTransfer(ctx, tx)
	if err != nil {
		t.log.Error("failed to create transfer", "tx", tx, "error", err)
		return uuid.Nil, fmt.Errorf("failed to create transfer: %w", err)
	}

	t.log.Info("transfer created", "tx", tx)

	job := queue.TransferJob{
		SenderId:      tx.SenderId,
		ReceiverId:    tx.ReceiverId,
		Amount:        tx.Amount,
		TransactionId: tx.Id,
	}

	t.log.Info("enqueuing transfer job", "job", job)

	queue.Enqueue(job)
	return tx.Id, nil
}
