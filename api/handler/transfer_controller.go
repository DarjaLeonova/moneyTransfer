package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"moneyTransfer/internal/domain/dtos"
	"moneyTransfer/internal/domain/service"
	"moneyTransfer/pkg/logger"
	"net/http"
)

type TransferController struct {
	TransferService service.TransferService
	log             logger.Logger
}

func NewTransferController(transferService service.TransferService, logger logger.Logger) *TransferController {
	return &TransferController{TransferService: transferService, log: logger}
}

// @Summary Get transactions by user ID
// @Description Get all transactions for a specific user
// @Tags transfers
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {array} dtos.TransactionResponseDto
// @Failure 400 {object} dtos.ErrorResponse
// @Router /transfers/{userId} [get]
func (c *TransferController) GetTransactionsByUserID(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userId"]
	if userID == "" {
		dtos.WriteErrorResponse(w, "User ID is required", "GetUserBalance: User ID is required", http.StatusBadRequest)
		return
	}

	transactions, err := c.TransferService.GetTransactionsByUserId(r.Context(), userID)
	if err != nil {
		dtos.WriteErrorResponse(w, "Error fetching transactions", err.Error(), http.StatusInternalServerError)
		return
	}

	response := dtos.TransactionResponseDto{Transactions: transactions}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	c.log.Info("transactions fetched successfully", "response", response)
}

// @Summary Create new transaction
// @Description Create a new money transfer transaction
// @Tags transfers
// @Accept json
// @Produce json
// @Param transaction body dtos.TransactionRequestDto true "Transaction details"
// @Success 200 {object} dtos.CreateTransactionResponseDto
// @Failure 400 {object} dtos.ErrorResponse
// @Router /transfers [post]
func (c *TransferController) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transactionRequestDto dtos.TransactionRequestDto
	err := json.NewDecoder(r.Body).Decode(&transactionRequestDto)
	if err != nil {
		dtos.WriteErrorResponse(w, "Error parsing request body", err.Error(), http.StatusBadRequest)
		return
	}

	id, err := c.TransferService.CreateTransfer(r.Context(), transactionRequestDto.From.String(), transactionRequestDto.To.String(), transactionRequestDto.Amount)
	if err != nil {
		dtos.WriteErrorResponse(w, "Failed to create transfer", err.Error(), http.StatusInternalServerError)
		return
	}

	response := dtos.CreateTransactionResponseDto{
		TransactionID: id,
		Status:        "completed",
		Message:       "Transaction was successful",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	c.log.Info("transaction was successful", "id", id)
}
