package controllers

import (
	"chess_server/database"
	"chess_server/models"
	"errors"
	"gorm.io/gorm"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSettings(c *gin.Context) {
	userData, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Could not get the user"})
		return
	}

	user, ok := userData.(models.User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Could not get the user"})
		return
	}

	var settings models.Setting
	if err := db.DB.Where("user_id = ?", user.ID).First(&settings).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Settings not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user settings",
		"data": gin.H{
			"BoardTheme":    settings.BoardTheme,
			"SystemMode":    settings.SystemMode,
			"PieceStyle":    settings.PieceStyle,
			"Notifications": settings.Notifications,
		},
	})
}
