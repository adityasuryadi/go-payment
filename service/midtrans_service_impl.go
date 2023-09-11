package service

import "payment/config"

type MidtransService interface {
	GenerateSnapRequest()
}

type MidtransServiceImpl struct {
	mConfig config.MidtransPayment
}

// GenerateSnapRequest implements MidtransService
func (*MidtransServiceImpl) GenerateSnapRequest() {
	panic("unimplemented")
}

func NewMidtransService(mConifg config.MidtransPayment) MidtransService {
	return &MidtransServiceImpl{
		mConfig: mConifg,
	}
}
