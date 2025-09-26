package controllers

import (
	"chess_server/database"
	"chess_server/models"
	"net/http"

	"chess_server/utils"
	"github.com/gin-gonic/gin"
	"strconv"
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

func SearchGame(c *gin.Context) {
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

	gameTypeId, _ := strconv.Atoi(c.Param("id"))
	utils.EnqueuePlayer(user.ID, gameTypeId)

	c.JSON(http.StatusOK, gin.H{
		"message": "Player added to matchmaking queue",
	})
}
