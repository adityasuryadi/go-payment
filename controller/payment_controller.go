package controller

import (
	"payment/config"
	"payment/model"
	"payment/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PaymentController struct {
	PaymentService service.PaymentService
}

func NewPaymentController(paymentService service.PaymentService) PaymentController {
	return PaymentController{
		PaymentService: paymentService,
	}
}

func (controller *PaymentController) Route(app *fiber.App) {
	app.Post("payment", controller.Create)
	app.Post("callback", controller.CallbackHandle)
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

func (controller *PaymentController) CallbackHandle(ctx *fiber.Ctx) error {
	var request model.CallbackFaspayRequest
	ctx.BodyParser(&request)
	code, response := controller.PaymentService.UpdatePayment(request)
	responseCode, _ := strconv.Atoi(code)
	return ctx.Status(responseCode).JSON(response)
}
