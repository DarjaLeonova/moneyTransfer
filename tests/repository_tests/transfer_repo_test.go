package repository_tests

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"moneyTransfer/internal/domain/model"
	"moneyTransfer/internal/repository"
	"moneyTransfer/tests"
	"testing"
	"time"
)

func TestTransferRepo_GetTransactionsByUserId_Success(t *testing.T) {
	db, mock := tests.SetupMockDB(t)
	repo := repository.NewTransferRepository(db)

	senderId := "c775d967-7b54-463f-9923-90f219d8224d"
	txId := "a5bceab4-9dab-4d7a-8cd5-4ba832ebf899"
	receiverId := "ed9c2b61-3908-413b-b355-a6c36d1a0cb3"

	mock.ExpectQuery(`SELECT id, sender_id, receiver_id, amount, status, created_at FROM transactions WHERE sender_id = \$1 OR receiver_id = \$1`).
		WithArgs(senderId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "sender_id", "receiver_id", "amount", "status", "created_at"}).
			AddRow(txId, senderId, receiverId, 100.0, model.StatusSuccess, time.Time{}))

	expectedTransactions := []model.Transaction{
		{
			Id:         uuid.MustParse(txId),
			SenderId:   uuid.MustParse(senderId),
			ReceiverId: uuid.MustParse(receiverId),
			Amount:     100,
			Status:     model.StatusSuccess,
			CreatedAt:  time.Time{},
		},
	}

	transactions, err := repo.GetTransactionsByUserId(context.Background(), senderId)
	require.NoError(t, err)
	require.Equal(t, expectedTransactions, transactions)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTransferRepo_GetTransactionsByUserId_Error(t *testing.T) {
	db, mock := tests.SetupMockDB(t)
	repo := repository.NewTransferRepository(db)

	mock.ExpectQuery(`SELECT id, sender_id, receiver_id, amount, status, created_at FROM transactions WHERE sender_id = \$1 OR receiver_id = \$1`).
		WithArgs("missing_user").
		WillReturnError(sql.ErrNoRows)

	transactions, err := repo.GetTransactionsByUserId(context.Background(), "missing_user")
	require.Error(t, err)
	require.Len(t, transactions, 0)
	require.ErrorIs(t, err, sql.ErrNoRows)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTransferRepo_CreateTransfer_Success(t *testing.T) {
	db, mock := tests.SetupMockDB(t)
	repo := repository.NewTransferRepository(db)

	tx := model.Transaction{
		Id:         uuid.MustParse("7141b92f-a8c8-471e-83e5-7fc72da61cb9"),
		SenderId:   uuid.MustParse("d489b057-aa2e-4d34-9020-d2b42294dc42"),
		ReceiverId: uuid.MustParse("ed9c2b61-3908-413b-b355-a6c36d1a0cb3"),
		Amount:     50.0,
		Status:     model.StatusPending,
		CreatedAt:  time.Now(),
	}

	mock.ExpectExec(`INSERT INTO transactions \(id, sender_id, receiver_id, amount, status, created_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\)`).
		WithArgs(tx.Id, tx.SenderId, tx.ReceiverId, tx.Amount, tx.Status, tx.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.CreateTransfer(context.Background(), tx)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTransferRepo_CreateTransfer_Error(t *testing.T) {
	db, mock := tests.SetupMockDB(t)
	repo := repository.NewTransferRepository(db)

	tx := model.Transaction{
		Id:         uuid.MustParse("7141b92f-a8c8-471e-83e5-7fc72da61cb9"),
		SenderId:   uuid.MustParse("d489b057-aa2e-4d34-9020-d2b42294dc42"),
		ReceiverId: uuid.MustParse("ed9c2b61-3908-413b-b355-a6c36d1a0cb3"),
		Amount:     50.0,
		Status:     model.StatusPending,
		CreatedAt:  time.Now(),
	}

	mock.ExpectExec(`INSERT INTO transactions \(id, sender_id, receiver_id, amount, status, created_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\)`).
		WithArgs(tx.Id, tx.SenderId, tx.ReceiverId, tx.Amount, tx.Status, tx.CreatedAt).
		WillReturnError(errors.New("db error"))

	err := repo.CreateTransfer(context.Background(), tx)
	require.Error(t, err)
	require.EqualError(t, err, "db error")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTransferRepo_UpdateTransactionStatus_Success(t *testing.T) {
	db, mock := tests.SetupMockDB(t)
	repo := repository.NewTransferRepository(db)

	txId := "7141b92f-a8c8-471e-83e5-7fc72da61cb9"
	newStatus := model.StatusSuccess

	mock.ExpectExec(`UPDATE transactions SET status = \$1 WHERE id = \$2`).
		WithArgs(newStatus, txId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.UpdateTransactionStatus(context.Background(), txId, newStatus)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTransferRepo_UpdateTransactionStatus_Error(t *testing.T) {
	db, mock := tests.SetupMockDB(t)
	repo := repository.NewTransferRepository(db)

	txId := "7141b92f-a8c8-471e-83e5-7fc72da61cb9"
	newStatus := model.StatusSuccess

	mock.ExpectExec(`UPDATE transactions SET status = \$1 WHERE id = \$2`).
		WithArgs(newStatus, txId).
		WillReturnError(errors.New("update failed"))

	err := repo.UpdateTransactionStatus(context.Background(), txId, newStatus)
	require.Error(t, err)
	require.EqualError(t, err, "update failed")

	require.NoError(t, mock.ExpectationsWereMet())
}
