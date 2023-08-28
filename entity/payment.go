package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	Id                uuid.UUID `gorm:"primaryKey;type:uuid;" column:"id"`
	Name              string    `gorm:"column:name"`
	Phone             string    `gorm:"column:phone"`
	Email             string    `gorm:"column:email"`
	BookingDate       time.Time `gorm:"column:booking_date"`
	RedirectUrl       string    `gorm:"column:redirect_url"`
	Qty               int64     `gorm:"column:qty"`
	ServiceId         string    `gorm:"column:service_id"`
	BillNo            string    `gorm:"column:bill_no"`
	BillTotal         float64   `gorm:"column:bill_total"`
	StatusId          int       `gorm:"column:status_id"`
	BillNoCounter     int       `gorm:"column:bill_no_counter"`
	TrxId             string    `gorm:"column:trx_id"`
	PaymentChannelUid int       `gorm:"column:payment_channel_uid"`
	PaymentChannel    string    `gorm:"column:payment_channel"`
	Signature         string    `gorm:"column:signature"`
	CreatedAt         time.Time `gorm:"column:created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at"`
}

func (Payment) TableName() string {
	return "payments"
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
