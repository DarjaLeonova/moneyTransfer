package service_tests

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"moneyTransfer/internal/domain/model"
	"moneyTransfer/internal/domain/service"
	"moneyTransfer/tests"
	"testing"
)

func TestUserService_GetBalance_Success(t *testing.T) {
	ctx, repo, svc, logger := initUserService()

	repo.On("GetBalance", ctx, "user123").Return(100.0, nil)
	logger.On("Info", "balance retrieved", "userId", "user123", "balance", 100.0).Return()

	balance, err := svc.GetBalance(ctx, "user123")
	assert.NoError(t, err)
	assert.Equal(t, 100.0, balance)

	repo.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestUserService_GetBalance_Error(t *testing.T) {
	ctx, repo, svc, logger := initUserService()

	repo.On("GetBalance", ctx, "user123").Return(0.0, errors.New("db error"))
	logger.On("Error", "failed to get balance", "userId", "user123", "error", mock.Anything).Return()

	balance, err := svc.GetBalance(ctx, "user123")
	assert.Error(t, err)
	assert.Equal(t, 0.0, balance)

	repo.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestUserService_GetById_Success(t *testing.T) {
	ctx, repo, svc, logger := initUserService()

	expectedUser := model.User{
		Id:        uuid.MustParse("7141b92f-a8c8-471e-83e5-7fc72da61cb9"),
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Balance:   200.0,
	}

	repo.On("GetById", ctx, "7141b92f-a8c8-471e-83e5-7fc72da61cb9").Return(expectedUser, nil)
	logger.On("Info", "user retrieved", "userId", expectedUser.Id).Return()

	user, err := svc.GetById(ctx, "7141b92f-a8c8-471e-83e5-7fc72da61cb9")
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	repo.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestUserService_GetById_Error(t *testing.T) {
	ctx, repo, svc, logger := initUserService()

	repo.On("GetById", ctx, "user123").Return(model.User{}, errors.New("db error"))
	logger.On("Error", "failed to get user", "userId", "user123", "error", mock.Anything).Return()

	user, err := svc.GetById(ctx, "user123")
	assert.Error(t, err)
	assert.Equal(t, model.User{}, user)

	repo.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func initUserService() (context.Context, *tests.MockUserRepo, service.UserService, *tests.MockLogger) {
	ctx := context.Background()
	repo := new(tests.MockUserRepo)
	logger := new(tests.MockLogger)
	svc := service.NewUserService(repo, logger)
	return ctx, repo, svc, logger
}
