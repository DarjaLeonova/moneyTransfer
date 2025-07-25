package repository

import (
	"context"
	"database/sql"
	"moneyTransfer/internal/domain/contracts"
	"moneyTransfer/internal/domain/model"
)

type TransferRepo struct {
	db *sql.DB
}

var _ contracts.TransferRepository = (*TransferRepo)(nil)

func NewTransferRepository(db *sql.DB) *TransferRepo {
	return &TransferRepo{db}
}

func (r *TransferRepo) GetTransactionsByUserId(ctx context.Context, userId string) ([]model.Transaction, error) {
	query := `SELECT * FROM transactions WHERE sender_id = $1 OR receiver_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []model.Transaction

	for rows.Next() {
		var transaction model.Transaction
		if err := rows.Scan(
			&transaction.Id,
			&transaction.SenderId,
			&transaction.ReceiverId,
			&transaction.Amount,
			&transaction.Status,
			&transaction.CreatedAt,
		); err != nil {
			return transfers, err
		}
		transfers = append(transfers, transaction)
	}
	if err := rows.Err(); err != nil {
		return transfers, err
	}

	return transfers, nil
}

func (r *TransferRepo) CreateTransfer(ctx context.Context, tx model.Transaction) error {
	query := `INSERT INTO transactions (id, sender_id, receiver_id, amount, status, created_at)
              VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query,
		tx.Id, tx.SenderId, tx.ReceiverId, tx.Amount, tx.Status, tx.CreatedAt)
	return err
}

func (r *TransferRepo) UpdateTransactionStatus(ctx context.Context, txID, status string) error {
	query := `UPDATE transactions SET status = $1 WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, status, txID)
	return err
}
