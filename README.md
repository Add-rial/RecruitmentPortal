# Recruitment Portal

A backend recruitment platform built using Go and Gin that supports applicants, recruiters, and administrators. The application provides secure authentication, role-based access control, and REST APIs for managing recruitment workflows.

## Features

- User registration and authentication
- JWT-based login and authorization
- Role-Based Access Control (RBAC)
- Job posting and management
- Skill management
- Applicant profile management
- Resume upload support
- PostgreSQL database integration
- Raw SQL queries for database operations

## Tech Stack

- Go
- Gin
- PostgreSQL
- JWT

## Project Structure

```
config/
handlers/
middleware/
models/
routes/
utils/
```

## Getting Started

### Clone the repository

```bash
git clone https://github.com/Add-rial/RecruitmentPortal.git
```

### Install dependencies

```bash
go mod tidy
```

### Configure environment variables

Create a `.env` file and add the required database credentials and JWT secret.

### Run the project

```bash
go run main.go
```

## Future Improvements

- Resume parsing
- Email notifications
- Pagination and filtering
- Unit testing
- Docker support
