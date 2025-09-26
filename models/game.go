package models

import (
	"time"
)

type Game struct {
	ID         uint     `gorm:"primaryKey;autoIncrement"`
	Player1ID  uint     `gorm:"not null"`
	Player2ID  uint     `gorm:"not null"`
	Player1    User     `gorm:"foreignKey:Player1ID"`
	Player2    User     `gorm:"foreignKey:Player2ID"`
	GameTypeID uint     `gorm:"not null"`
	GameType   GameType `gorm:"foreignKey:GameTypeID"`
	Status     string   `gorm:"type:varchar(50);default:'pending'"`
	WinnerID   *uint    `gorm:"default:null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
