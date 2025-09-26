package models

import (
	"time"
)

type Setting struct {
	ID                         uint   `gorm:"primaryKey"`
	UserID                     uint   `gorm:"uniqueIndex"`
	BoardTheme                 string `gorm:"type:varchar(50);default:'standard'"`
	SystemMode                 string `gorm:"type:varchar(50);default:'light'"`
	PieceStyle                 string `gorm:"type:varchar(50);default:'standard'"`
	Notifications              bool   `gorm:"default:true"`
	LowerBoundPlayerRatingDiff uint   `gorm:"default:100"`
	UpperBoundPlayerRatingDiff uint   `gorm:"default:100"`
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
}
