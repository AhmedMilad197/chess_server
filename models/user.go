package models

import (
	"time"
)

type User struct {
	ID        uint             `gorm:"primaryKey;autoIncrement"`
	UserName  string           `gorm:"unique;not null"`
	Password  string           `gorm:"not null"`
	Email     string           `gorm:"unique;not null"`
	Ratings   []UserGameRating `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
