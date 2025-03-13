package repository

import "github.com/jackc/pgx/v5/pgxpool"

type BalanceRepository struct {
	db *pgxpool.Pool
}

func NewBalanceRepository(db *pgxpool.Pool) *BalanceRepository {
	return &BalanceRepository{db}
}
