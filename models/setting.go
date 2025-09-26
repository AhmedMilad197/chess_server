package models

import (
	"time"
)

type Setting struct {
	ID            uint   `gorm:"primaryKey"`
	UserID        uint   `gorm:"uniqueIndex"`
	BoardTheme    string `gorm:"type:varchar(50);default:'standard'"`
	SystemMode    string `gorm:"type:varchar(50);default:'light'"`
	PieceStyle    string `gorm:"type:varchar(50);default:'standard'"`
	Notifications bool   `gorm:"default:true"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
