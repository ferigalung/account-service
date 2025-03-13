package handler

import (
	"time"

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

func (h *TransactionHandler) GetTrxHistory(c *fiber.Ctx) error {
	var startDate *time.Time
	rawStartDate := c.Query("startDate")
	if rawStartDate != "" {
		s, err := time.Parse(time.DateOnly, rawStartDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"data":    nil,
				"message": "Invalid start date",
			})
		}
		startDate = &s

	}
	var endDate *time.Time
	rawEndDate := c.Query("endDate")
	if rawEndDate != "" {
		e, err := time.Parse(time.DateOnly, rawEndDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"data":    nil,
				"message": "Invalid end date",
			})
		}
		endDate = &e
	}

	payload := transactions.ListPayload{
		Page:          c.QueryInt("page"),
		Size:          c.QueryInt("size"),
		StartDate:     startDate,
		EndDate:       endDate,
		AccountNumber: c.Params("an"),
	}
	if err := h.validator.Validate(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"data":    nil,
			"message": err.Error(),
		})
	}

	// extended validation
	if startDate != nil && endDate != nil && startDate.After(*endDate) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"data":    nil,
			"message": "Start date cannot be later than end date",
		})
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.UTC().Location())
	if (startDate != nil && startDate.After(today)) || (endDate != nil && endDate.After(today)) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"data":    nil,
			"message": "Date cannot be later than today",
		})
	}

	res, err := h.service.GetTrxHistory(c.Context(), payload)
	if err != nil {
		if err.Error() == "Invalid account number" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"data":    res.Data,
				"meta":    res.Meta,
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"data":    res.Data,
			"meta":    res.Meta,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"data":    res.Data,
		"meta":    res.Meta,
		"message": "Successfully get transaction history",
	})
}
