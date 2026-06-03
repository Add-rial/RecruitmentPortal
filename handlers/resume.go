package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"recruitmentportal/config"
)

func UploadResume(c *gin.Context) {
	id, _ := c.Get("user_id")
	userID := uint(id.(float64))

	file, err := c.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No file uploaded",
		})
		return
	}

	if file.Size > 1*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File size exceeds 1MB limit",
		})
		return
	}

	if filepath.Ext(file.Filename) != ".pdf" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Only PDF files are allowed",
		})
		return
	}

	uploadDir := "./uploads/resumes"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create upload directory",
		})
		return
	}

	filename := fmt.Sprintf("%d_%d.pdf", userID, time.Now().Unix())
	filePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save file",
		})
		return
	}

	resumeURL := fmt.Sprintf("/applicant/resumes/%s", filename)
	_, err = config.DB.Exec("UPDATE users SET resume_url = $1 WHERE id = $2", resumeURL, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error updating resume URL",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Resume uploaded successfully",
		"resume_url": resumeURL,
	})
}

func ViewResume(c *gin.Context) {
	id, _ := c.Get("user_id")
	userID := uint(id.(float64))

	roleVal, _ := c.Get("role")
	role := roleVal.(string)

	filename := c.Param("filename")

	parts := strings.Split(filename, "_")
	if len(parts) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid filename format",
		})
		return
	}
	ownerIDStr := parts[0]
	ownerID, err := strconv.Atoi(ownerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid filename format",
		})
		return
	}

	if role == "applicant" && uint(ownerID) != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not authorized to view this resume",
		})
		return
	}

	filePath := filepath.Join("./uploads/resumes", filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Resume file not found",
		})
		return
	}

	c.File(filePath)
}

