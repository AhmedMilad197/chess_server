package routes

import (
	"chess_server/controllers"
	"chess_server/middlewares"
	"github.com/gin-gonic/gin"
)

func GameRoutes(api *gin.RouterGroup) {
	gameRoutes := api.Group("/games")
	gameRoutes.Use(middleware.Auth())
	gameRoutes.GET("/", controllers.GetGameTypes)
	gameRoutes.GET("/:id/search", controllers.SearchGame)
}
