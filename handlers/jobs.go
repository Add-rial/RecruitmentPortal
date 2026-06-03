package handlers

import (
	"database/sql"
	"net/http"
	"recruitmentportal/config"
	"recruitmentportal/models"
	"recruitmentportal/utils"
	"strconv"

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

func GetJobs(c *gin.Context){
	rows, err := config.DB.Query(`
		SELECT
			j.id,
			j.title,
			j.description,
			j.company,
			j.company_description,
			j.company_contact_mail,
			j.created_by,
			s.id,
			s.name
		FROM jobs j
		LEFT JOIN job_skills js
			ON j.id = js.job_id
		LEFT JOIN skills s
			ON s.id = js.skill_id
		ORDER BY j.id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch jobs",
		})
		return
	}
	defer rows.Close()

	jobs, err := utils.JobsJSONCreator(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to scan jobs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs": jobs,
	})
}

func GetRecruiterJobs(c * gin.Context){
	id, _ := c.Get("user_id")
	userID := uint(id.(float64))

	rows, err := config.DB.Query(`
		SELECT
			j.id,
			j.title,
			j.description,
			j.company,
			j.company_description,
			j.company_contact_mail,
			j.created_by,
			s.id,
			s.name
		FROM jobs j
		LEFT JOIN job_skills js
			ON j.id = js.job_id
		LEFT JOIN skills s
			ON s.id = js.skill_id
		WHERE j.created_by = $1
		ORDER BY j.id`, userID)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Couldn't fetch jobs",
		})
		return
	}
	defer rows.Close()
	
	jobs, err := utils.JobsJSONCreator(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to scan jobs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs": jobs,
	})
}

func ApplyJob(c *gin.Context){
	jobIDstr := c.Param("id")
	jobID, err := strconv.Atoi(jobIDstr)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Couldn't fetch job id",
		})
		return
	}

	id, _ := c.Get("user_id")
	userID := uint(id.(float64))

	_, err = config.DB.Exec(`
		INSERT INTO applications (user_id, job_id) values($1, $2)
	`, userID, jobID)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Couldn't apply for job",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully applied for job",
	})
}

func UpdateJob(c *gin.Context) {
	id, _ := c.Get("user_id")
	userID := uint(id.(float64))

	jobIDstr := c.Param("id")
	jobID, err := strconv.Atoi(jobIDstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid job ID",
		})
		return
	}

	var createdBy uint
	err = config.DB.QueryRow("SELECT created_by FROM jobs WHERE id = $1", jobID).Scan(&createdBy)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Job not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error checking ownership",
		})
		return
	}

	if createdBy != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not authorized to update this job",
		})
		return
	}

	var req models.CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var skills []models.Skill
	for _, s := range req.Skills {
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

	tx, err := config.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to start database transaction",
		})
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE jobs
		SET title = $1, description = $2, company = $3, company_description = $4, company_contact_mail = $5
		WHERE id = $6`,
		req.Title, req.Description, req.Company, req.CompanyDescription, req.CompanyContactMail, jobID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update job details",
		})
		return
	}

	_, err = tx.Exec("DELETE FROM job_skills WHERE job_id = $1", jobID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to clear old skills mapping",
		})
		return
	}

	for _, skill := range skills {
		_, err = tx.Exec("INSERT INTO job_skills (job_id, skill_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", jobID, skill.ID)
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
		"message": "Job updated successfully",
	})
}

func DeleteJob(c *gin.Context) {
	id, _ := c.Get("user_id")
	userID := uint(id.(float64))

	jobIDstr := c.Param("id")
	jobID, err := strconv.Atoi(jobIDstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid job ID",
		})
		return
	}

	var createdBy uint
	err = config.DB.QueryRow("SELECT created_by FROM jobs WHERE id = $1", jobID).Scan(&createdBy)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Job not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error checking ownership",
		})
		return
	}

	if createdBy != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not authorized to delete this job",
		})
		return
	}

	_, err = config.DB.Exec("DELETE FROM jobs WHERE id = $1", jobID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete job",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Job deleted successfully",
	})
}

func GetJobApplicants(c *gin.Context) {
	id, _ := c.Get("user_id")
	userID := uint(id.(float64))

	jobIDstr := c.Param("id")
	jobID, err := strconv.Atoi(jobIDstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid job ID",
		})
		return
	}

	var createdBy uint
	err = config.DB.QueryRow("SELECT created_by FROM jobs WHERE id = $1", jobID).Scan(&createdBy)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Job not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error checking ownership",
		})
		return
	}

	if createdBy != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not authorized to view applicants for this job",
		})
		return
	}

	rows, err := config.DB.Query(`
		SELECT u.id, u.name, u.email, u.role, u.resume_url
		FROM users u
		JOIN applications a ON u.id = a.user_id
		WHERE a.job_id = $1`, jobID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch applicants",
		})
		return
	}
	defer rows.Close()

	var applicants []models.UserResponse
	for rows.Next() {
		var app models.UserResponse
		err := rows.Scan(&app.ID, &app.Name, &app.Email, &app.Role, &app.ResumeURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to scan applicant",
			})
			return
		}
		applicants = append(applicants, app)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error reading applicants",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"applicants": applicants,
	})
}


