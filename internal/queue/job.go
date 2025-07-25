package queue

import "github.com/google/uuid"

type TransferJob struct {
	SenderId      uuid.UUID
	ReceiverId    uuid.UUID
	Amount        float64
	TransactionId uuid.UUID
}
