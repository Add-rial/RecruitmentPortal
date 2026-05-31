package handlers

import (
	"database/sql"
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

	for _, s := range j.Skills {
		var skill models.Skill
		row := config.DB.QueryRow("SELECT id, name FROM skills WHERE name = $1", s)
		err := row.Scan(&skill.ID, &skill.Name)
		if err != nil {
			if err == sql.ErrNoRows {
				skill.Name = s
				err = config.DB.QueryRow("INSERT INTO skills (name) VALUES ($1) RETURNING id", s).Scan(&skill.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "Failed to create skill",
					})
					return
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Database error checking skill",
				})
				return
			}
		}
		skills = append(skills, skill)
	}

	job := models.Job{
		Title:              j.Title,
		Description:        j.Description,
		Company:            j.Company,
		CompanyDescription: j.CompanyDescription,
		CompanyContactMail: j.CompanyContactMail,
		CreatedBy:          userID,
		Skills:             skills,
	}

	tx, err := config.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to start database transaction",
		})
		return
	}
	defer tx.Rollback()

	err = tx.QueryRow(
		"INSERT INTO jobs (title, description, company, company_description, company_contact_mail, created_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		job.Title, job.Description, job.Company, job.CompanyDescription, job.CompanyContactMail, job.CreatedBy,
	).Scan(&job.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create job",
		})
		return
	}

	for _, skill := range skills {
		_, err = tx.Exec("INSERT INTO job_skills (job_id, skill_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", job.ID, skill.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to map job to skills",
			})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to commit transaction",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Job created successfully",
		"job_id": job.ID,
	})
}