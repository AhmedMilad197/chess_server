package controllers

import (
	"chess_server/database"
	"chess_server/models"
	"chess_server/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var payload struct {
		Email    string `form:"email" json:"email" xml:"email"  binding:"required,email"`
		UserName string `form:"username" json:"username" xml:"username"  binding:"required"`
		Password string `form:"password" json:"password" xml:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var existingUser models.User
	err := db.DB.Where("user_name = ? OR email = ?", payload.UserName, payload.Email).
		First(&existingUser).Error

	if err == nil {
		if existingUser.UserName == payload.UserName {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}
		if existingUser.Email == payload.Email {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
	}
	hash, hashErr := utils.GenerateFromPassword(payload.Password)
	if hashErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong please try again later"})
		return
	}
	newUser := models.User{
		UserName: payload.UserName,
		Email:    payload.Email,
		Password: hash,
	}
	if err := db.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	token, createTokenErr := utils.CreateToken(
		map[string]interface{}{
			"id":       newUser.ID,
			"username": newUser.UserName,
		},
		time.Hour*24*30,
	)
	if createTokenErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong please try again later"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Registered successfully",
		"token":   token,
	})
}

func Login(c *gin.Context) {
	var payload struct {
		UserName string `form:"username" json:"username" xml:"username"  binding:"required"`
		Password string `form:"password" json:"password" xml:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := db.DB.Where("user_name = ?", payload.UserName).First(&user).Error

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "There is no user with this username"})
		return
	}
	match, compareErr := utils.ComparePasswordAndHash(payload.Password, user.Password)
	if compareErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong please try again later"})
		return
	}
	if !match {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	token, createTokenErr := utils.CreateToken(
		map[string]interface{}{
			"id":       user.ID,
			"username": user.UserName,
		},
		time.Hour*24*30,
	)
	if createTokenErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong please try again later"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in successfully",
		"token":   token,
	})
}
