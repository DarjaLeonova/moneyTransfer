package controller_tests

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"moneyTransfer/api/handler"
	"moneyTransfer/internal/domain/dtos"
	"moneyTransfer/tests"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserController_GetUserBalance_Success(t *testing.T) {
	svc, logger, controller := initUserCOntroller()

	userId := "7141b92f-a8c8-471e-83e5-7fc72da61cb9"
	expectedBalance := 100.50
	dtoBalance := dtos.BalanceResponseDto{Balance: expectedBalance}

	svc.On("GetBalance", mock.Anything, userId).Return(expectedBalance, nil)
	logger.On("Info", "balance fetched successfully", "response", dtoBalance).Return()

	req := httptest.NewRequest(http.MethodGet, "/balance/"+userId, nil)
	rr := httptest.NewRecorder()

	req = mux.SetURLVars(req, map[string]string{"userId": userId})

	controller.GetUserBalance(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp dtos.BalanceResponseDto
	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)

	assert.Equal(t, expectedBalance, resp.Balance)
	svc.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestUserController_GetUserBalance_Error(t *testing.T) {
	svc, _, controller := initUserCOntroller()

	expectedErr := errors.New("database connection failed")
	userId := "thgh"

	svc.On("GetBalance", mock.Anything, userId).Return(0.0, expectedErr)

	req := httptest.NewRequest(http.MethodGet, "/balance/"+userId, nil)
	req = mux.SetURLVars(req, map[string]string{"userId": userId})
	rr := httptest.NewRecorder()

	controller.GetUserBalance(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var errResp dtos.ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&errResp)
	require.NoError(t, err)
	assert.Equal(t, "Error fetching balance", errResp.Message)
}

func TestUserController_GetUserBalance_MissingUserId(t *testing.T) {
	_, _, controller := initUserCOntroller()

	req := httptest.NewRequest(http.MethodGet, "/balance/", nil)
	rr := httptest.NewRecorder()

	controller.GetUserBalance(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var errResp dtos.ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&errResp)
	require.NoError(t, err)
	assert.Equal(t, "User Id is required", errResp.Message)
}

func initUserCOntroller() (*tests.MockUserService, *tests.MockLogger, *handler.UserController) {
	svc := new(tests.MockUserService)
	logger := new(tests.MockLogger)
	controller := handler.NewUserController(svc, logger)
	return svc, logger, controller
}
