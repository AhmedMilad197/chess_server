package routes

import (
	"chess_server/controllers"
	"github.com/gin-gonic/gin"
	"chess_server/middlewares"
)

func UserRoutes(api *gin.RouterGroup) {
	userRoutes := api.Group("/users")
	userRoutes.Use(middleware.Auth())
	userRoutes.GET("/settings", controllers.GetSettings)
	userRoutes.GET("/info", controllers.GetInfo)
}
