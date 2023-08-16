package controller

import (
	"ANTRIQUE/payment/config"
	"ANTRIQUE/payment/model"
	"ANTRIQUE/payment/service"
	"strconv"

	"github.com/go-playground/validator/v10"
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
	validate := validator.New()
	err := validate.Struct(request)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)

		out := make([]config.ErrorMessage, len(validationErrors))
		for i, fieldError := range validationErrors {
			out[i] = config.ErrorMessage{
				Field:   fieldError.Field(),
				Message: config.GetErrorMsg(fieldError),
			}
		}
		return ctx.JSON(model.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   out,
		})
	}

	ctx.BodyParser(&request)

	code, response := controller.PaymentService.CreatePayment(request)
	responseCode, _ := strconv.Atoi(code)
	return ctx.JSON(model.WebResponse{
		Code:   int(responseCode),
		Status: "OK",
		Data:   response,
	})
}
