package repository

import (
	"payment/entity"

	"gorm.io/gorm"
)

type BookingRepositoryImpl struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &BookingRepositoryImpl{
		db: db,
	}
}

// Create implements BookingRepository.
func (repository *BookingRepositoryImpl) Create(booking *entity.Booking) (*entity.Booking, error) {
	err := repository.db.Create(&booking).Error
	if err != nil {
		return nil, err
	}
	return booking, nil
}
