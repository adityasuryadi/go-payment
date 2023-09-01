package service

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"payment/config"
	"payment/entity"
	"payment/helper"
	"payment/model"
	"payment/repository"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/natefinch/lumberjack"
	"gorm.io/gorm"
)

const price float64 = 2000

func NewPaymentService(repository repository.PaymentRepository, db *gorm.DB, faspayService FaspayService, pointRepository repository.PointRespository, midtransPayment config.MidtransPayment) PaymentService {
	return &PaymentServiceImpl{
		PaymentRepository: repository,
		db:                db,
		FaspayService:     faspayService,
		PointRepository:   pointRepository,
		MidtransPayment:   midtransPayment,
	}
}

type PaymentServiceImpl struct {
	PaymentRepository repository.PaymentRepository
	db                *gorm.DB
	FaspayService     FaspayService
	PointRepository   repository.PointRespository
	MidtransPayment   config.MidtransPayment
}

// UpdatePayment implements PaymentService
func (paymentService *PaymentServiceImpl) UpdatePayment(request model.CallbackFaspayRequest) (string, interface{}) {

	billNo := request.BillNo
	payment, err := paymentService.PaymentRepository.FindPaymentByBillNo(billNo)
	merchantId := os.Getenv("FASPAY_MERCHANT_ID")

	if err != nil {
		return "404", payment
	}

	if payment.StatusId == 2 {
		return "400", "has already pay"
	}

	// log callback
	log.SetOutput(&lumberjack.Logger{
		Filename:   "./var/log/faspaycallback.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     1,    //days
		Compress:   true, // disabled by default
	})
	log.Print(request)

	if payment.Signature != request.Signature {
		return "400", "invalid signature"
	}

	// create ticket api

	type RequestGenerateTicket struct {
		UserId      int    `json:"user_id"`
		ServiceCode string `json:"string"`
		QueueName   string `json:"queue_name"`
		ServiceTime string `json:"service_time"`
		QueueHp     string `json:"queue_hp"`
		Note        string `json:"note"`
	}

	url := os.Getenv("CREATE_TICKET_URL")
	requestTicket := RequestGenerateTicket{
		UserId:      64,
		ServiceCode: "SPnjWy",
		QueueName:   payment.Name,
		ServiceTime: payment.BookingDate.String(),
		QueueHp:     payment.Phone,
		Note:        "Catatan",
	}

	client := resty.New()
	resp, err := client.R().
		SetFormData(map[string]string{
			"user_id":      "64",
			"service_code": requestTicket.ServiceCode,
			"queue_name":   requestTicket.QueueName,
			"service_time": payment.BookingDate.String(),
			"note":         "-",
			"queue_hp":     payment.Phone,
		}).
		SetHeader("Accept", "application/json").
		Post(url)

	if err != nil {
		return "500", err.Error()
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

	responseTicket := make(map[string]interface{})
	json.Unmarshal(resp.Body(), &responseTicket)
	if responseTicket["error"] == true {
		log.SetOutput(&lumberjack.Logger{
			Filename:   "./var/log/ticket.log",
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     1,    //days
			Compress:   true, // disabled by default
		})
		responseTicket["trx_id"] = payment.TrxId
		log.Print(responseTicket)

		// restore payment to point
		userPoint, err := paymentService.PointRepository.FindPointByUserId(30)
		if err != nil && !!errors.Is(err, gorm.ErrRecordNotFound) {
			return "500", err.Error()
		}
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			entityPoint := &entity.Point{
				UserId: 30,
				Point:  payment.BillTotal,
			}
			_, err := paymentService.PointRepository.InsertOrUpdate(&entity.Point{
				UserId: 30,
				Point:  entityPoint.Point + payment.BillTotal,
			})
			if err != nil {
				return "500", err.Error()
			}
		}
		paymentService.PointRepository.InsertOrUpdate(&entity.Point{
			UserId: 30,
			Point:  userPoint.Point + payment.BillTotal,
		})
		return "400", responseTicket["error_msg"]
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
		BillDesc:    "Booking Antrian " + request.ServiceName,
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
				Product: request.ServiceName,
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
			ServiceId:     request.ServiceId,
			BookingDate:   bookingDate,
			RedirectUrl:   faspayResponse.RedirectUrl,
			BillNoCounter: billNoCounter,
			Qty:           1,
			BillNo:        faspayResponse.BillNo,
			BillTotal:     price,
			StatusId:      1,
			Signature:     string(fmt.Sprintf("%x", signature)),
			ServiceCode:   request.ServiceCode,
		}

		err := paymentService.PaymentRepository.Store(tx, &payment)
		if err != nil {
			return "500", err.Error()
		}
	}

	if faspayResponse.ResponseCode != "00" {
		return "400", nil
	}
	tx.Commit()
	return "200", faspayResponse
}

