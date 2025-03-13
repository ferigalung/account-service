package handler

import (
	"github.com/ferigalung/account-service/internal/service"
	cv "github.com/ferigalung/account-service/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type BalanceHandler struct {
	service   *service.BalanceService
	validator *cv.ValidatorImpl
}

func NewBalanceHandler(svc *service.BalanceService, cv *cv.ValidatorImpl) *BalanceHandler {
	return &BalanceHandler{service: svc, validator: cv}
}

func (h *BalanceHandler) GetBalance(c *fiber.Ctx) error {
	accountNumber := c.Params("an")

	if accountNumber == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"data":    nil,
			"message": "Account number is required in param",
		})
	}

	detailBalance, err := h.service.GetBalance(c.Context(), accountNumber)
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
		"message": "Successfully get balance",
	})
}
