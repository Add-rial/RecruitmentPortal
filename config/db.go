package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

var (
	DB        *sql.DB
	JWTSecret []byte
)

func init() {
	godotenv.Load()
	JWTSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
}

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	var err error
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Failed to open Postgres connection:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to connect to Postgres:", err)
	}

	// Initialize tables if they do not exist
	err = initSchema()
	if err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	log.Println("Database connection and schema initialization successful.")
}

func initSchema() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL,
			resume_url TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS jobs (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT NOT NULL,
			company VARCHAR(255) NOT NULL,
			company_description TEXT,
			company_contact_mail VARCHAR(255) NOT NULL,
			created_by INT
		);`,
		`CREATE INDEX IF NOT EXISTS idx_jobs_company ON jobs(company);`,
		`CREATE INDEX IF NOT EXISTS idx_jobs_created_by ON jobs(created_by);`,
		`CREATE TABLE IF NOT EXISTS skills (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) UNIQUE NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS job_skills (
			job_id INT REFERENCES jobs(id) ON DELETE CASCADE,
			skill_id INT REFERENCES skills(id) ON DELETE CASCADE,
			PRIMARY KEY (job_id, skill_id)
		);`,
		`CREATE TABLE IF NOT EXISTS applications (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			job_id INT NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(user_id, job_id)
		);`,
		`CREATE TABLE messages (
			id SERIAL PRIMARY KEY,
			sender_id INT NOT NULL REFERENCES users(id),
			receiver_id INT NOT NULL REFERENCES users(id),
			content TEXT NOT NULL,
			sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			return err
		}
	}
	return nil
}