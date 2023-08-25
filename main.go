package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := InitializeApp(".env")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	log.Fatal(app.Listen(":5000"))
}
