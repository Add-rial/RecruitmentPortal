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

    r.POST("/register", handlers.RegisterUser)
    r.POST("/login", handlers.LoginUser)
    r.GET("/profile", middleware.AuthMiddleware(),handlers.Profile)

    r.POST("/jobs", middleware.AuthMiddleware(), middleware.RequireRole(string(models.Recruiter), string(models.Admin)), handlers.CreateJob)

    r.PUT("/admin/approve/:id", middleware.AuthMiddleware(), middleware.RequireRole("admin"),handlers.ApproveRecruiter,)
}