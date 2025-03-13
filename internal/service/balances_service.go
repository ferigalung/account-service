package service

import (
	"context"
	"errors"

	"github.com/ferigalung/account-service/internal/model/balances"
	"github.com/ferigalung/account-service/internal/repository"
	"github.com/ferigalung/account-service/pkg/logger"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gofiber/fiber/v2"
)

type BalanceService struct {
	repo *repository.BalanceRepository
}

func NewBalanceService(repo *repository.BalanceRepository) *BalanceService {
	return &BalanceService{repo: repo}
}

func (s *BalanceService) GetBalance(ctx context.Context, accountNumber string) (*balances.DetailBalance, error) {
	detailBalance, err := s.repo.GetDetailBalance(ctx, accountNumber)
	if err != nil {
		if pgxscan.NotFound(err) {
			msg := "Invalid account number"
			logger.Log("info", msg, fiber.Map{"accountNumber": accountNumber})
			return nil, errors.New(msg)
		}
		logger.Log("error", "Failed to get balance", fiber.Map{"error": err.Error()})
		return nil, err
	}

	return detailBalance, nil
}
