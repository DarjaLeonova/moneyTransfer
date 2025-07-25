package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"moneyTransfer/internal/domain/contracts"
	"moneyTransfer/internal/domain/model"
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
	logger.Log.Info("starting transfer", "from", from, "to", to)
	if amount <= 0 {
		logger.Log.Error("amount must be greater than zero", "amount", amount)
		return uuid.Nil, fmt.Errorf("amount must be greater than zero")
	}

	senderBalance, err := t.userRepo.GetBalance(ctx, from)
	if err != nil {
		logger.Log.Error("failed to get sender balance", "err", err)
		return uuid.Nil, fmt.Errorf("failed to get sender balance: %w", err)
	}

	if senderBalance < amount {
		logger.Log.Error("insufficient funds", "balance", senderBalance, "amount", amount, "err", err)
		return uuid.Nil, fmt.Errorf("insufficient funds: balance=%.2f, required=%.2f", senderBalance, amount)
	}

	receiverBalance, err := t.userRepo.GetBalance(ctx, to)
	if err != nil {
		logger.Log.Error("failed to get receiver balance", "err", err)
		return uuid.Nil, fmt.Errorf("failed to get receiver balance: %w", err)
	}

	senderBalance -= amount
	receiverBalance += amount

	err = t.userRepo.UpdateBalance(ctx, from, senderBalance)
	if err != nil {
		logger.Log.Error("failed to update sender balance", "err", err)
		return uuid.Nil, fmt.Errorf("failed to update sender balance: %w", err)
	}

	err = t.userRepo.UpdateBalance(ctx, to, receiverBalance)
	if err != nil {
		logger.Log.Error("failed to update receiver balance", "err", err)
		return uuid.Nil, fmt.Errorf("failed to update receiver balance: %w", err)
	}

	tx := model.Transaction{
		Id:         uuid.New(),
		SenderId:   uuid.MustParse(from),
		ReceiverId: uuid.MustParse(to),
		Amount:     amount,
		Status:     "SUCCESS",
		CreatedAt:  time.Now(),
	}

	err = t.transferRepo.CreateTransfer(ctx, tx)
	if err != nil {
		logger.Log.Error("failed to create transfer", "err", err)
		return uuid.Nil, fmt.Errorf("failed to create transfer: %w", err)
	}

	logger.Log.Info("Created transaction", "id", tx.Id)
	return tx.Id, nil
}
