package controller

import (
	"ANTRIQUE/payment/config"
	"ANTRIQUE/payment/model"
	"ANTRIQUE/payment/service"

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
	valid := config.NewValidation()
	errValidation := valid.ValidateRequest(request)
	if errValidation != nil {
		return ctx.JSON(model.GetResponse("400", errValidation))
	}

	code, response := controller.PaymentService.CreatePayment(request)

	return ctx.JSON(model.GetResponse(code, response))
}
