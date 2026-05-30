package handlers

import (
	"net/http"
	"recruitmentportal/config"
	"recruitmentportal/models"

	"github.com/gin-gonic/gin"
)

func ApproveRecruiter(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := config.DB.Where("id = ?", userID).First(&user); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	if user.Role != models.PendingRecruiter {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User is not a pending recruiter",
		})
		return
	}

	user.Role = models.Recruiter
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
	})
}