package model

type CreatePaymentRequest struct {
	CustName    string `json:"cust_name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Phone       string `json:"phone" validate:"required,number,min=8"`
	ServiceName string `json:"service_name" validate:"required"`
	ServiceCode string `json:"service_code" validate:"required"`
	ServiceId   string `json:"service_id" validate:"required"`
	BookingDate string `json:"booking_date" validate:"required,gtenow"`
}

type CreateFaspayPaymentRequest struct {
	MerchantId  string  `json:"merchant_id"`
	BillNo      string  `json:"bill_no"`
	BillDate    string  `json:"bill_date"`
	BillExpired string  `json:"bill_expired"`
	BillTotal   float64 `json:"bill_total"`
	BillDesc    string  `json:"bill_desc"`
	CustNo      string  `json:"cust_no"`
	CustName    string  `json:"cust_name"`
	ReturnUrl   string  `json:"return_url"`
	Product     string  `json:"product"`
	Qty         int64   `json:"qty"`
	Amount      float64 `json:"amount"`
	Signature   string  `json:"signature"`
	Msisdn      string  `json:"msisdn"`
	Email       string  `json:"email"`
	Item        []Item  `json:"item"`
	PayType     string  `json:"pay_type"`
	Terminal    string  `json:"terminal"`
}

type CreatePaymentFaspayResponse struct {
	BillNo       string `json:"bill_no"`
	MerchantId   string `json:"merchant_id"`
	Merchant     string `json:"merchant"`
	ResponseCode string `json:"response_code"`
	ResponseDesc string `json:"response_desc"`
	RedirectUrl  string `json:"redirect_url"`
}

type Item struct {
	Product string  `json:"product"`
	Qty     int     `json:"qty"`
	Amount  float64 `json:"amount"`
}

type CallbackFaspayRequest struct {
	Request           string `json:"request"`
	TrxId             string `json:"trx_id"`
	Merchant          string `json:"merchant"`
	MerchantId        string `json:"merchant_id"`
	BillNo            string `json:"bill_no"`
	PaymentReff       string `json:"payment_reff"`
	PaymentDate       string `json:"payment_date"`
	PaymentStatusCode string `json:"payment_status_code"`
	PaymentStatusDesc string `json:"payment_status_desc"`
	BillTotal         string `json:"bill_total"`
	PaymentTotal      string `json:"payment_total"`
	PaymentChannelUid string `json:"payment_channel_uid"`
	PaymentChannel    string `json:"payment_channel"`
	Amount            string `json:"amount"`
	Signature         string `json:"signature"`
}

type FaspayNotifResponse struct {
	Response     string `json:"response"`
	TrxId        string `json:"trx_id"`
	MerchantId   string `json:"merchant_id"`
	Merchant     string `json:"merchant"`
	BillNo       string `json:"bill_no"`
	ResponseCode string `json:"response_code"`
	ResponseDesc string `json:"response_desc"`
	ResponseDate string `json:"response_date"`
}
