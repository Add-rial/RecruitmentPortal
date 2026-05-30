package config

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/joho/godotenv"

	"recruitmentportal/models"
)

var (
	DB *gorm.DB
	JWTSecret []byte
)

func init(){
	godotenv.Load()
	JWTSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
}

func ConnectDB(){
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to Postgres:", err)
	}

	err = DB.AutoMigrate(&models.User{}, &models.Job{}, &models.Skill{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database schema: %v", err)
	}

	log.Println("Database connection and migration successful.")
}