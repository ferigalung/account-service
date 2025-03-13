package service

import (
	"context"

	"github.com/ferigalung/account-service/internal/model"
	"github.com/ferigalung/account-service/internal/repository"
)

type AccountService struct {
	repo *repository.AccountRepository
}

func NewAccountService(repo *repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) CreateAccount(ctx context.Context, payload model.Register) (*model.Account, error) {
	return nil, nil
}
