package repository

import "github.com/jackc/pgx/v5/pgxpool"

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db}
}
