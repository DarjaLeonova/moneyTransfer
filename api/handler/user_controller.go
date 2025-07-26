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
	log         logger.Logger
}

func NewUserController(service service.UserService, logger logger.Logger) *UserController {
	return &UserController{UserService: service, log: logger}
}

// @Summary Get user balance
// @Description Get current balance for a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param userId path string true "User Id"
// @Success 200 {object} dtos.BalanceResponseDto
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /balance/{userId} [get]
func (c *UserController) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["userId"]
	if userId == "" {
		dtos.WriteErrorResponse(w, "User Id is required", "GetUserBalance: User Id is required", http.StatusBadRequest)
		return
	}

	balance, err := c.UserService.GetBalance(r.Context(), userId)
	if err != nil {
		dtos.WriteErrorResponse(w, "Error fetching balance", err.Error(), http.StatusInternalServerError)
		return
	}

	response := dtos.BalanceResponseDto{Balance: balance}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	c.log.Info("balance fetched successfully", "response", response)
}
