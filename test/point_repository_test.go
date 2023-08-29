package test

import (
	"payment/config"
	"payment/entity"
	"payment/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInsertPointIfUserNotExist(t *testing.T) {
	now := time.Now()
	configApp := config.New(`\.env.test`)
	db := config.NewPostgresDB(configApp)
	point := &entity.Point{
		UserId:    30,
		Point:     2500,
		CreatedAt: now,
	}
	// db,mock,err := sqlmock.New()
	// if err != nil {
	// 	t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	// }
	repository := repository.NewPointRepository(db)
	result, err := repository.InsertOrUpdate(point)
	assert.Nil(t, err)
	assert.Equal(t, 30, result.UserId)
	assert.Equal(t, float64(2500), result.Point)
}

func TestInsertPointIfUserExist(t *testing.T) {
	now := time.Now()
	const point_user = 2500
	configApp := config.New(`\.env.test`)
	db := config.NewPostgresDB(configApp)
	repository := repository.NewPointRepository(db)

	resultPoint, _ := repository.FindPointByUserId(30)

	point := &entity.Point{
		UserId:    30,
		Point:     resultPoint.Point + point_user,
		CreatedAt: now,
	}
	// db,mock,err := sqlmock.New()
	// if err != nil {
	// 	t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	// }
	result, err := repository.InsertOrUpdate(point)
	assert.Nil(t, err)
	assert.Equal(t, 30, result.UserId)
	assert.Equal(t, resultPoint.Point+point_user, result.Point)
}

func TestFindPointByUserIdNotExist(t *testing.T) {
	configApp := config.New(`\.env.test`)
	db := config.NewPostgresDB(configApp)
	repository := repository.NewPointRepository(db)
	user_id := 10000
	point, err := repository.FindPointByUserId(user_id)
	assert.Nil(t, point)
	assert.NotNil(t, err)
}

func TestFindPointByUserIsExist(t *testing.T) {
	configApp := config.New(`\.env.test`)
	db := config.NewPostgresDB(configApp)
	repository := repository.NewPointRepository(db)
	user_id := 30
	point, err := repository.FindPointByUserId(user_id)
	assert.Nil(t, err)
	assert.NotNil(t, point)
}
