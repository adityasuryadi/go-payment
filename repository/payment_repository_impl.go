package repository

import (
	"ANTRIQUE/payment/entity"
	"fmt"

	"gorm.io/gorm"
)

func NewPaymentRepository(database *gorm.DB) PaymentRepository {
	return &PaymentReposirotyImpl{
		db: database,
	}
}

type PaymentReposirotyImpl struct {
	db *gorm.DB
}

// Store implements PaymentRepository
func (repository *PaymentReposirotyImpl) Store(tx *gorm.DB, payment *entity.Payment) error {
	result := tx.Create(&payment)
	if result.Error != nil {
		return result.Error
	}
	fmt.Println(result)
	return nil
}
