package model

type CreatePaymentRequest struct {
	CustName    string `json:"cust_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	ServiceName string `json:"service_name"`
	ServiceId   string `json:"service_id"`
	BookingDate string `json:"booking_date"`
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
