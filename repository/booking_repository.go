package repository

import (
	"payment/entity"
)


type BookingRepository interface {
	Create(booking *entity.Booking) (*entity.Booking, error)
}