func (paymentService *PaymentServiceImpl) GenerateBillNo(tx *gorm.DB) (billNo string, billNoCounter int) {
	payment, err := paymentService.PaymentRepository.GetLastPaymentToday(tx)
	curdate := time.Now().Format("20060102")
	if errors.Is(err, gorm.ErrRecordNotFound) {
		billNo = fmt.Sprintf("INV-%s%d", curdate, 1)
		billNoCounter = 1
	} else {
		billNoCounter += payment.BillNoCounter + 1
		billNo = fmt.Sprintf("INV-%s%d", curdate, billNoCounter)
	}
	return billNo, billNoCounter
}

func (paymentService *PaymentServiceImpl) GenerateSnapToken(request model.CreatePaymentRequest) (string, interface{}) {
	tx := paymentService.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	billNo, billNoCounter := paymentService.GenerateBillNo(tx)
	bookingDate, _ := time.Parse("2006-01-02", request.BookingDate)

	custAddress := &midtrans.CustomerAddress{
		FName:       request.CustName,
		LName:       "",
		Phone:       request.Phone,
		Address:     "",
		City:        "",
		Postcode:    "",
		CountryCode: "IDN",
	}

	fmt.Println("service_name", request.ServiceName)
	// Initiate Snap Request
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  billNo,
			GrossAmt: int64(price),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName:    request.CustName,
			LName:    "",
			Email:    request.Email,
			Phone:    request.Phone,
			BillAddr: custAddress,
			ShipAddr: custAddress,
		},
		EnabledPayments: snap.AllSnapPaymentType,
		Items: &[]midtrans.ItemDetails{
			{
				ID:    request.ServiceCode,
				Price: int64(price),
				Qty:   1,
				Name:  request.ServiceName,
			},
		},
	}
	token, err := paymentService.MidtransPayment.CreateTokenTransactionWithGateway(snapReq)
	if err != nil {
		return "400", err.Error()
	}
	payment := entity.Payment{
		Name:          request.CustName,
		Phone:         request.Phone,
		Email:         request.Email,
		ServiceId:     request.ServiceId,
		BookingDate:   bookingDate,
		RedirectUrl:   "",
		BillNoCounter: billNoCounter,
		Qty:           1,
		BillNo:        billNo,
		BillTotal:     price,
		StatusId:      1,
		ServiceCode:   request.ServiceCode,
		SnapToken:     token,
	}

	err = paymentService.PaymentRepository.Store(tx, &payment)
	if err != nil {
		return "500", err.Error()
	}
	tx.Commit()

	return "200", token
}

func (paymentService *PaymentServiceImpl) CallbackMidtrans(request model.MidtransNotificationRequest) (string, interface{}) {
	status, err := paymentService.MidtransPayment.Notification(request)
	if err != nil && errors.Is(err, err.(*midtrans.Error)) {
		return "400", err.(*midtrans.Error).GetMessage()
	}
	payment, err := paymentService.PaymentRepository.FindPaymentByBillNo(request.OrderId)
	if err != nil {
		return "404", payment
	}

	payment.StatusId = status

	err = paymentService.PaymentRepository.Update(payment)

	if err != nil {
		return "500", err
	}
	return "200", payment
}
