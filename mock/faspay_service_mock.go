package mocks

import (
	"errors"
	"payment/model"

	"github.com/stretchr/testify/mock"
)

type FaspayServiceMock struct {
	Mock mock.Mock
}

func (service *FaspayServiceMock) CreatePaymentExpress(request model.CreateFaspayPaymentRequest) (response *model.CreatePaymentFaspayResponse, err error) {
	arguments := service.Mock.Called(request)
	if arguments.Get(0) == nil {
		return nil, errors.New("not found")
	} else {
		response := arguments.Get(0).(*model.CreatePaymentFaspayResponse)
		return response, nil
	}
}
