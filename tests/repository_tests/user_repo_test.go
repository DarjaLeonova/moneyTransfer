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
)

func TestUserRepo_GetBalance_Success(t *testing.T) {
	db, mock := tests.SetupMockDB(t)

	repo := repository.NewUserRepository(db)
	userId := "ed9c2b61-3908-413b-b355-a6c36d1a0cb3"

	mock.ExpectQuery(`SELECT balance FROM users WHERE id = \$1`).
		WithArgs(userId).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(250.0))

	balance, err := repo.GetBalance(context.Background(), userId)
	require.NoError(t, err)
	require.Equal(t, 250.0, balance)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_GetBalance_Error(t *testing.T) {
	db, mock := tests.SetupMockDB(t)

	repo := repository.NewUserRepository(db)
	mock.ExpectQuery(`SELECT balance FROM users WHERE id = \$1`).
		WithArgs("nonexistent-user").
		WillReturnError(sql.ErrNoRows)

	balance, err := repo.GetBalance(context.Background(), "nonexistent-user")

	require.Error(t, err)
	require.Equal(t, 0.0, balance)
	require.ErrorIs(t, err, sql.ErrNoRows)
	require.NoError(t, mock.ExpectationsWereMet())

	require.Error(t, err)
}

func TestUserRepo_GetById_Success(t *testing.T) {
	db, mock := tests.SetupMockDB(t)

	repo := repository.NewUserRepository(db)
	expectedUser := model.User{
		Id:        uuid.MustParse("7141b92f-a8c8-471e-83e5-7fc72da61cb9"),
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Balance:   200.0,
	}

	mock.ExpectQuery(`SELECT id, first_name, last_name, email, balance FROM users WHERE id = \$1`).
		WithArgs(expectedUser.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "balance"}).
			AddRow(expectedUser.Id, expectedUser.FirstName, expectedUser.LastName, expectedUser.Email, expectedUser.Balance))

	user, err := repo.GetById(context.Background(), expectedUser.Id.String())
	require.NoError(t, err)
	require.Equal(t, expectedUser, user)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_GetById_Error(t *testing.T) {
	db, mock := tests.SetupMockDB(t)

	repo := repository.NewUserRepository(db)

	mock.ExpectQuery(`SELECT id, first_name, last_name, email, balance FROM users WHERE id = \$1`).
		WithArgs("missing_user").
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetById(context.Background(), "missing_user")
	require.Error(t, err)
	require.Equal(t, model.User{}, user)
	require.ErrorIs(t, err, sql.ErrNoRows)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_UpdateBalance_Success(t *testing.T) {
	db, mock := tests.SetupMockDB(t)

	repo := repository.NewUserRepository(db)
	userId := "ed9c2b61-3908-413b-b355-a6c36d1a0cb3"

	mock.ExpectExec(`UPDATE users SET balance = \$1 WHERE id = \$2`).
		WithArgs(300.0, userId).
		WillReturnResult(sqlmock.NewResult(1, 1)) // 1 row affected

	err := repo.UpdateBalance(context.Background(), userId, 300.0)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_UpdateBalance_Error(t *testing.T) {
	db, mock := tests.SetupMockDB(t)

	repo := repository.NewUserRepository(db)
	userId := "ed9c2b61-3908-413b-b355-a6c36d1a0cb3"

	mock.ExpectExec(`UPDATE users SET balance = \$1 WHERE id = \$2`).
		WithArgs(300.0, userId).
		WillReturnError(errors.New("update failed"))

	err := repo.UpdateBalance(context.Background(), userId, 300.0)
	require.Error(t, err)
	require.EqualError(t, err, "update failed")

	require.NoError(t, mock.ExpectationsWereMet())
}
