package handlers

import (
	"net/http"
	"recruitmentportal/config"
	"recruitmentportal/models"

	"github.com/gin-gonic/gin"
)

func CreateJob(c *gin.Context){
	var j models.CreateJobRequest
	var skills []models.Skill
	id, _ := c.Get("user_id")
	userID := uint(id.(float64))

	if err := c.ShouldBindJSON(&j); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	for _, s := range(j.Skills){
		var skill models.Skill
		if err := config.DB.Where("name = ?", s).First(&skill).Error; err != nil{
			skill = models.Skill{Name: s}
			if err := config.DB.Create(&skill).Error; err != nil{
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to create skill",
				})
				return
			}
		}
		skills = append(skills, skill)
	}

	job := models.Job{
		Title: j.Title,
		Description: j.Description,
		Company: j.Company,
		CompanyDescription: j.CompanyDescription,
		CompanyContactMail: j.CompanyContactMail,
		CreatedBy: userID,
		Skills: skills,
	}

	if err := config.DB.Create(&job).Error; err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create job",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Job created successfully",
		"job_id": job.ID,
	})
}