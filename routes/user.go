package routes

import (
	"chess_server/controllers"
	"chess_server/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(api *gin.RouterGroup) {
	userRoutes := api.Group("/users")
	userRoutes.Use(middleware.Auth())
	userRoutes.GET("/", controllers.ListUsers)
	userRoutes.GET("/settings", controllers.GetSettings)
	userRoutes.GET("/info", controllers.GetInfo)
	userRoutes.PUT("/update", controllers.UpdateUser)
	userRoutes.PUT("/settings/update", controllers.UpdateUserSettings)
}
