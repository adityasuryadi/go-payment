package entity

import (
	"time"

	"gorm.io/gorm"
)

// type Point struct {
// 	Id        uuid.UUID `gorm:"primaryKey;type:uuid;" column:"id"`
// 	UserId    int       `gorm:"column:user_id"`
// 	Point     float64   `gorm:"column:point"`
// 	CreatedAt time.Time `gorm:"column:created_at"`
// 	UpdatedAt time.Time `gorm:"column:updated_at"`
// }

type Point struct {
	Id        int       `gorm:"primaryKey;type:int;column:token_id"`
	UserId    int       `gorm:"column:customer_id"`
	Point     float64   `gorm:"column:total_token"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_date"`
}

func (Point) TableName() string {
	return "token"
}

func (entity *Point) BeforeCreate(db *gorm.DB) error {
	entity.CreatedAt = time.Now().Local()
	return nil
}

func (entity *Point) BeforeUpdate(db *gorm.DB) error {
	entity.UpdatedAt = time.Now().Local()
	return nil
}
