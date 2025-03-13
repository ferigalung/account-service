package service

import (
	"github.com/ferigalung/account-service/internal/repository"
)

type TransactionService struct {
	repo *repository.AccountRepository
}

func NewTransactionService(repo *repository.AccountRepository) *TransactionService {
	return &TransactionService{repo: repo}
}
