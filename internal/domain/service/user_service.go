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
	log      logger.Logger
}

func NewUserService(userRepo contracts.UserRepository, logger logger.Logger) UserService {
	return &userService{userRepo: userRepo, log: logger}
}

func (u *userService) GetBalance(ctx context.Context, userId string) (float64, error) {
	balance, err := u.userRepo.GetBalance(ctx, userId)
	if err != nil {
		u.log.Error("failed to get balance", "userId", userId, "error", err)
		return 0.0, fmt.Errorf("failed to get balance: %w", err)
	}

	u.log.Info("balance retrieved", "userId", userId, "balance", balance)
	return balance, nil
}

func (u *userService) GetById(ctx context.Context, userId string) (model.User, error) {
	user, err := u.userRepo.GetById(ctx, userId)
	if err != nil {
		u.log.Error("failed to get user", "userId", userId, "error", err)
		return model.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	u.log.Info("user retrieved", "userId", user.Id)
	return user, nil
}
