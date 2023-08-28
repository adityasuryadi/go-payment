package repository

import (
	"errors"
	"payment/entity"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (repository *PaymentReposirotyImpl) GetLastPaymentToday(tx *gorm.DB) (*entity.Payment, error) {
	var payment entity.Payment
	today := time.Now().Format("2006-01-02")
	result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(clause.Expr{SQL: "DATE(created_at) = ?", Vars: []interface{}{today}}).
		Order("bill_no_counter desc").First(&payment)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.Error != nil {
		return nil, result.Error
	}
	return &payment, nil
}
