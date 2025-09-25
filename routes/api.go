package routes

import (
	"github.com/gin-gonic/gin"
)

func ApiRoutes(router *gin.Engine) {
	api := router.Group("/api")
	AuthRoutes(api)
	UserRoutes(api)
}
