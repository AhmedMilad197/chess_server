package routes

import (
	"chess_server/controllers"
	"github.com/gin-gonic/gin"
	"chess_server/middlewares"
)

func GameRoutes(api *gin.RouterGroup) {
	gameRoutes := api.Group("/games")
	gameRoutes.Use(middleware.Auth())
	gameRoutes.GET("/types", controllers.GetGameTypes)
}
