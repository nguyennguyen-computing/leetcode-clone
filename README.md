# LeetCode Clone

A full-stack LeetCode clone built with Angular frontend and Go backend.

## Tech Stack

### Frontend
- **Angular 17** with standalone components (zoneless)
- **NgRx** for state management
- **NG-ZORRO** for UI components
- **Tailwind CSS** for styling
- **Monaco Editor** for code editing

### Backend
- **Go 1.21** with Gin framework
- **PostgreSQL 15** database
- **JWT** authentication
- **Docker** containerization

### Development Tools
- **Colima** for Docker container management on macOS
- **Docker Compose** for orchestration

## Prerequisites

- Node.js 20+
- Go 1.21+
- Docker and Docker Compose
- Colima (for macOS users)

## Quick Start

### 1. Setup Colima (macOS only)
```bash
./setup-colima.sh
```

### 2. Development Environment
```bash
# Start all services in development mode
docker-compose -f docker-compose.dev.yml up

# Or start individual services
docker-compose -f docker-compose.dev.yml up postgres
docker-compose -f docker-compose.dev.yml up backend-dev
docker-compose -f docker-compose.dev.yml up frontend-dev
```

### 3. Production Environment
```bash
# Build and start all services
docker-compose up --build
```

## Services

- **Frontend**: http://localhost:4200
- **Backend API**: http://localhost:8080
- **PostgreSQL**: localhost:5432

## API Endpoints

### Health Check
- `GET /api/v1/health` - Service health status

### Authentication (Placeholder)
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login

### Problems (Placeholder)
- `GET /api/v1/problems` - List all problems
- `GET /api/v1/problems/:id` - Get specific problem

### Submissions (Placeholder)
- `POST /api/v1/submissions` - Submit solution
- `GET /api/v1/submissions` - Get user submissions

## Database Schema

The PostgreSQL database includes tables for:
- Users (authentication and profiles)
- Problems (coding challenges)
- Test Cases (problem validation)
- Submissions (user solutions)
- User Progress (tracking solved problems)

## Development

### Frontend Development
```bash
cd frontend
npm install
npm start
```

### Backend Development
```bash
cd backend
go mod download
go run main.go
```

### Database Migrations
Database schema is automatically applied when PostgreSQL container starts using the migration files in `backend/migrations/`.

## Project Structure

```
├── frontend/                 # Angular application
│   ├── src/
│   ├── package.json
│   ├── angular.json
│   ├── tailwind.config.js
│   └── Dockerfile
├── backend/                  # Go API server
│   ├── main.go
│   ├── go.mod
│   ├── migrations/
│   └── Dockerfile
├── docker-compose.yml        # Production compose
├── docker-compose.dev.yml    # Development compose
└── setup-colima.sh          # Colima setup script
```

## Next Steps

This is the initial project setup. The following features will be implemented in subsequent tasks:
1. User authentication system
2. Problem management
3. Code execution environment
4. Submission system
5. User interface components
6. Real-time features