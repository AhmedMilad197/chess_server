package routes

import (
	"chess_server/controllers"
	"chess_server/middlewares"
	"github.com/gin-gonic/gin"
)

func WSRoutes(router *gin.Engine) {
	wsRoutes := router.Group("/ws")
	wsRoutes.Use(middleware.Auth())
	wsRoutes.GET("/", controllers.ConnectWS)
}
