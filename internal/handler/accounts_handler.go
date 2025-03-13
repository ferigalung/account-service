package handler

import (
	"github.com/ferigalung/account-service/internal/model"
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
	var payload model.Register

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	if err := h.validator.Validate(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	account, err := h.service.CreateAccount(c.Context(), payload)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Account not found"})
	}

	return c.JSON(account)
}
