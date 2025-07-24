package dtos

import "moneyTransfer/internal/domain/model"

type TransactionResponseDto struct {
	Transactions []model.Transaction `json:"transactions"`
}
