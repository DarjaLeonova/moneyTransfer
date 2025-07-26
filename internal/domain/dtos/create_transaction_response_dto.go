package dtos

import "github.com/google/uuid"

type CreateTransactionResponseDto struct {
	TransactionId uuid.UUID `json:"transaction_id"`
	Status        string    `json:"status"`
	Message       string    `json:"message"`
}
