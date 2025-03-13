package repository

import (
	"context"
	"fmt"

	accountModel "github.com/ferigalung/account-service/internal/model/accounts"
	balanceModel "github.com/ferigalung/account-service/internal/model/balances"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{db}
}

func (r *AccountRepository) GetUniqueAccount(ctx context.Context, params accountModel.UniqueAccount) (*accountModel.Account, error) {
	var account accountModel.Account
	query := fmt.Sprintf("SELECT * FROM %s WHERE nik = $1 OR phone = $2", account.TableName())
	if err := pgxscan.Get(ctx, r.db, &account, query, params.NIK, params.Phone); err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *AccountRepository) GetAccountByAccountNumber(ctx context.Context, accountNumber string) (*accountModel.Account, error) {
	var account accountModel.Account
	query := fmt.Sprintf("SELECT * FROM %s WHERE account_number $1", account.TableName())
	if err := pgxscan.Select(ctx, r.db, &account, query, accountNumber); err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *AccountRepository) InsertAccount(ctx context.Context, doc accountModel.Account) (*accountModel.Account, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var newAccount accountModel.Account
	query := fmt.Sprintf("INSERT INTO %s (account_number, name, nik, phone) VALUES ($1, $2, $3, $4) RETURNING *", doc.TableName())
	if err := tx.QueryRow(ctx, query, doc.AccountNumber, doc.Name, doc.NIK, doc.Phone).Scan(
		&newAccount.ID, &newAccount.AccountNumber, &newAccount.Name, &newAccount.NIK, &newAccount.Phone, &newAccount.CreatedAt, &newAccount.UpdatedAt); err != nil {
		return nil, err
	}

	var balance balanceModel.Balance
	if _, err := tx.Exec(ctx, fmt.Sprintf("INSERT INTO %s (account_id, balance) VALUES ($1, $2)", balance.TableName()), newAccount.ID, 0); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &newAccount, nil
}
