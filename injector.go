//go:build wireinject
// +build wireinject

package main

import (
	"ANTRIQUE/payment/config"
	"ANTRIQUE/payment/controller"
	"ANTRIQUE/payment/exception"
	"ANTRIQUE/payment/repository"
	"ANTRIQUE/payment/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/wire"
)

var (
	paymentSet = wire.NewSet(repository.NewPaymentRepository, service.NewPaymentService, controller.NewPaymentController)
)

func InitializeApp(filenames ...string) *fiber.App {
	wire.Build(
		config.New,
		config.NewPostgresDB,
		paymentSet,
		NewServer,
	)
	return nil
}

func NewServer(paymentController controller.PaymentController) *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: exception.ErrorHandler})
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "*",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
	paymentController.Route(app)
	return app
}
