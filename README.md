# 🚀 LeetCode Clone

A comprehensive, full-stack LeetCode clone featuring a modern Angular frontend and robust Go backend with real-time code execution, user authentication, and problem management.

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Angular Version](https://img.shields.io/badge/Angular-17+-red.svg)](https://angular.io)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-blue.svg)](https://postgresql.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://docker.com)

## ✨ Features

### 🔐 Authentication & User Management
- JWT-based authentication system
- User registration and login
- Password reset functionality
- Admin role management
- Rate limiting for security

### 📚 Problem Management
- **Complete CRUD operations** for coding problems
- **Advanced filtering** by difficulty, tags, and search queries
- **Rich problem descriptions** with examples and constraints
- **Multi-language support** (JavaScript, Python, Java)
- **Template code** generation for each language
- **Public and hidden test cases**

### 💻 Code Execution Engine
- **Sandboxed code execution** using Docker containers
- **Multi-language support** with secure runtime environments
- **Performance monitoring** (runtime and memory usage)
- **Timeout and memory limit enforcement**
- **Comprehensive error handling** and feedback

### 📊 Submission System
- **Real-time code submission** and evaluation
- **Detailed execution results** with test case feedback
- **Submission history** with pagination and filtering
- **Performance metrics** and statistics tracking
- **User progress tracking** and problem completion status
- **Acceptance rate calculations**

### 🎨 Modern Frontend
- **Angular 17** with standalone components (zoneless)
- **NgRx** for state management
- **NG-ZORRO** UI component library
- **Monaco Editor** for code editing with syntax highlighting
- **Tailwind CSS** for responsive design
- **Real-time updates** and optimistic UI

## 🛠 Tech Stack

### Frontend
- **Angular 17** - Modern web framework with standalone components
- **NgRx** - Reactive state management
- **NG-ZORRO** - Enterprise-class UI components
- **Tailwind CSS** - Utility-first CSS framework
- **Monaco Editor** - VS Code editor for the web
- **TypeScript** - Type-safe JavaScript

### Backend
- **Go 1.21** - High-performance backend language
- **Gin** - Fast HTTP web framework
- **PostgreSQL 15** - Robust relational database
- **JWT** - Secure authentication tokens
- **Docker** - Containerized code execution
- **GORM** - Go ORM for database operations

### Infrastructure
- **Docker & Docker Compose** - Containerization and orchestration
- **Colima** - Docker container management for macOS
- **PostgreSQL** - Primary database
- **Nginx** - Production web server

## 🚀 Quick Start

### Prerequisites
- **Node.js 20+**
- **Go 1.21+**
- **Docker & Docker Compose**
- **Colima** (for macOS users)

### 1. Clone and Setup
```bash
git clone <repository-url>
cd leetcode-clone

# Setup Colima (macOS only)
./setup-colima.sh
```

### 2. Development Environment
```bash
# Quick setup with sample data
make setup-dev

# Or manual setup
docker-compose -f docker-compose.dev.yml up -d
make seed-db
```

### 3. Access the Application
- **Frontend**: http://localhost:4200
- **Backend API**: http://localhost:8080
- **API Documentation**: http://localhost:8080/api/v1/health

## 📖 API Documentation

### 🔐 Authentication Endpoints
```
POST   /api/v1/auth/register              - User registration
POST   /api/v1/auth/login                 - User login
POST   /api/v1/auth/password-reset        - Request password reset
POST   /api/v1/auth/password-reset/confirm - Confirm password reset
```

### 📚 Problem Management
```
# Public Endpoints
GET    /api/v1/problems                   - List problems with filtering
GET    /api/v1/problems/search            - Search problems
GET    /api/v1/problems/:id               - Get problem by ID
GET    /api/v1/problems/slug/:slug        - Get problem by slug
GET    /api/v1/problems/:id/testcases     - Get test cases

# Admin Endpoints (Authentication Required)
POST   /api/v1/admin/problems             - Create problem
PUT    /api/v1/admin/problems/:id         - Update problem
DELETE /api/v1/admin/problems/:id         - Delete problem
POST   /api/v1/admin/problems/:id/testcases - Create test case
PUT    /api/v1/admin/testcases/:id        - Update test case
DELETE /api/v1/admin/testcases/:id        - Delete test case
```

### 💻 Code Execution
```
POST   /api/v1/execute/run                - Run code against public test cases
POST   /api/v1/execute/submit             - Submit code for evaluation
POST   /api/v1/execute/validate           - Validate code syntax
GET    /api/v1/execute/languages          - Get supported languages
```

### 📊 Submission Management
```
POST   /api/v1/submissions                - Create submission
GET    /api/v1/submissions/:id            - Get submission by ID
GET    /api/v1/submissions/me             - Get current user submissions
GET    /api/v1/submissions/stats/me       - Get user statistics
GET    /api/v1/problems/:id/submissions   - Get problem submissions
```

### Query Parameters
- **Pagination**: `page`, `page_size` (max 100)
- **Filtering**: `difficulty`, `tags`, `problem_id`
- **Sorting**: `sort_by`, `sort_order`
- **Search**: `q` (query string)

## 🏗 Project Structure

```
leetcode-clone/
├── 📁 frontend/                    # Angular Application
│   ├── 📁 src/app/
│   │   ├── 📁 auth/               # Authentication module
│   │   │   ├── 📁 components/     # Login, register components
│   │   │   ├── 📁 guards/         # Route guards
│   │   │   ├── 📁 services/       # Auth services
│   │   │   └── 📁 store/          # NgRx auth state
│   │   ├── 📁 problems/           # Problem management
│   │   │   ├── 📁 components/     # Problem list, solve, editor
│   │   │   ├── 📁 services/       # Problem services
│   │   │   └── 📁 store/          # NgRx problem state
│   │   ├── 📄 app.config.ts       # App configuration
│   │   └── 📄 app.routes.ts       # Routing configuration
│   ├── 📄 package.json
│   ├── 📄 angular.json
│   ├── 📄 tailwind.config.js
│   └── 📄 Dockerfile
├── 📁 backend/                     # Go API Server
│   ├── 📁 pkg/
│   │   ├── 📁 auth/               # JWT authentication
│   │   ├── 📁 database/           # Database connection
│   │   ├── 📁 execution/          # Code execution engine
│   │   ├── 📁 handlers/           # HTTP handlers
│   │   ├── 📁 middleware/         # HTTP middleware
│   │   ├── 📁 models/             # Data models
│   │   ├── 📁 repository/         # Data access layer
│   │   └── 📁 services/           # Business logic
│   ├── 📁 migrations/             # Database migrations
│   ├── 📁 scripts/                # Utility scripts
│   ├── 📄 main.go                 # Application entry point
│   ├── 📄 go.mod                  # Go dependencies
│   └── 📄 Dockerfile
├── 📄 docker-compose.yml          # Production compose
├── 📄 docker-compose.dev.yml      # Development compose
├── 📄 Makefile                    # Development commands
└── 📄 setup-colima.sh            # macOS Docker setup
```

## 🔧 Development

### Available Make Commands
```bash
make help          # Show all available commands
make setup-dev     # Complete development environment setup
make seed-db       # Seed database with sample problems
make clean-db      # Clear all database data (WARNING!)
```

### Frontend Development
```bash
cd frontend
npm install
npm start          # Development server on :4200
npm run build      # Production build
npm test           # Run unit tests
npm run lint       # Code linting
```

### Backend Development
```bash
cd backend
go mod download
go run main.go     # Development server on :8080
go test ./...      # Run all tests
go build           # Build binary
```

### Database Operations
```bash
# View database schema
cd backend && go run scripts/validate_schema.go

# Seed with sample data
make seed-db

# Connect to database
docker exec -it leetcode-postgres psql -U leetcode -d leetcode
```

## 🧪 Testing

### Backend Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/services -v
go test ./pkg/handlers -v
```

### Frontend Testing
```bash
# Unit tests
npm test

# E2E tests
npm run e2e

# Test coverage
npm run test:coverage
```

## 🚀 Deployment

### Production Build
```bash
# Build and start all services
docker-compose up --build -d

# View logs
docker-compose logs -f
```

### Environment Variables
```bash
# Backend (.env)
DB_HOST=postgres
DB_PORT=5432
DB_USER=leetcode
DB_PASSWORD=password
DB_NAME=leetcode
JWT_SECRET=your-secret-key
PORT=8080

# Frontend (environment.prod.ts)
API_URL=http://localhost:8080/api/v1
```

## 📊 Database Schema

### Core Tables
- **users** - User accounts and authentication
- **problems** - Coding problems and metadata
- **test_cases** - Problem validation test cases
- **submissions** - User code submissions
- **user_progress** - Problem completion tracking

### Key Features
- **Automatic migrations** on startup
- **Foreign key constraints** for data integrity
- **Indexes** for query performance
- **JSON fields** for flexible data storage

## 🔒 Security Features

### Authentication & Authorization
- **JWT tokens** with expiration
- **Password hashing** with bcrypt
- **Rate limiting** on auth endpoints
- **Role-based access control** (admin/user)

### Code Execution Security
- **Docker sandboxing** for code execution
- **Resource limits** (CPU, memory, time)
- **Network isolation** (no internet access)
- **Input validation** and sanitization
- **Dangerous pattern detection**

## 🎯 Current Implementation Status

### ✅ Completed Features
- [x] **Authentication System** - JWT-based auth with registration/login
- [x] **Problem Management** - Full CRUD with advanced filtering
- [x] **Code Execution Engine** - Sandboxed multi-language execution
- [x] **Submission System** - Complete submission workflow
- [x] **Database Schema** - Comprehensive data model
- [x] **API Documentation** - Detailed endpoint documentation
- [x] **Testing Suite** - Unit and integration tests
- [x] **Docker Setup** - Development and production containers

### 🚧 In Progress
- [ ] **Frontend Components** - Angular UI implementation
- [ ] **Real-time Features** - WebSocket integration
- [ ] **Advanced Analytics** - User performance insights

### 📋 Planned Features
- [ ] **Discussion Forum** - Problem discussions and solutions
- [ ] **Contest System** - Timed coding competitions
- [ ] **Social Features** - User profiles and following
- [ ] **Mobile App** - React Native mobile client

## 🤝 Contributing

1. **Fork the repository**
2. **Create a feature branch** (`git checkout -b feature/amazing-feature`)
3. **Commit your changes** (`git commit -m 'Add amazing feature'`)
4. **Push to the branch** (`git push origin feature/amazing-feature`)
5. **Open a Pull Request**

### Development Guidelines
- Follow **Go** and **Angular** best practices
- Write **comprehensive tests** for new features
- Update **documentation** for API changes
- Use **conventional commits** for commit messages

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **LeetCode** for inspiration
- **Angular Team** for the amazing framework
- **Go Community** for excellent libraries
- **Docker** for containerization technology

---

**Built with ❤️ by the development team**