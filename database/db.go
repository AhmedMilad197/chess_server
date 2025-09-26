package db

import (
	"fmt"
	"log"

	"chess_server/config"
	"chess_server/models"
	"chess_server/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	host := config.Config.DBHost
	user := config.Config.DBUser
	password := config.Config.DBPassword
	dbname := config.Config.DBName
	port := config.Config.DBPort
	sslmode := config.Config.DBSSLMode

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode,
	)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(0)

	DB.AutoMigrate(
		&models.User{},
		&models.Setting{},
		&models.GameType{},
		&models.UserGameRating{},
	)

	utils.SeedGameTypes(DB)
	fmt.Println("Database connection established")
}
