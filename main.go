package main

import (
	"chess_server/database"
	"chess_server/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"chess_server/config"
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
	api := router.Group("/api")
	routes.AuthRoutes(api)
	router.Run()
}
