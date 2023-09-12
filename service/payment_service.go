package service

import (
	"payment/model"

	"gorm.io/gorm"
)

type PaymentService interface {
	CreatePayment(request model.CreatePaymentRequest) (string, interface{})
	UpdatePayment(request model.CallbackFaspayRequest) (string, interface{})
	GenerateBillNo(tx *gorm.DB) (billNo string, billNoCounter int)
	GenerateSnapToken(request model.CreatePaymentRequest) (string, interface{})
	CallbackMidtrans(request model.MidtransNotificationRequest) (string, interface{})
	GetListPaymentType() (string, interface{})
}
