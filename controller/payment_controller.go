package controller

import (
	"ANTRIQUE/payment/model"
	"ANTRIQUE/payment/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PaymentController struct {
	PaymentService service.PaymentService
}

func NewPaymentController(paymentService *service.PaymentService) PaymentController {
	return PaymentController{
		PaymentService: *paymentService,
	}
}

func (controller *PaymentController) Route(app *fiber.App) {
	app.Post("payment", controller.Create)
}

func (controller *PaymentController) Create(ctx *fiber.Ctx) error {
	var request model.CreatePaymentRequest
	ctx.BodyParser(&request)

	code, response := controller.PaymentService.CreatePayment(request)
	responseCode, _ := strconv.Atoi(code)
	return ctx.JSON(model.WebResponse{
		Code:   responseCode,
		Status: "OK",
		Data:   response,
	})
}
