package main

import (
	"chess_server/config"
	"chess_server/database"
	"chess_server/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

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
	routes.ApiRoutes(router)
	routes.WSRoutes(router)
	router.Run()
}
