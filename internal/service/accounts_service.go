package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	accountModel "github.com/ferigalung/account-service/internal/model/accounts"
	"github.com/ferigalung/account-service/internal/repository"
	"github.com/ferigalung/account-service/pkg/logger"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gofiber/fiber/v2"
)

type AccountService struct {
	repo *repository.AccountRepository
}

func NewAccountService(repo *repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func generateAccountNumber() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bankCode := "1023"

	// generate 6 digit unique number
	accountNumber := r.Intn(900000) + 100000

	// calculate checksum (mod 10)
	rawNumber := fmt.Sprintf("%s%d", bankCode, accountNumber)
	sum := 0
	for _, digit := range rawNumber {
		sum += int(digit - '0')
	}
	checksum := sum % 10

	return fmt.Sprintf("%s%d%d", bankCode, accountNumber, checksum)
}

func (s *AccountService) CreateAccount(ctx context.Context, payload accountModel.Register) (*accountModel.Account, error) {
	params := accountModel.UniqueAccount{NIK: payload.NIK, Phone: payload.Phone}
	account, err := s.repo.GetUniqueAccount(ctx, params)
	fmt.Println(pgxscan.NotFound(err))
	if err != nil && !pgxscan.NotFound(err) {
		logger.Log("error", "Failed to get account", fiber.Map{"error": err.Error()})
		return nil, err
	}

	if account != nil {
		msg := "NIK or phone number already registered"
		logger.Log("info", msg, fiber.Map{"nik": payload.NIK, "phone": payload.Phone})
		return nil, errors.New(msg)
	}

	newAccount, err := s.repo.InsertAccount(ctx, accountModel.Account{
		AccountNumber: generateAccountNumber(),
		Name:          payload.Name,
		NIK:           payload.NIK,
		Phone:         payload.Phone,
	})
	if err != nil {
		return nil, err
	}

	return newAccount, nil
}
