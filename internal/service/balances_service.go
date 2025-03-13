package service

import (
	"github.com/ferigalung/account-service/internal/repository"
)

type BalanceService struct {
	repo *repository.AccountRepository
}

func NewBalanceService(repo *repository.AccountRepository) *BalanceService {
	return &BalanceService{repo: repo}
}
