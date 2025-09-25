package models

import (
	"time"
)

type GameType struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"not null"`
	Duration  uint   `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
