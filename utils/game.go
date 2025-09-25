package utils

import (
	"log"
	"time"

	"chess_server/models"
	"gorm.io/gorm"
)

func SeedGameTypes(db *gorm.DB) {
	gameTypes := []models.GameType{
		{Name: "Bullet", Duration: 1},
		{Name: "Blitz", Duration: 5},
		{Name: "Rapid", Duration: 10},
		{Name: "Classical", Duration: 30},
	}

	for _, gt := range gameTypes {
		var existing models.GameType
		if err := db.Where("name = ?", gt.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				gt.CreatedAt = time.Now()
				gt.UpdatedAt = time.Now()
				if err := db.Create(&gt).Error; err != nil {
					log.Printf("Failed to create game type %s: %v", gt.Name, err)
				}
			}
		}
	}
}
