package service

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"payment/entity"
	"payment/helper"
	"payment/model"
	"payment/repository"
	"strconv"
	"time"

	"github.com/natefinch/lumberjack"
	"gorm.io/gorm"
)

const price float64 = 2000

func NewPaymentService(repository repository.PaymentRepository, db *gorm.DB, faspayService FaspayService) PaymentService {
	return &PaymentServiceImpl{
		PaymentRepository: repository,
		FaspayService:     faspayService,
		db:                db,
	}
}

type PaymentServiceImpl struct {
	PaymentRepository repository.PaymentRepository
	db                *gorm.DB
	FaspayService     FaspayService
}

// UpdatePayment implements PaymentService
func (paymentService *PaymentServiceImpl) UpdatePayment(request model.CallbackFaspayRequest) (string, interface{}) {

	billNo := request.BillNo
	payment, err := paymentService.PaymentRepository.FindPaymentByBillNo(billNo)
	merchantId := os.Getenv("FASPAY_MERCHANT_ID")

	// log callback
	log.SetOutput(&lumberjack.Logger{
		Filename:   "./var/log/faspaycallback.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     1,    //days
		Compress:   true, // disabled by default
	})
	log.Print(request)

	if err != nil {
		return "404", nil
	}

	if payment.Signature != request.Signature {
		return "400", nil
	}

	paymentStatus, _ := strconv.Atoi(request.PaymentStatusCode)
	paymentChannelUid, _ := strconv.Atoi(request.PaymentChannelUid)

	payment.StatusId = paymentStatus
	payment.TrxId = request.TrxId
	payment.PaymentChannel = request.PaymentChannel
	payment.PaymentChannelUid = paymentChannelUid
	err = paymentService.PaymentRepository.Update(payment)

	if err != nil {
		return "500", err
	}

	response := model.FaspayNotifResponse{
		Response:     request.Request,
		TrxId:        payment.TrxId,
		MerchantId:   merchantId,
		Merchant:     request.Merchant,
		BillNo:       payment.BillNo,
		ResponseCode: "00",
		ResponseDesc: "Success",
		ResponseDate: helper.GetNowStringFormat(),
	}

	return "200", response
}

// CreatePayment implements PaymentService
func (paymentService *PaymentServiceImpl) CreatePayment(request model.CreatePaymentRequest) (string, interface{}) {
	tx := paymentService.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()
	nowString := now.Format("2006-01-02 15:04:05")
	billNo, billNoCounter := paymentService.GenerateBillNo(tx)
	requestFaspay := model.CreateFaspayPaymentRequest{
		BillNo:      billNo,
		BillDate:    nowString,
		BillExpired: now.Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
		BillDesc:    "Booking Antrian",
		BillTotal:   float64(price),
		CustNo:      request.Phone,
		CustName:    request.CustName,
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

	faspayResponse, err := paymentService.FaspayService.CreatePaymentExpress(requestFaspay)
	if err != nil {
		log.Fatal(err)
	}
	userId := os.Getenv("FASPAY_USER_ID")
	password := os.Getenv("FASPAY_PASSWORD")
	shaEncrypt := sha1.New()
	md5Encrypt := md5.New()

	plainSignature := userId + password + faspayResponse.BillNo + strconv.Itoa(int(price))

	md5Encrypt.Write([]byte(plainSignature))
	md5Signature := md5Encrypt.Sum(nil)

	shaEncrypt.Write([]byte(string(fmt.Sprintf("%x", md5Signature))))
	signature := shaEncrypt.Sum(nil)
	bookingDate, _ := time.Parse("2006-01-02", request.BookingDate)
	if faspayResponse.ResponseCode == "00" {
		payment := entity.Payment{
			Name:          request.CustName,
			Phone:         request.Phone,
			Email:         request.Email,
			BookingDate:   bookingDate,
			RedirectUrl:   faspayResponse.RedirectUrl,
			BillNoCounter: billNoCounter,
			Qty:           1,
			BillNo:        faspayResponse.BillNo,
			BillTotal:     price,
			StatusId:      1,
			Signature:     string(fmt.Sprintf("%x", signature)),
		}

		paymentService.PaymentRepository.Store(tx, &payment)
		if err != nil {
			return "500", err.Error()
		}
	}

	if faspayResponse.ResponseCode != "00" {
		return "400", nil
	}
	tx.Commit()
	return "200", faspayResponse
	// return "wkwkwk", nil
}

func (paymentService *PaymentServiceImpl) GenerateBillNo(tx *gorm.DB) (billNo string, billNoCounter int) {
	billNo, billNoCounter = helper.GenerateBillNo(tx)
	return billNo, billNoCounter
}
