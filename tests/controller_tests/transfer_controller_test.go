package controller_tests

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"moneyTransfer/api/handler"
	"moneyTransfer/internal/domain/dtos"
	"moneyTransfer/internal/domain/model"
	"moneyTransfer/tests"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTransferController_GetTransactionsByUserId_Success(t *testing.T) {
	svc, logger, controller := initTransferController()

	userId := "7141b92f-a8c8-471e-83e5-7fc72da61cb9"
	var expectedTransactions = dtos.TransactionResponseDto{
		Transactions: []model.Transaction{
			{
				Id:         uuid.MustParse("a5bceab4-9dab-4d7a-8cd5-4ba832ebf899"),
				SenderId:   uuid.MustParse(userId),
				ReceiverId: uuid.MustParse("9d02adbc-27ca-4695-9d92-10cb35db67f4"),
				Amount:     100,
				Status:     model.StatusSuccess,
				CreatedAt:  time.Time{},
			},
		},
	}

	svc.On("GetTransactionsByUserId", mock.Anything, userId).Return(expectedTransactions.Transactions, nil)
	logger.On("Info", "transactions fetched successfully", "response", expectedTransactions).Return()

	req := httptest.NewRequest(http.MethodGet, "/transfers/"+userId, nil)
	rr := httptest.NewRecorder()

	req = mux.SetURLVars(req, map[string]string{"userId": userId})

	controller.GetTransactionsByUserId(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp dtos.TransactionResponseDto
	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)

	assert.Equal(t, expectedTransactions.Transactions, resp.Transactions)
	svc.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestTransferController_GetTransactionsByUserId_Error(t *testing.T) {
	svc, _, controller := initTransferController()

	expectedErr := errors.New("database connection failed")
	userId := "thgh"

	svc.On("GetTransactionsByUserId", mock.Anything, userId).Return([]model.Transaction{}, expectedErr)

	req := httptest.NewRequest(http.MethodGet, "/transfers/"+userId, nil)
	req = mux.SetURLVars(req, map[string]string{"userId": userId})
	rr := httptest.NewRecorder()

	controller.GetTransactionsByUserId(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var errResp dtos.ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&errResp)
	require.NoError(t, err)
	assert.Equal(t, "Error fetching transactions", errResp.Message)
}

func TestTransferController_GetTransactionsByUserId_MissingUserId(t *testing.T) {
	_, _, controller := initTransferController()

	req := httptest.NewRequest(http.MethodGet, "/transfers/", nil)
	rr := httptest.NewRecorder()

	controller.GetTransactionsByUserId(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var errResp dtos.ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&errResp)
	require.NoError(t, err)
	assert.Equal(t, "User Id is required", errResp.Message)
}

func initTransferController() (*tests.MockTransferService, *tests.MockLogger, *handler.TransferController) {
	transferSvc := new(tests.MockTransferService)
	logger := new(tests.MockLogger)
	controller := handler.NewTransferController(transferSvc, logger)
	return transferSvc, logger, controller
}
