package handlers

import (
	"database/sql"
	"net/http"
	"recruitmentportal/config"
	"recruitmentportal/models"

	"github.com/gin-gonic/gin"
)

func Profile(c *gin.Context){
	userID, _ := c.Get("user_id")
	id := uint(userID.(float64))

	var u models.User

	row := config.DB.QueryRow("SELECT id, name, email, password, role FROM users WHERE id = $1", id)
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found, create a new account",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Database error",
			})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": u,
	})
}