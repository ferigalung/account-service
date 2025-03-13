package repository

import (
	"context"
	"fmt"

	"github.com/ferigalung/account-service/internal/model/accounts"
	"github.com/ferigalung/account-service/internal/model/balances"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BalanceRepository struct {
	db *pgxpool.Pool
}

func NewBalanceRepository(db *pgxpool.Pool) *BalanceRepository {
	return &BalanceRepository{db}
}

func (r *BalanceRepository) GetDetailBalance(ctx context.Context, accountNumber string) (*balances.DetailBalance, error) {
	var detailBalance balances.DetailBalance
	var balance balances.Balance
	var account accounts.Account
	query := fmt.Sprintf("SELECT b.balance, a.account_number, a.name FROM %s b JOIN %s a ON b.account_id = a.id WHERE a.account_number = $1", balance.TableName(), account.TableName())
	if err := pgxscan.Get(ctx, r.db, &detailBalance, query, accountNumber); err != nil {
		return nil, err
	}

	return &detailBalance, nil
}
