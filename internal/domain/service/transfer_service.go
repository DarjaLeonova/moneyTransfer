package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"moneyTransfer/internal/domain/contracts"
	"moneyTransfer/internal/domain/model"
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
	return t.transferRepo.GetTransactionsByUserId(ctx, userId)
}

func (t *transferService) CreateTransfer(ctx context.Context, from, to string, amount float64) (uuid.UUID, error) {
	if amount <= 0 {
		return uuid.Nil, fmt.Errorf("amount must be greater than zero")
	}

	senderBalance, err := t.userRepo.GetBalance(ctx, from)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get sender balance: %w", err)
	}

	if senderBalance < amount {
		return uuid.Nil, fmt.Errorf("insufficient funds: balance=%.2f, required=%.2f", senderBalance, amount)
	}

	receiverBalance, err := t.userRepo.GetBalance(ctx, to)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get receiver balance: %w", err)
	}

	senderBalance -= amount
	receiverBalance += amount

	err = t.userRepo.UpdateBalance(ctx, from, senderBalance)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to update sender balance: %w", err)
	}

	err = t.userRepo.UpdateBalance(ctx, to, receiverBalance)
	if err != nil {
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
		return uuid.Nil, fmt.Errorf("failed to create transfer: %w", err)
	}

	return tx.Id, nil
}
