package service

import (
	"errors"
	"payment/config"
	"payment/model"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransServiceImpl struct {
	Client        snap.Client
	CoreApiClient coreapi.Client
}

func NewMidtransService(configuration config.Config) MidtransService {
	mid := &MidtransServiceImpl{}
	mid.Client.New(configuration.Get("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)
	mid.CoreApiClient.New(configuration.Get("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)
	return mid
}

// createTokenTransactionWithGateway implements MidtransPayment
func (m *MidtransServiceImpl) CreateTokenTransactionWithGateway(request *snap.Request) (string, error) {

	resp, err := m.Client.CreateTransactionToken(request)
	if err != nil {
		return "", err
	}
	return resp, nil
}

/*
@param array
@return void
function untuk update setelah selesai pembayaran dan callback dari midtrans
status pembayaran
0 = belum proses
1 = dalam proses/pending
2 = sukses
3 = gagal
4 = expired
5 = cancel
*/
func (m *MidtransServiceImpl) Notification(request model.MidtransNotificationRequest) (int, error) {
	status := make(chan int, 1)
	if request.OrderId == "" {
		// do something when key `order_id` not found
		return 0, errors.New("transaction not found")
	}

	// 4. Check transaction to Midtrans with param orderId
	transactionStatusResp, e := m.CoreApiClient.CheckTransaction(request.OrderId)
	if e != nil {
		// http.Error(w, e.GetMessage(), http.StatusInternalServerError)
		return 0, e
	} else {
		if transactionStatusResp != nil {
			// 5. Do set transaction status based on response from check transaction status
			if transactionStatusResp.TransactionStatus == "capture" {
				if transactionStatusResp.FraudStatus == "challenge" {
					// TODO set transaction status on your database to 'challenge'
					// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
					status <- 8
				} else if transactionStatusResp.FraudStatus == "accept" {
					// TODO set transaction status on your database to 'success'
					status <- 8
				}
			} else if transactionStatusResp.TransactionStatus == "settlement" {
				// TODO set transaction status on your databaase to 'success'
				status <- 2
			} else if transactionStatusResp.TransactionStatus == "deny" {
				// TODO you can ignore 'deny', because most of the time it allows payment retries
				// and later can become success
				status <- 3
			} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
				// TODO set transaction status on your databaase to 'failure'
				status <- 4
			} else if transactionStatusResp.TransactionStatus == "pending" {
				// TODO set transaction status on your databaase to 'pending' / waiting payment
				status <- 1
			}
		}
	}
	return <-status, nil
}
