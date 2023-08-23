package test

import (
	"errors"
	"fmt"
	"payment/entity"
	"payment/model"
	"payment/service"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PaymentRepositoryMock struct {
	Mock mock.Mock
}

func (repository *PaymentRepositoryMock) FindPaymentByBillNo(billNo string) (payment *entity.Payment, err error) {
	arguments := repository.Mock.Called(billNo)
	if arguments.Get(0) == nil {
		return nil, errors.New("not found")
	} else {
		payment := arguments.Get(0).(entity.Payment)
		return &payment, nil
	}
}

func (repository *PaymentRepositoryMock) Store(tx *gorm.DB, payment *entity.Payment) error {
	panic("")
}

func (repository *PaymentRepositoryMock) Update(payment *entity.Payment) error {
	panic("")
}

var paymentRepository = &PaymentRepositoryMock{Mock: mock.Mock{}}

func TestFindPaymentNotFound(t *testing.T) {
	paymentRepository.Mock.On("FindPaymentByBillNo", "1").Return(nil, mock.Anything)

	payment, err := paymentRepository.FindPaymentByBillNo("1")
	assert.Nil(t, payment)
	assert.NotNil(t, err)
}

func TestFindPaymentFound(t *testing.T) {
	paymentRepository.Mock.On("FindPaymentByBillNo", "1").Return(entity.Payment{
		Email: "adit@mail.com",
	}, nil)

	payment, err := paymentRepository.FindPaymentByBillNo("1")
	assert.NotNil(t, payment)
	assert.Nil(t, err)
}

func TestUpdatePaymentBillNotFound(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		// t.Fatal("error creating mock database: %v",err)
		fmt.Println(err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("error opening GORM database: %v", err)
	}
	paymentRepository.Mock.On("FindPaymentByBillNo", "INV-202108222").Return(nil, errors.New("record not found"))
	var paymentService = service.NewPaymentService(paymentRepository, gormDB)

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
		// t.Fatal("error creating mock database: %v",err)
		fmt.Println(err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("error opening GORM database: %v", err)
	}
	paymentRepository.Mock.On("FindPaymentByBillNo", "INV-202308222").Return(entity.Payment{
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
	var paymentService = service.NewPaymentService(paymentRepository, gormDB)

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
