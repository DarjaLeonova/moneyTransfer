package repository

import (
	"context"
	"moneyTransfer/internal/domain/model"
)

type TransferRepository interface {
	GetTransactionsByUserId(ctx context.Context, userId string) ([]model.Transaction, error)
	CreateTransfer(ctx context.Context, tx model.Transaction) error
	UpdateTransactionStatus(ctx context.Context, txID, status string) error
}

type UserRepository interface {
	GetBalance(ctx context.Context, userId string) (float64, error)
	GetById(ctx context.Context, userId string) (model.User, error)
	UpdateBalance(ctx context.Context, userID string, newBalance float64) error
}
