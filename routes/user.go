package routes

import (
	"chess_server/controllers"
	"github.com/gin-gonic/gin"
	"chess_server/middleware"
)

func UserRoutes(api *gin.RouterGroup) {
	userRoutes := api.Group("/users")
	userRoutes.Use(middleware.Auth())
	userRoutes.GET("/settings", controllers.GetSettings)
}
