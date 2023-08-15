package main

import (
	"ANTRIQUE/payment/config"
	"ANTRIQUE/payment/controller"
	"ANTRIQUE/payment/repository"
	"ANTRIQUE/payment/service"

	"github.com/gofiber/fiber/v2"
)

func main() {
	configuration := config.New(".env")
	database := config.NewPostgresDB(configuration)
	paymentRepository := repository.NewPaymentRepository(database)
	paymentService := service.NewPaymentService(&paymentRepository, database)
	paymentController := controller.NewPaymentController(&paymentService)

	app := fiber.New()
	paymentController.Route(app)

	app.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":5000")
}
