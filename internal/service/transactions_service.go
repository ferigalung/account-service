package service

import (
	"context"
	"errors"
	"math"

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

func (s *TransactionService) GetTrxHistory(ctx context.Context, params trxModel.ListPayload) (trxModel.Pagination, error) {
	account, err := s.accountRepo.GetAccountByAccountNumber(ctx, params.AccountNumber)
	if err != nil {
		if pgxscan.NotFound(err) {
			msg := "Invalid account number"
			logger.Log("info", msg, fiber.Map{"accountNumber": params.AccountNumber})
			return trxModel.Pagination{}, errors.New(msg)
		}
		logger.Log("error", "Failed to get account", fiber.Map{"error": err.Error()})
		return trxModel.Pagination{}, err
	}
	params.AccountID = account.ID

	result, err := s.repo.GetTransactions(ctx, params)
	if err != nil {
		return trxModel.Pagination{}, err
	}

	totalPage := 0
	totalDataOnPage := len(result.Data)
	if result.TotalRecords < int64(params.Size) {
		totalPage = 1
	} else {
		totalPage = int(math.Ceil(float64(result.TotalRecords) / float64(params.Size)))
	}

	return trxModel.Pagination{
		Data: result.Data,
		Meta: trxModel.PaginationMeta{
			Page:            params.Page,
			TotalData:       result.TotalRecords,
			TotalPage:       totalPage,
			TotalDataOnPage: totalDataOnPage,
		},
	}, nil
}
