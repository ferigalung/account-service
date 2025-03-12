package main

import (
	"flag"
	"fmt"

	"github.com/ferigalung/account-service/config"
	"github.com/ferigalung/account-service/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
)

func main() {
	port := flag.Int("port", 3000, "REST API Port")
	flag.Parse()

	// load env
	cfg := config.LoadConfig()

	// init fiber
	app := fiber.New(fiber.Config{
		AppName: "Account Service - Ihsan Solusi",
	})

	// middleware
	app.Use(healthcheck.New())
	app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			cfg.BasicAuth.Username: cfg.BasicAuth.Password,
		},
	}))

	// routes
	g := app.Group("/api/v1")
	g.Post("/register", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})

	logger.Log("info", fmt.Sprintf("%s is running properly on port %d", app.Config().AppName, *port), nil)
	app.Listen(fmt.Sprintf(":%d", *port))
}
