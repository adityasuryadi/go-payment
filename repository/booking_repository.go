package repository

import (
	"payment/entity"

	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(tx *gorm.DB, booking *entity.Booking) (*entity.Booking, error)
}
