package repository

import (
	"context"
	"database/sql"
	"moneyTransfer/internal/domain/contracts"
	"moneyTransfer/internal/domain/model"
)

type UserRepo struct {
	db *sql.DB
}

var _ contracts.UserRepository = (*UserRepo)(nil)

func NewUserRepository(db *sql.DB) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) GetBalance(ctx context.Context, userId string) (float64, error) {
	query := `SELECT balance FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, userId)

	var balance float64

	err := row.Scan(&balance)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func (r *UserRepo) GetById(ctx context.Context, userId string) (model.User, error) {
	query := `SELECT id, first_name, last_name, email, balance FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, userId)

	var user model.User
	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Balance)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepo) UpdateBalance(ctx context.Context, userId string, newBalance float64) error {
	query := `UPDATE users SET balance = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, newBalance, userId)

	return err
}
