package repository

import "payment/entity"

type PointRespository interface {
	InsertOrUpdate(point *entity.Point) (*entity.Point, error)
	FindPointByUserId(user_id int) (*entity.Point, error)
}
