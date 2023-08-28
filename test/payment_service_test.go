package test

import (
	"crypto/md5"
	"crypto/sha1"
	"errors"
	"fmt"
	"payment/config"
	"payment/entity"
	"payment/helper"
	mocks "payment/mock"
	"payment/model"
	"payment/repository"
	"payment/service"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const projectDirName = "payment"

var paymentRepositoryMock = &mocks.PaymentRepositoryMock{Mock: mock.Mock{}}

var faspayServiceMock = &mocks.FaspayServiceMock{Mock: mock.Mock{}}

func TestFindPaymentNotFound(t *testing.T) {
	paymentRepositoryMock.Mock.On("FindPaymentByBillNo", "1").Return(nil, mock.Anything)

	payment, err := paymentRepositoryMock.FindPaymentByBillNo("1")
	assert.Nil(t, payment)
	assert.NotNil(t, err)
}

func TestFindPaymentFound(t *testing.T) {
	paymentRepositoryMock.Mock.On("FindPaymentByBillNo", "2").Return(entity.Payment{
		Email: "adit@mail.com",
	}, nil)

	payment, err := paymentRepositoryMock.FindPaymentByBillNo("2")
	assert.NotNil(t, payment)
	assert.Nil(t, err)
}

func TestCreatePayment(t *testing.T) {
	const price = 2000
	configApp := config.New(`\.env.test`)
	db := config.NewPostgresDB(configApp)
	paymentRepository := repository.NewPaymentRepository(db)
	billNo, _ := helper.GenerateBillNo(db)

	merchantId := configApp.Get("FASPAY_MERCHANT_ID")
	userId := configApp.Get("FASPAY_USER_ID")
	password := configApp.Get("FASPAY_PASSWORD")

	shaEncrypt := sha1.New()
	md5Encrypt := md5.New()

	plainSignature := userId + password + billNo + strconv.Itoa(int(price))

	md5Encrypt.Write([]byte(plainSignature))
	md5Signature := md5Encrypt.Sum(nil)

	shaEncrypt.Write([]byte(string(fmt.Sprintf("%x", md5Signature))))

	now := time.Now()
	nowString := now.Format("2006-01-02 15:04:05")

	request := model.CreatePaymentRequest{
		CustName:    "Adit",
		Email:       "adit@mail.com",
		Phone:       "08961243124",
		BookingDate: "2028-08-30",
	}

	mockRequestFaspay := model.CreateFaspayPaymentRequest{
		BillNo:      billNo,
		BillDate:    nowString,
		BillExpired: now.Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
		BillTotal:   price,
		BillDesc:    "Booking Antrian",
		CustNo:      request.Phone,
		CustName:    request.CustName,
		Product:     request.ServiceName,
		Amount:      float64(price),
		Email:       request.Email,
		Msisdn:      request.Phone,
		Item: []model.Item{
			{
				Product: "antrian",
				Qty:     1,
				Amount:  float64(price),
			},
		},
		PayType:  "1",
		Terminal: "10",
	}
	var paymentService = service.NewPaymentService(paymentRepository, db, faspayServiceMock)
	faspayServiceMock.Mock.On("CreatePaymentExpress", mockRequestFaspay).Return(&model.CreatePaymentFaspayResponse{
		BillNo:       billNo,
		MerchantId:   merchantId,
		Merchant:     "Antrique",
		ResponseCode: "00",
		ResponseDesc: "Sukses",
		RedirectUrl:  "https://s.faspay.co.id/2FK",
	}, nil)

	responseCode, response := paymentService.CreatePayment(request)
	assert.Equal(t, responseCode, "200")
	assert.NotNil(t, response)
	faspayServiceMock.Mock.AssertExpectations(t)

}

func TestUpdatePaymentBillNotFound(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		fmt.Print(err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("error opening GORM database: %v", err)
	}
	paymentRepositoryMock.Mock.On("FindPaymentByBillNo", "INV-202108222").Return(nil, errors.New("record not found"))
	var paymentService = service.NewPaymentService(paymentRepositoryMock, gormDB, faspayServiceMock)

	request := model.CallbackFaspayRequest{
		Request:           "Payment Notification",
		TrxId:             "3183540500001172",
		Merchant:          "31835",
		MerchantId:        "Sophia Store",
		BillNo:            "INV-202108222",
		PaymentReff:       "null",
		PaymentDate:       "2017-10-04 15:46:35",
		PaymentStatusCode: "2",
		PaymentStatusDesc: "Payment Success",
		BillTotal:         "5000000",
		PaymentTotal:      "5000000",
		PaymentChannelUid: "402",
		PaymentChannel:    "Permata Virtual Account",
		Signature:         "f0275409443913ec563ef2307897c233ce109455",
	}
	errCode, response := paymentService.UpdatePayment(request)
	assert.Equal(t, errCode, "404")
	assert.Nil(t, response)
}
func TestUpdatePaymentSignatureNotMatch(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		fmt.Print(err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("error opening GORM database: %v", err)
	}
	paymentRepositoryMock.Mock.On("FindPaymentByBillNo", "INV-202308222").Return(entity.Payment{
		Id:                uuid.MustParse("33a61a3d-88e8-484d-8061-3db0bff92e3a"),
		Name:              "Antrian",
		Phone:             "089656234771",
		Email:             "adit@mail.com",
		BookingDate:       time.Now(),
		RedirectUrl:       "http://antrique.com",
		Qty:               1,
		BillNo:            "INV-202308222",
		BillTotal:         2000,
		StatusId:          2,
		BillNoCounter:     22,
		TrxId:             "3183540500001172",
		PaymentChannelUid: 402,
		PaymentChannel:    "Permata Virtual Account",
		Signature:         "f1275409443913ec563ef2307897c233ce109456",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}, nil)
	var paymentService = service.NewPaymentService(paymentRepositoryMock, gormDB, faspayServiceMock)

	request := model.CallbackFaspayRequest{
		Request:           "Payment Notification",
		TrxId:             "3183540500001172",
		Merchant:          "31835",
		MerchantId:        "Sophia Store",
		BillNo:            "INV-202308222",
		PaymentReff:       "null",
		PaymentDate:       "2017-10-04 15:46:35",
		PaymentStatusCode: "2",
		PaymentStatusDesc: "Payment Success",
		BillTotal:         "5000000",
		PaymentTotal:      "5000000",
		PaymentChannelUid: "402",
		PaymentChannel:    "Permata Virtual Account",
		Signature:         "f0275409443913ec563ef2307897c233ce109455",
	}
	errCode, response := paymentService.UpdatePayment(request)
	assert.Equal(t, errCode, "400")
	assert.Nil(t, response)
}

func TestGenerateBillNoLastPaymentTodayIsNull(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		fmt.Print(err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("error opening GORM database: %v", err)
	}
	curdate := time.Now().Format("20060102")
	paymentRepositoryMock.Mock.On("GetLastPaymentToday", gormDB).Return(nil, err)
	var paymentService = service.NewPaymentService(paymentRepositoryMock, gormDB, faspayServiceMock)
	billNo, billNoCounter := paymentService.GenerateBillNo(gormDB)
	resultBillNo := fmt.Sprintf("INV-%s%d", curdate, 1)
	assert.Equal(t, billNoCounter, 1)
	assert.Equal(t, billNo, resultBillNo)
}

func TestGenerateBillNoLastPaymentTodayExist(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		fmt.Print(err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("error opening GORM database: %v", err)
	}
	curdate := time.Now().Format("20060102")
	paymentRepositoryMock.Mock.On("GetLastPaymentToday", gormDB).Return(entity.Payment{BillNo: fmt.Sprintf("INV-%s%d", curdate, 2), BillNoCounter: 2}, nil)
	var paymentService = service.NewPaymentService(paymentRepositoryMock, gormDB, faspayServiceMock)
	billNo, billNoCounter := paymentService.GenerateBillNo(gormDB)
	resultBillNo := fmt.Sprintf("INV-%s%d", curdate, 2)
	assert.Equal(t, billNoCounter, 2)
	assert.Equal(t, billNo, resultBillNo)
}
