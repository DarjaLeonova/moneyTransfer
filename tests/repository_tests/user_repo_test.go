package repository_tests

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"moneyTransfer/internal/domain/model"
	"moneyTransfer/internal/repository"
	"moneyTransfer/tests"
	"testing"
)

func TestUserRepo_GetBalance_Success(t *testing.T) {
	db, mock := tests.SetupMockDB(t)

	repo := repository.NewUserRepository(db)

	mock.ExpectQuery(`SELECT balance FROM users WHERE id = \$1`).
		WithArgs("user123").
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(250.0))

	balance, err := repo.GetBalance(context.Background(), "user123")
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

	mock.ExpectQuery(`SELECT id, first_name, last_name, email, balance FROM users WHERE id = \$1`).
		WithArgs("a5bceab4-9dab-4d7a-8cd5-4ba832ebf899").
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "balance"}).
			AddRow("a5bceab4-9dab-4d7a-8cd5-4ba832ebf899", "John", "Doe", "john@example.com", 200.0))

	user, err := repo.GetById(context.Background(), "a5bceab4-9dab-4d7a-8cd5-4ba832ebf899")
	require.NoError(t, err)
	require.Equal(t, "a5bceab4-9dab-4d7a-8cd5-4ba832ebf899", user.Id.String())
	require.Equal(t, "John", user.FirstName)
	require.Equal(t, "Doe", user.LastName)
	require.Equal(t, "john@example.com", user.Email)
	require.Equal(t, 200.0, user.Balance)

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

	mock.ExpectExec(`UPDATE users SET balance = \$1 WHERE id = \$2`).
		WithArgs(300.0, "user1").
		WillReturnResult(sqlmock.NewResult(1, 1)) // 1 row affected

	err := repo.UpdateBalance(context.Background(), "user1", 300.0)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_UpdateBalance_Error(t *testing.T) {
	db, mock := tests.SetupMockDB(t)

	repo := repository.NewUserRepository(db)

	mock.ExpectExec(`UPDATE users SET balance = \$1 WHERE id = \$2`).
		WithArgs(300.0, "user1").
		WillReturnError(errors.New("update failed"))

	err := repo.UpdateBalance(context.Background(), "user1", 300.0)
	require.Error(t, err)
	require.EqualError(t, err, "update failed")

	require.NoError(t, mock.ExpectationsWereMet())
}
