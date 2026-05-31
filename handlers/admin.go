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
	row := config.DB.QueryRow("SELECT id, name, email, password, role FROM users WHERE id = $1", userID)
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role); err != nil {
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
	_, err := config.DB.Exec("UPDATE users SET role = $1 WHERE id = $2", user.Role, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
	})
}