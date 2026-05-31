package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"recruitmentportal/models"
	"recruitmentportal/config"
	"recruitmentportal/utils"
)

func RegisterUser(c *gin.Context){
	var u models.CreateUserRequest

	if err := c.ShouldBindJSON(&u); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return 
	}

	var role models.Role

	if u.Role == models.Recruiter {
		role = models.PendingRecruiter
	} else {
		role = models.Applicant
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(u.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Enter a shorter password"})
		return
	}

	user := models.User{
		Name:  u.Name,
		Email: u.Email,
		Password: string(hashedPassword),
		Role: role,
	}

	err = config.DB.QueryRow(
		"INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Name, user.Email, user.Password, user.Role,
	).Scan(&user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	userResponse := models.UserResponse{
		Name: u.Name,
		Email: u.Email,
		Role: user.Role,
		ID: user.ID,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"user": userResponse,
	})
}

func LoginUser(c *gin.Context){
	var u models.CreateUserRequest

	if err := c.ShouldBindJSON(&u); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.User
	row := config.DB.QueryRow("SELECT id, name, email, password, role FROM users WHERE email = $1", u.Email)
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role); err != nil {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(u.Password),
	)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Login successful",
		"token": token,
		"user": gin.H{	
			"id": user.ID,
			"name": user.Name,
			"email": user.Email,
			"role": user.Role,
		},
	})
}