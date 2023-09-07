package controller

import (
	"fmt"
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
	app.Post("midtrans", controller.createTokenTransactionWithGateway)
	app.Post("midtranscallback", controller.CallbackMidtrans)
}

func (controller *PaymentController) Create(ctx *fiber.Ctx) error {
	var request model.CreatePaymentRequest
	ctx.BodyParser(&request)
	valid := config.NewValidation()
	errValidation := valid.ValidateRequest(request)
	if errValidation != nil {
		return ctx.Status(400).JSON(model.GetResponse("400", errValidation))
	}

	code, response := controller.PaymentService.CreatePayment(request)
	responseCode, _ := strconv.Atoi(code)
	return ctx.Status(responseCode).JSON(model.GetResponse(code, response))
}

func (controller *PaymentController) CallbackHandle(ctx *fiber.Ctx) error {
	var request model.CallbackFaspayRequest
	ctx.BodyParser(&request)
	code, response := controller.PaymentService.UpdatePayment(request)
	responseCode, _ := strconv.Atoi(code)
	return ctx.Status(responseCode).JSON(response)
}

func (controller *PaymentController) createTokenTransactionWithGateway(ctx *fiber.Ctx) error {
	var request model.CreatePaymentRequest
	ctx.BodyParser(&request)

	valid := config.NewValidation()
	errValidation := valid.ValidateRequest(request)
	if errValidation != nil {
		return ctx.Status(400).JSON(model.GetResponse("400", errValidation))
	}

	code, resp := controller.PaymentService.GenerateSnapToken(request)
	responseCode, _ := strconv.Atoi(code)
	return ctx.Status(responseCode).JSON(model.GetResponse(code, model.SnapResponse{Token: resp.(string)}))
}

func (controller *PaymentController) CallbackMidtrans(ctx *fiber.Ctx) error {
	var request model.MidtransNotificationRequest
	ctx.BodyParser(&request)
	fmt.Println(request)
	code, resp := controller.PaymentService.CallbackMidtrans(request)
	responseCode, _ := strconv.Atoi(code)
	return ctx.Status(responseCode).JSON(model.GetResponse(code, resp))
}
