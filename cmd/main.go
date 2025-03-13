package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ferigalung/account-service/config"
	"github.com/ferigalung/account-service/internal/handler"
	"github.com/ferigalung/account-service/internal/repository"
	"github.com/ferigalung/account-service/internal/service"
	"github.com/ferigalung/account-service/pkg/database"
	"github.com/ferigalung/account-service/pkg/logger"
	Validator "github.com/ferigalung/account-service/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// set default service port, can be overrided by argumen parser
	port := flag.Int("port", 9000, "REST API Port")
	flag.Parse()

	// load env
	cfg := config.LoadConfig()

	// init fiber
	app := fiber.New(fiber.Config{
		AppName: "Account Service - Ihsan Solusi",
	})

	// middleware
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	app.Use(healthcheck.New())
	app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			cfg.BasicAuth.Username: cfg.BasicAuth.Password,
		},
		Unauthorized: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"data":    nil,
				"message": "Unauthorized access",
			})
		},
	}))

	// helper init
	validator := Validator.NewValidator()
	db := database.NewConnectionPool(cfg.DB)

	// repo, service, and handler init
	accountRepo := repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepo)
	accountHandler := handler.NewAccountHandler(accountService, validator)
	balanceRepo := repository.NewBalanceRepository(db)
	balanceService := service.NewBalanceService(balanceRepo)
	balanceHandler := handler.NewBalanceHandler(balanceService, validator)
	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := service.NewTransactionService(transactionRepo, accountRepo, balanceRepo)
	transactionHandler := handler.NewTransactionHandler(transactionService, validator)

	// routes
	g := app.Group("/api/v1")
	g.Post("/daftar", accountHandler.CreateAccount)
	g.Put("/tabung", transactionHandler.Deposit)
	g.Put("/tarik", transactionHandler.Withdraw)
	g.Get("/saldo/:an", balanceHandler.GetBalance)
	g.Get("/mutasi/:an", transactionHandler.GetTrxHistory)

	// graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sig
		logger.Log("info", "Shutting down service...", nil)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		db.Close()

		if err := app.ShutdownWithContext(ctx); err != nil {
			logger.Log("error", "Error during shutdown", fiber.Map{"error": err.Error()})
		} else {
			logger.Log("info", "Service shut down gracefully", nil)
		}
	}()

	logger.Log("info", fmt.Sprintf("%s is running properly on port %d", app.Config().AppName, *port), nil)
	if err := app.Listen(fmt.Sprintf(":%d", *port)); err != nil {
		logger.Log("fatal", "Server is not running", fiber.Map{"port": *port})
		os.Exit(1)
	}
}
