package mocks

import (
	"errors"
	"payment/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type PaymentRepositoryMock struct {
	Mock mock.Mock
}

func (repository *PaymentRepositoryMock) FindPaymentByBillNo(billNo string) (payment *entity.Payment, err error) {
	arguments := repository.Mock.Called(billNo)
	if arguments.Get(0) == nil {
		return nil, errors.New("not found")
	} else {
		return arguments.Get(0).(*entity.Payment), nil
	}
}

func (repository *PaymentRepositoryMock) Store(tx *gorm.DB, payment *entity.Payment) error {
	panic("")
}

func (repository *PaymentRepositoryMock) Update(payment *entity.Payment) error {
	arguments := repository.Mock.Called(payment)
	if arguments.Get(0) == nil {
		return errors.New("failed update payment")
	} else {
		return nil
	}
}

func (repository *PaymentRepositoryMock) GetLastPaymentToday(tx *gorm.DB) (*entity.Payment, error) {
	arguments := repository.Mock.Called(tx)
	if arguments.Get(0) == nil {
		return nil, gorm.ErrRecordNotFound
	} else {
		payment := arguments.Get(0).(entity.Payment)
		return &payment, nil
	}
}
