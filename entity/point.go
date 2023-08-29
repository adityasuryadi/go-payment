package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Point struct {
	Id        uuid.UUID `gorm:"primaryKey;type:uuid;" column:"id"`
	UserId    int       `gorm:"column:user_id"`
	Point     float64   `gorm:"column:point"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Point) TableName() string {
	return "points"
}

func (entity *Point) BeforeCreate(db *gorm.DB) error {
	entity.Id = uuid.New()
	entity.CreatedAt = time.Now().Local()
	return nil
}

func (entity *Point) BeforeUpdate(db *gorm.DB) error {
	entity.UpdatedAt = time.Now().Local()
	return nil
}
