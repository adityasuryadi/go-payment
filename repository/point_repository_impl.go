package repository

import (
	"payment/entity"

	"gorm.io/gorm"
)

func NewPointRepository(db *gorm.DB) PointRespository {
	return &PointRespositoryImpl{
		db: db,
	}
}

type PointRespositoryImpl struct {
	db *gorm.DB
}

// FindPointByUserId implements PointRespository
func (repository *PointRespositoryImpl) FindPointByUserId(user_id int) (*entity.Point, error) {
	var point *entity.Point
	err := repository.db.Where("user_id = ?", user_id).First(&point).Error
	if err != nil {
		return nil, err
	}
	return point, nil
}

// InsertOrUpdate implements PointRespository
func (repository *PointRespositoryImpl) InsertOrUpdate(point *entity.Point) (*entity.Point, error) {
	err := repository.db.Debug().Where(entity.Point{UserId: point.UserId}).Assign(entity.Point{Point: point.Point}).FirstOrCreate(&point).Error
	if err != nil {
		return nil, err
	}
	return point, nil
}
