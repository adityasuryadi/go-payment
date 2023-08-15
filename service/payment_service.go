package service

import "ANTRIQUE/payment/model"

type PaymentService interface {
	CreatePayment(request model.CreatePaymentRequest) (string, interface{})
}
