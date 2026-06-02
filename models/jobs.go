package models

type Job struct {
	ID          uint     `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`

	Company string `json:"company"`

	CompanyDescription string `json:"company_description"`
	CompanyContactMail string `json:"company_contact_mail"`

	CreatedBy uint `json:"created_by"`

	Skills []Skill `json:"skills"`
}

type Skill struct{
	ID		uint	`gorm:"primaryKey"`
	Name	string	`gorm:"unique;not null"`
}

type CreateJobRequest struct {
	Title                 string   `json:"title" binding:"required"`
	Description           string   `json:"description" binding:"required"`
	Company               string   `json:"company" binding:"required"`
	CompanyDescription    string   `json:"company_description"`
	CompanyContactMail    string   `json:"company_contact_mail" binding:"required"`
	Skills                []string `json:"skills" binding:"required"`
}

type Application struct {
	ID uint
	UserID uint
	JobID uint
}