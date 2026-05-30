package handlers

import (
	"net/http"
	"recruitmentportal/config"
	"recruitmentportal/models"

	"github.com/gin-gonic/gin"
)

func Profile(c *gin.Context){
	userID, _ := c.Get("user_id")
	id := uint(userID.(float64))

	var u models.User

	if err := config.DB.Where("id = ?", id).First(&u).Error; err != nil{
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found, create a new account",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": u,
	})
}