package repository

import (
	"payment/entity"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	Store(tx *gorm.DB, payment *entity.Payment) error
	FindPaymentByBillNo(billNo string) (payment *entity.Payment, err error)
	Update(payment *entity.Payment) error
}
