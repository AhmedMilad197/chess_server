package routes

import (
	"chess_server/controllers"
	// "chess_server/middlewares"
	"github.com/gin-gonic/gin"
)

func GameRoutes(api *gin.RouterGroup) {
	gameRoutes := api.Group("/games")
	gameRoutes.GET("/", controllers.GetGameTypes)
	gameRoutes.GET("/:id/play", controllers.PlayGame)
}
