package models

type Role string

const (
    Applicant Role = "applicant"
    Recruiter Role = "recruiter"
    Admin     Role = "admin"
    PendingRecruiter    Role = "pending_recruiter"
)

type CreateUserRequest struct {
    Name     string `json:"name"`
    Email    string `json:"email" binding:"required"`
    Password string `json:"password" binding:"required"`
    Role     Role   `json:"role"`
}

type UserResponse struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Role  Role   `json:"role"`
}

type User struct {
    ID       int    `gorm:"primaryKey" json:"id"`
    Name     string `json:"name"`
    Email    string `gorm:"unique;not null" json:"email"`
    Password string `json:"-"`
    Role     Role   `json:"role"`
}

type LoginRequest struct {
    Email string `json:"email" binding:"required"`
    Password string `json:"password" binding:"required"`
}