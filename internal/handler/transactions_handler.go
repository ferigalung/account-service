package handler

import (
	"github.com/ferigalung/account-service/internal/model/transactions"
	"github.com/ferigalung/account-service/internal/service"
	cv "github.com/ferigalung/account-service/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	service   *service.TransactionService
	validator *cv.ValidatorImpl
}

func NewTransactionHandler(svc *service.TransactionService, cv *cv.ValidatorImpl) *TransactionHandler {
	return &TransactionHandler{service: svc, validator: cv}
}

func (h *TransactionHandler) Deposit(c *fiber.Ctx) error {
	var payload transactions.Payload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"data":    nil,
			"message": "Invalid request body",
		})
	}
	payload.Type = "deposit"

	if err := h.validator.Validate(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"data":    nil,
			"message": err.Error(),
		})
	}

	detailBalance, err := h.service.CreateTransaction(c.Context(), payload)
	if err != nil {
		if err.Error() == "Invalid account number" {
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
		"data":    detailBalance,
		"message": "Deposit successful",
	})
}

func (h *TransactionHandler) Withdraw(c *fiber.Ctx) error {
	var payload transactions.Payload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"data":    nil,
			"message": "Invalid request body",
		})
	}
	payload.Type = "withdraw"

	if err := h.validator.Validate(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"data":    nil,
			"message": err.Error(),
		})
	}

	detailBalance, err := h.service.CreateTransaction(c.Context(), payload)
	if err != nil {
		if err.Error() == "Invalid account number" || err.Error() == "Insufficient balance" {
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
		"data":    detailBalance,
		"message": "Withdrawal successful",
	})
}
