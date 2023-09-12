package service

import (
	"payment/model"

	"github.com/midtrans/midtrans-go/snap"
)

type MidtransService interface {
	CreateTokenTransactionWithGateway(requset *snap.Request) (string, error)
	Notification(request model.MidtransNotificationRequest) (int, error)
}
