package service

import (
	"context"
	"fmt"
	"moneyTransfer/internal/domain/contracts"
	"moneyTransfer/internal/domain/model"
	"moneyTransfer/pkg/logger"
)

type UserService interface {
	GetBalance(ctx context.Context, userId string) (float64, error)
	GetById(ctx context.Context, userId string) (model.User, error)
}

type userService struct {
	userRepo contracts.UserRepository
}

func NewUserService(userRepo contracts.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (u *userService) GetBalance(ctx context.Context, userId string) (float64, error) {
	balance, err := u.userRepo.GetBalance(ctx, userId)
	if err != nil {
		logger.Log.Error("failed to get balance", "err", err)
		return 0.0, fmt.Errorf("failed to get balance: %w", err)
	}

	logger.Log.Info("Balance retrieved", "balance", balance)
	return balance, nil
}

func (u *userService) GetById(ctx context.Context, userId string) (model.User, error) {
	user, err := u.userRepo.GetById(ctx, userId)
	if err != nil {
		logger.Log.Error("failed to get user", "err", err)
		return model.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	logger.Log.Info("User retrieved", "user", user)
	return user, nil
}
