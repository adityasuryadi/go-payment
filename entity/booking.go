package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Booking struct {
	Id          uuid.UUID `gorm:"primaryKey;type:uuid;" column:"id"`
	Name        string    `gorm:"column:name"`
	Phone       string    `gorm:"column:phone"`
	Email       string    `gorm:"column:email"`
	BookingDate time.Time `gorm:"column:bookingDate"`
	ServiceCode string    `gorm:"column:uniqueCode"`
	ServiceId   string    `gorm:"column:serviceId"`
	StatusId    int       `gorm:"column:status"`
	CreatedAt   time.Time `gorm:"column:createdAt"`
	UpdatedAt   time.Time `gorm:"column:updatedAt"`
}

func (Booking) TableName() string {
	return "booking"
}

func (entity *Booking) BeforeCreate(db *gorm.DB) error {
	entity.Id = uuid.New()
	entity.CreatedAt = time.Now().Local()
	return nil
}

func (entity *Booking) BeforeUpdate(db *gorm.DB) error {
	entity.UpdatedAt = time.Now().Local()
	return nil
}
