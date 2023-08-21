package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := InitializeApp(".env")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":5000")
}
