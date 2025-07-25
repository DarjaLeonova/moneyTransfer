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
}

func NewTransferService(transferRepo contracts.TransferRepository, userRepo contracts.UserRepository) TransferService {
	return &transferService{transferRepo: transferRepo, userRepo: userRepo}
}

func (t *transferService) GetTransactionsByUserId(ctx context.Context, userId string) ([]model.Transaction, error) {
	transactions, err := t.transferRepo.GetTransactionsByUserId(ctx, userId)
	if err != nil {
		logger.Log.Error("failed to get transactions by user id", "userId", userId, "error", err)
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	logger.Log.Info("Transactions retrieved", "transactions", transactions)

	return transactions, nil
}

func (t *transferService) CreateTransfer(ctx context.Context, from, to string, amount float64) (uuid.UUID, error) {
	if amount <= 0 {
		return uuid.Nil, fmt.Errorf("amount must be greater than zero")
	}

	tx := model.Transaction{
		Id:         uuid.New(),
		SenderId:   uuid.MustParse(from),
		ReceiverId: uuid.MustParse(to),
		Amount:     amount,
		Status:     "PENDING",
		CreatedAt:  time.Now(),
	}

	err := t.transferRepo.CreateTransfer(ctx, tx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create transfer: %w", err)
	}

	job := queue.TransferJob{
		SenderId:      tx.SenderId,
		ReceiverId:    tx.ReceiverId,
		Amount:        tx.Amount,
		TransactionId: tx.Id,
	}

	queue.Enqueue(job)
	return tx.Id, nil
}
