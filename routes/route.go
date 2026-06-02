package routes

import (
	"recruitmentportal/handlers"
	"recruitmentportal/middleware"
    "recruitmentportal/models"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Routes working",
        })
    })

    auth := r.Group("/auth")
    {
        auth.POST("/register", handlers.RegisterUser)
        auth.POST("/login", handlers.LoginUser)
    }

    protected := r.Group("/")
    protected.Use(middleware.AuthMiddleware())
    {
        protected.GET("/profile", handlers.Profile)
    }

    jobs := r.Group("/jobs")
    {
        jobs.GET("", handlers.GetJobs)
    }

    recruiter := r.Group("/recruiter")
    recruiter.Use(
        middleware.AuthMiddleware(),
        middleware.RequireRole(
            string(models.Recruiter),
            string(models.Admin),
        ),
    )
    {
        recruiter.POST("/jobs", handlers.CreateJob)
        recruiter.GET("/jobs", handlers.GetRecruiterJobs)
        recruiter.PUT("/jobs/:id", handlers.UpdateJob)
        recruiter.DELETE("/jobs/:id", handlers.DeleteJob)
        recruiter.GET("/jobs/:id/applicants", handlers.GetJobApplicants)
    }

    applicant := r.Group("/applicant")
    applicant.Use(
        middleware.AuthMiddleware(),
        middleware.RequireRole(
            string(models.Applicant),
        ),
    )
    {
        applicant.POST("/jobs/:id/apply", handlers.ApplyJob)
    }

    admin := r.Group("/admin")
    admin.Use(
        middleware.AuthMiddleware(),
        middleware.RequireRole(string(models.Admin)),
    )
    {
        admin.PUT("/approve/:id", handlers.ApproveRecruiter)
    }
}