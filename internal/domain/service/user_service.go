package service

import (
	"context"
	"moneyTransfer/internal/domain/model"
	"moneyTransfer/internal/repository"
)

type UserService interface {
	GetBalance(ctx context.Context, userId string) (float64, error)
	GetById(ctx context.Context, userId string) (model.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (u *userService) GetBalance(ctx context.Context, userId string) (float64, error) {
	return u.userRepo.GetBalance(ctx, userId)
}

func (u *userService) GetById(ctx context.Context, userId string) (model.User, error) {
	return u.userRepo.GetById(ctx, userId)
}
