package model

import (
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	Id         uuid.UUID `json:"id"`
	SenderId   uuid.UUID `json:"sender_id"`
	ReceiverId uuid.UUID `json:"receiver_id"`
	Amount     float64   `json:"amount"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}
