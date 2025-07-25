package service_tests

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"moneyTransfer/internal/domain/model"
	"moneyTransfer/internal/domain/service"
	"testing"
)

func TestGetBalance_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepo)
	svc := service.NewUserService(repo)

	repo.On("GetBalance", ctx, "user123").Return(100.0, nil)

	balance, err := svc.GetBalance(ctx, "user123")
	assert.NoError(t, err)
	assert.Equal(t, 100.0, balance)

	repo.AssertExpectations(t)
}

func TestGetBalance_Error(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepo)
	svc := service.NewUserService(repo)

	repo.On("GetBalance", ctx, "user123").Return(0.0, errors.New("db error"))

	balance, err := svc.GetBalance(ctx, "user123")
	assert.Error(t, err)
	assert.Equal(t, 0.0, balance)

	repo.AssertExpectations(t)
}

func TestGetById_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepo)
	svc := service.NewUserService(repo)

	expectedUser := model.User{
		Id:        uuid.MustParse("7141b92f-a8c8-471e-83e5-7fc72da61cb9"),
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Balance:   200.0,
	}

	repo.On("GetById", ctx, "7141b92f-a8c8-471e-83e5-7fc72da61cb9").Return(expectedUser, nil)

	user, err := svc.GetById(ctx, "7141b92f-a8c8-471e-83e5-7fc72da61cb9")
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	repo.AssertExpectations(t)
}

func TestGetById_Error(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepo)
	svc := service.NewUserService(repo)

	repo.On("GetById", ctx, "user123").Return(model.User{}, errors.New("db error"))

	user, err := svc.GetById(ctx, "user123")
	assert.Error(t, err)
	assert.Equal(t, model.User{}, user)

	repo.AssertExpectations(t)
}
