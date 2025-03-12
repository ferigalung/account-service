package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	port := flag.Int("port", 3000, "REST API Port")
	flag.Parse()

	app := fiber.New(fiber.Config{
		AppName: "Account Service - Ihsan Solusi",
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})

	log.Infof("Server starting at %s", time.Now().Format(time.RFC3339))
	app.Listen(fmt.Sprintf(":%d", *port))
}
