package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"moneyTransfer/internal/domain/dtos"
	"moneyTransfer/internal/domain/service"
	"moneyTransfer/pkg/logger"
	"net/http"
)

type UserController struct {
	UserService service.UserService
}

func NewUserController(service service.UserService) *UserController {
	return &UserController{UserService: service}
}

// @Summary Get user balance
// @Description Get current balance for a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} dtos.BalanceResponseDto
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /balance/{userId} [get]
func (c *UserController) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	if userID == "" {
		logger.Log.Error("UserId is required")
		dtos.WriteErrorResponse(w, "User ID is required", "GetUserBalance: User ID is required", http.StatusBadRequest)
		return
	}

	balance, err := c.UserService.GetBalance(r.Context(), userID)
	if err != nil {
		logger.Log.Error("Error fetching user balance", "err", err)
		dtos.WriteErrorResponse(w, "Error fetching user balance", err.Error(), http.StatusInternalServerError)
		return
	}

	b := dtos.BalanceResponseDto{Balance: balance}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(b)

	logger.Log.Info("Balance fetched successfully", "balance", balance)
}
