package main

import (
	"chess_server/config"
	"chess_server/database"
	"chess_server/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	config.LoadConfig()
	db.Init()
	router := gin.Default()
	routes.ApiRoutes(router)
	router.Run()
}
