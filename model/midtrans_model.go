package model

type SnapResponse struct {
	Token string `json:"token"`
}

type MidtransPaymentRequest struct {
	BillNo    string  `json:"bill_no"`
	BillTotal float64 `json:"bill_total"`
	Email     string  `json:"email"`
	PHone     string  `json:"phone"`
	Item      string  `json:"item"`
	Qty       int     `json:"qty"`
}

type MidtransNotificationRequest struct {
	TransactionTime        string  `json:"transaction_time"`
	TransactionStatus      string  `json:"transaction_status"`
	TransactionId          string  `json:"transaction_id"`
	StatusMessage          string  `json:"status_message"`
	StatusCode             string  `json:"status_code"`
	SigantureKey           string  `json:"signature_key"`
	PaymentType            string  `json:"payment_type"`
	OrderId                string  `json:"order_id"`
	MerchantId             string  `json:"merchant_id"`
	MaskedCard             string  `json:"masked_card"`
	GrossAmount            float64 `json:"gross_amount"`
	FraudStatus            string  `json:"fraud_status"`
	Eci                    string  `json:"eci"`
	Currency               string  `json:"currency"`
	ChannelResponseMessage string  `json:"channel_response_message"`
	ChannelResponseCode    string  `json:"channel_response_code"`
	CardType               string  `json:"card_type"`
	Bank                   string  `json:"bank"`
	ApprovalCode           string  `json:"approval_code"`
}
