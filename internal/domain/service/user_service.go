package service

import (
	"context"
	"moneyTransfer/internal/domain/contracts"
	"moneyTransfer/internal/domain/model"
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
	return u.userRepo.GetBalance(ctx, userId)
}

func (u *userService) GetById(ctx context.Context, userId string) (model.User, error) {
	return u.userRepo.GetById(ctx, userId)
}
