package main

import (
	"chess_server/config"
	"chess_server/database"
	"chess_server/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"time"

	"chess_server/utils"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	config.LoadConfig()
	db.Init()
	utils.InitGame()
	utils.InitRedis()
	go utils.MatchmakingWorker()
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.Config.Client},
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.ApiRoutes(router)
	router.Run()
}
