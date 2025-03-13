package handler

import (
	"github.com/ferigalung/account-service/internal/model/accounts"
	"github.com/ferigalung/account-service/internal/service"
	cv "github.com/ferigalung/account-service/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type AccountHandler struct {
	service   *service.AccountService
	validator *cv.ValidatorImpl
}

func NewAccountHandler(svc *service.AccountService, cv *cv.ValidatorImpl) *AccountHandler {
	return &AccountHandler{service: svc, validator: cv}
}

func (h *AccountHandler) CreateAccount(c *fiber.Ctx) error {
	var payload accounts.Register

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"data":    nil,
			"message": "Invalid request body",
		})
	}

	if err := h.validator.Validate(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"data":    nil,
			"message": err.Error(),
		})
	}

	account, err := h.service.CreateAccount(c.Context(), payload)
	if err != nil {
		if err.Error() == "NIK or phone number already registered" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"data":    nil,
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"data":    nil,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"data":    account,
		"message": "Successfully create new account",
	})
}
