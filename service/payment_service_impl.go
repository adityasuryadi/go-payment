package service

import (
	"ANTRIQUE/payment/config"
	"ANTRIQUE/payment/entity"
	"ANTRIQUE/payment/helper"
	"ANTRIQUE/payment/model"
	"ANTRIQUE/payment/repository"
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
)

func NewPaymentService(repository *repository.PaymentRepository, db *gorm.DB) PaymentService {
	return &PaymentServiceImpl{
		PaymentRepository: *repository,
		db:                db,
	}
}

type PaymentServiceImpl struct {
	PaymentRepository repository.PaymentRepository
	db                *gorm.DB
}

// CreatePayment implements PaymentService
func (paymentService *PaymentServiceImpl) CreatePayment(request model.CreatePaymentRequest) (string, interface{}) {
	tx := paymentService.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	response_code := make(chan string, 1)

	merchantId := os.Getenv("FASPAY_MERCHANT_ID")
	userId := os.Getenv("FASPAY_USER_ID")
	password := os.Getenv("FASPAY_PASSWORD")

	url := "https://xpress.faspay.co.id/v4/post"
	if os.Getenv("APP_ENV") == "dev" {
		url = "https://xpress-sandbox.faspay.co.id/v4/post"
	}

	var price float64 = 2000
	billNo, billNoCounter := helper.GenerateBillNo(tx)

	shaEncrypt := sha1.New()
	md5Encrypt := md5.New()

	plainSignature := userId + password + billNo + strconv.Itoa(int(price))

	md5Encrypt.Write([]byte(plainSignature))
	md5Signature := md5Encrypt.Sum(nil)

	shaEncrypt.Write([]byte(string(fmt.Sprintf("%x", md5Signature))))
	signature := shaEncrypt.Sum(nil)

	now := time.Now()
	nowString := now.Format("2006-01-02 15:04:05")
	requestFaspay := model.CreateFaspayPaymentRequest{
		MerchantId:  merchantId,
		BillNo:      billNo,
		BillDate:    nowString,
		BillExpired: now.Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
		Signature:   string(fmt.Sprintf("%x", signature)),
		BillDesc:    "Booking Antrian",
		BillTotal:   float64(price),
		CustNo:      request.Phone,
		CustName:    request.CustName,
		ReturnUrl:   os.Getenv("FASPAY_CALLBACK_URL"),
		Product:     request.ServiceName,
		Amount:      float64(price),
		Msisdn:      request.Phone,
		Terminal:    "10",
		PayType:     "1",
		Email:       request.Email,
		Item: []model.Item{
			{
				Product: "antrian",
				Qty:     1,
				Amount:  float64(price),
			},
		},
	}
	client := resty.New()
	resp, err := client.R().
		SetBody(requestFaspay).
		Post(url)
	if err != nil {
		fmt.Println(err)
	}

	response := make(map[string]interface{})
	json.Unmarshal(resp.Body(), &response)
	bookingDate, _ := time.Parse("2006-01-02", request.BookingDate)

	response_code <- "200"
	if response["response_code"] == "00" {
		payment := entity.Payment{
			Name:          request.CustName,
			Phone:         request.Phone,
			Email:         request.Email,
			BookingDate:   bookingDate,
			RedirectUrl:   response["redirect_url"].(string),
			BillNoCounter: billNoCounter,
			Qty:           1,
			BillNo:        billNo,
			BillTotal:     price,
			StatusId:      1,
		}
		err = paymentService.PaymentRepository.Store(tx, &payment)
		if err != nil {
			response_code <- "500"
			response = nil
			config.WriteLog(err)
		}
	}

	if response["response_code"] != "00" {
		response_code <- "400"
	}
	tx.Commit()
	return <-response_code, response
}
