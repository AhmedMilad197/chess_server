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

func GetInfo(c *gin.Context) {
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

	type GameRatingResponse struct {
		GameType string `json:"game_type"`
		Duration uint   `json:"duration"`
		Rating   int    `json:"rating"`
	}

	type SettingResponse struct {
		BoardTheme    string `json:"board_theme"`
		SystemMode    string `json:"system_mode"`
		PieceStyle    string `json:"piece_style"`
		Notifications bool   `json:"notifications"`
	}

	type UserInfoResponse struct {
		ID       uint                 `json:"id"`
		UserName string               `json:"username"`
		Email    string               `json:"email"`
		Settings SettingResponse      `json:"settings"`
		Ratings  []GameRatingResponse `json:"ratings"`
	}

	var dbUser models.User
	if err := db.DB.
		Preload("Setting").
		Preload("Ratings.GameType").
		First(&dbUser, user.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User not found"})
		return
	}

	resp := UserInfoResponse{
		ID:       dbUser.ID,
		UserName: dbUser.UserName,
		Email:    dbUser.Email,
		Settings: SettingResponse{
			BoardTheme:    dbUser.Setting.BoardTheme,
			SystemMode:    dbUser.Setting.SystemMode,
			PieceStyle:    dbUser.Setting.PieceStyle,
			Notifications: dbUser.Setting.Notifications,
		},
	}

	for _, r := range dbUser.Ratings {
		resp.Ratings = append(resp.Ratings, GameRatingResponse{
			GameType: r.GameType.Name,
			Duration: r.GameType.Duration,
			Rating:   r.Rating,
		})
	}

	c.JSON(http.StatusOK, resp)
}
