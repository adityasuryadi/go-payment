package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// type Payment struct {
// 	Id                uuid.UUID `gorm:"primaryKey;type:uuid;" column:"id"`
// 	Name              string    `gorm:"column:name"`
// 	Phone             string    `gorm:"column:phone"`
// 	Email             string    `gorm:"column:email"`
// 	BookingDate       time.Time `gorm:"column:booking_date"`
// 	RedirectUrl       string    `gorm:"column:redirect_url"`
// 	ServiceCode       string    `gorm:"column:service_code"`
// 	Qty               int64     `gorm:"column:qty"`
// 	ServiceId         string    `gorm:"column:service_id"`
// 	BillNo            string    `gorm:"column:bill_no"`
// 	BillTotal         float64   `gorm:"column:bill_total"`
// 	StatusId          int       `gorm:"column:status_id"`
// 	BillNoCounter     int       `gorm:"column:bill_no_counter"`
// 	TrxId             string    `gorm:"column:trx_id"`
// 	PaymentChannelUid int       `gorm:"column:payment_channel_uid"`
// 	PaymentChannel    string    `gorm:"column:payment_channel"`
// 	Signature         string    `gorm:"column:signature"`
// 	SnapToken         string    `gorm:"column:snap_token"`
// 	CreatedAt         time.Time `gorm:"column:created_at"`
// 	UpdatedAt         time.Time `gorm:"column:updated_at"`
// }

type Payment struct {
	Id                uuid.UUID `gorm:"primaryKey;type:uuid;" column:"id"`
	Name              string    `gorm:"column:name"`
	Phone             string    `gorm:"column:phone"`
	Email             string    `gorm:"column:email"`
	BookingId         uuid.UUID `gorm:"column:bookingId"`
	BookingDate       time.Time `gorm:"column:bookingDate"`
	RedirectUrl       string    `gorm:"column:redirectUrl"`
	ServiceCode       string    `gorm:"column:uniqueCode"`
	Qty               int64     `gorm:"column:qty"`
	ServiceId         string    `gorm:"column:serviceId"`
	BillNo            string    `gorm:"column:billNo"`
	BillTotal         float64   `gorm:"column:price"`
	StatusId          int       `gorm:"column:status"`
	BillNoCounter     int       `gorm:"column:bill_no_counter"`
	TrxId             string    `gorm:"column:trxId"`
	PaymentChannelUid int       `gorm:"column:paymentChannelCode"`
	PaymentChannel    string    `gorm:"column:paymentChannelName"`
	Signature         string    `gorm:"column:signature"`
	SnapToken         string    `gorm:"column:snapToken"`
	CreatedAt         time.Time `gorm:"column:createdAt"`
	UpdatedAt         time.Time `gorm:"column:updatedAt"`
}

func (Payment) TableName() string {
	return "payment"
}

func (entity *Payment) BeforeCreate(db *gorm.DB) error {
	entity.Id = uuid.New()
	entity.CreatedAt = time.Now().Local()
	return nil
}

func (entity *Payment) BeforeUpdate(db *gorm.DB) error {
	entity.UpdatedAt = time.Now().Local()
	return nil
}
