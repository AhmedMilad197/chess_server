package routes

import (
	"chess_server/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(api *gin.RouterGroup) {
	api.POST("/register", controllers.Register)
	api.POST("/login", controllers.Login)
}
