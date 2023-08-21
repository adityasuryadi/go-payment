package repository

import (
	"ANTRIQUE/payment/entity"
	"errors"

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

// Update implements PaymentRepository
func (repository *PaymentReposirotyImpl) Update(payment *entity.Payment) error {
	err := repository.db.Save(&payment).Error
	if err != nil {
		return err
	}
	return nil
}

// Store implements PaymentRepository
func (repository *PaymentReposirotyImpl) Store(tx *gorm.DB, payment *entity.Payment) error {
	result := tx.Create(&payment)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repository *PaymentReposirotyImpl) FindPaymentByBillNo(billNo string) (payment *entity.Payment, err error) {
	result := repository.db.Where("bill_no = ?", billNo).First(&payment)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return payment, nil
}
