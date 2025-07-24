package dtos

import "github.com/google/uuid"

type TransactionRequestDto struct {
	From   uuid.UUID `json:"from"`
	To     uuid.UUID `json:"to"`
	Amount float64   `json:"amount"`
}
