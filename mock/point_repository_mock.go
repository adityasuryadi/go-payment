package mocks

import (
	"errors"
	"payment/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type PointRepositoryMock struct {
	Mock mock.Mock
}

func (repository *PointRepositoryMock) FindPointByUserId(user_id int) (*entity.Point, error) {
	arguments := repository.Mock.Called(user_id)
	if arguments.Get(0) == nil {
		return nil, gorm.ErrRecordNotFound
	} else {
		point := arguments.Get(0).(entity.Point)
		return &point, nil
	}
}

func (repository *PointRepositoryMock) InsertOrUpdate(point *entity.Point) (*entity.Point, error) {
	arguments := repository.Mock.Called(point)
	if arguments.Get(0) == nil {
		return nil, errors.New("error")
	} else {
		point := arguments.Get(0).(entity.Point)
		return &point, nil
	}
}
