package controllers

import (
	"chess_server/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSettings(c *gin.Context) {
	userData, ok := c.Get("user");
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"messaage": "Could not get the user"})
		return
	}
	user, ok := userData.(models.User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"messaage": "Could not get the user"})
		return
	}
	println(user.Email)
	
	c.JSON(http.StatusOK, gin.H{"messaage": "user settings"})
}
