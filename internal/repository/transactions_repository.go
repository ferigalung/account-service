package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/ferigalung/account-service/internal/model/balances"
	trxModel "github.com/ferigalung/account-service/internal/model/transactions"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db}
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, doc trxModel.Transaction) (*float64, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := fmt.Sprintf("INSERT INTO %s (account_id, type, amount) VALUES ($1, $2, $3)", doc.TableName())
	if _, err := tx.Exec(ctx, query, doc.AccountID, doc.Type, doc.Amount); err != nil {
		return nil, err
	}

	// define operator according to transaction type
	opt := "-"
	if doc.Type == "deposit" {
		opt = "+"
	}

	var balances balances.Balance
	query = fmt.Sprintf("UPDATE %s SET balance = balance %s $1 WHERE account_id = $2 RETURNING (balance)", balances.TableName(), opt)
	if err := tx.QueryRow(ctx, query, doc.Amount, doc.AccountID).Scan(&balances.Balance); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23514" {
			return nil, errors.New("Insufficient balance")
		}
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &balances.Balance, nil
}
