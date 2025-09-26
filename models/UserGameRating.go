package models

import (
	"time"
)

type UserGameRating struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint `gorm:"index;not null"`
	GameTypeID uint `gorm:"index;not null"`
	Rating     int  `gorm:"default:1200"`

	User      User     `gorm:"foreignKey:UserID"`
	GameType  GameType `gorm:"foreignKey:GameTypeID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
