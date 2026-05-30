package models

type Job struct{
	ID          uint	`gorm:"primaryKey"`
	Title       string	`gorm:"not null"`
	Description string	`gorm:"not null"`
	Company     string	`gorm:"not null;index"`
	CompanyDescription string	
	CompanyContactMail string	`gorm:"not null"`
	CreatedBy   uint	`gorm:"index"`

	Skills		[]Skill `gorm:"many2many:job_skills;"`
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