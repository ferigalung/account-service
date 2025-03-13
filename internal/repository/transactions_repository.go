package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/ferigalung/account-service/internal/model/balances"
	trxModel "github.com/ferigalung/account-service/internal/model/transactions"
	"github.com/georgysavva/scany/v2/pgxscan"
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

func (r *TransactionRepository) GetTransactions(ctx context.Context, params trxModel.ListPayload) (trxModel.ListRepo, error) {
	var trx trxModel.Transaction
	var trxList []trxModel.Transaction
	var totalRecords int64

	query := fmt.Sprintf("SELECT * FROM %s WHERE account_id = $1", trx.TableName())
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE account_id = $1", trx.TableName())
	args := []interface{}{}
	args = append(args, params.AccountID)
	argIdx := 2

	if params.StartDate != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argIdx)
		countQuery += fmt.Sprintf(" AND created_at >= $%d", argIdx)
		args = append(args, *params.StartDate)
		argIdx++
	}
	if params.EndDate != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argIdx)
		countQuery += fmt.Sprintf(" AND created_at <= $%d", argIdx)
		args = append(args, *params.EndDate)
		argIdx++
	}
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, params.Size)
	args = append(args, (params.Page-1)*params.Size)

	if err := pgxscan.Select(ctx, r.db, &trxList, query, args...); err != nil {
		return trxModel.ListRepo{}, err
	}
	if err := pgxscan.Get(ctx, r.db, &totalRecords, countQuery, args[:len(args)-2]...); err != nil {
		return trxModel.ListRepo{}, err
	}

	return trxModel.ListRepo{
		Data:         trxList,
		TotalRecords: totalRecords,
	}, nil

}
