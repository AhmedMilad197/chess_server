package controllers

import (
	"chess_server/database"
	"chess_server/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetGameTypes(c *gin.Context) {
	var gameTypes []models.GameType
	result := db.DB.Find(&gameTypes)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving game types"})
		return
	}

	type GameTypeResponse struct {
		Name     string `json:"name"`
		Duration uint   `json:"duration"`
	}

	responses := make([]GameTypeResponse, 0, len(gameTypes))
	for _, gt := range gameTypes {
		responses = append(responses, GameTypeResponse{
			Name:     gt.Name,
			Duration: gt.Duration,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "game types",
		"data":    responses,
	})
}
