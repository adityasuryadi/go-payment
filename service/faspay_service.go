package service

import "payment/model"

type FaspayService interface {
	CreatePaymentExpress(request model.CreateFaspayPaymentRequest) (*model.CreatePaymentFaspayResponse, error)
}
