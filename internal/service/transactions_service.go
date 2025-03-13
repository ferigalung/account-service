package service

import (
	"context"
	"errors"

	balanceModel "github.com/ferigalung/account-service/internal/model/balances"
	trxModel "github.com/ferigalung/account-service/internal/model/transactions"
	"github.com/ferigalung/account-service/internal/repository"
	"github.com/ferigalung/account-service/pkg/logger"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gofiber/fiber/v2"
)

type TransactionService struct {
	repo        *repository.TransactionRepository
	accountRepo *repository.AccountRepository
	balanceRepo *repository.BalanceRepository
}

func NewTransactionService(
	repo *repository.TransactionRepository,
	accountRepo *repository.AccountRepository,
	balanceRepo *repository.BalanceRepository,
) *TransactionService {
	return &TransactionService{repo: repo, accountRepo: accountRepo, balanceRepo: balanceRepo}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, payload trxModel.Payload) (*balanceModel.DetailBalance, error) {
	account, err := s.accountRepo.GetAccountByAccountNumber(ctx, payload.AccountNumber)
	if err != nil {
		if pgxscan.NotFound(err) {
			msg := "Invalid account number"
			logger.Log("info", msg, fiber.Map{"accountNumber": payload.AccountNumber})
			return nil, errors.New(msg)
		}
		logger.Log("error", "Failed to get account", fiber.Map{"error": err.Error()})
		return nil, err
	}

	currentBalance, err := s.repo.CreateTransaction(ctx, trxModel.Transaction{
		AccountID: account.ID,
		Type:      payload.Type,
		Amount:    payload.Amount,
	})
	if err != nil {
		logger.Log("error", "Failed to create transaction", fiber.Map{"error": err.Error()})
		return nil, err
	}

	return &balanceModel.DetailBalance{
		AccountNumber: account.AccountNumber,
		Name:          account.Name,
		Balance:       *currentBalance,
	}, nil
}
