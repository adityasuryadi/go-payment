package repository

import (
	"ANTRIQUE/payment/entity"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	Store(tx *gorm.DB, payment *entity.Payment) error
}
