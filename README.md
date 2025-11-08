# Full-Stack Application with Comprehensive Documentation and Testing

A modern web application with Go backend, React frontend, and complete testing/documentation infrastructure.

## Architecture

### Backend (Go)
- **Framework**: Gin HTTP framework
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT-based authentication
- **Documentation**: Swagger/OpenAPI 3.0 with swaggo
- **Testing**: Unit tests, integration tests

### Frontend (React/TypeScript)
- **Framework**: React 18 with TypeScript
- **Testing**: Vitest for unit/component tests, Cypress for e2e
- **State Management**: React Context API
- **HTTP Client**: Axios

### Infrastructure
- **CI/CD**: GitHub Actions
- **Containerization**: Docker and Docker Compose
- **Database**: PostgreSQL

## Features

- User authentication (login/register)
- Content browsing and search
- Recommendation system
- Admin CRUD operations
- Role-based access control

## Setup

### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- Docker (optional)

### Local Development

1. **Clone and setup**
   ```bash
   git clone <repository-url>
   cd <repository-name>
   ```

2. **Backend setup**
   ```bash
   cd backend
   go mod init example.com/app
   go mod tidy
   export DATABASE_URL="postgres://user:password@localhost/dbname?sslmode=disable"
   export JWT_SECRET="your-secret-key"
   go run main.go
   ```

3. **Frontend setup**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

4. **Database setup**
   ```bash
   # Using Docker Compose (recommended)
   docker-compose up -d postgres
   ```

## API Documentation

Once the backend is running, visit:
- Swagger UI: `http://localhost:8080/docs/index.html`
- OpenAPI JSON: `http://localhost:8080/docs/doc.json`

## Testing

### Backend Tests
```bash
cd backend
go test ./...                    # Unit tests
go test -tags=integration ./... # Integration tests
```

### Frontend Tests
```bash
cd frontend
npm run test                     # Unit tests (Vitest)
npm run test:e2e                 # E2E tests (Cypress)
```

### All Tests
```bash
npm run test:all                 # Run all test suites
```

## Deployment

### Docker Deployment
```bash
docker-compose up -d
```

### Environment Variables

#### Backend
- `DATABASE_URL`: PostgreSQL connection string
- `JWT_SECRET`: Secret key for JWT tokens
- `PORT`: Server port (default: 8080)
- `ENVIRONMENT`: Environment (development/production)

#### Frontend
- `VITE_API_URL`: Backend API URL
- `VITE_ENVIRONMENT`: Environment (development/production)

## Project Structure

```
├── backend/                 # Go backend application
│   ├── cmd/                # Application entry points
│   ├── internal/           # Internal packages
│   │   ├── api/           # HTTP handlers and routes
│   │   ├── auth/          # Authentication logic
│   │   ├── models/        # Data models
│   │   ├── services/      # Business logic
│   │   └── database/      # Database configuration
│   ├── tests/             # Test files
│   └── docs/              # Generated API docs
├── frontend/               # React frontend
│   ├── src/
│   │   ├── components/    # React components
│   │   ├── pages/         # Page components
│   │   ├── hooks/         # Custom hooks
│   │   ├── services/      # API services
│   │   └── utils/         # Utility functions
│   ├── tests/             # Test files
│   └── cypress/           # E2E test specs
├── docs/                   # Documentation
├── .github/                # GitHub Actions workflows
└── docker-compose.yml      # Docker configuration
```

## Development Guidelines

- Follow Go best practices and golangci-lint recommendations
- Use TypeScript strict mode for frontend development
- Write tests for all new features
- Update documentation when adding new API endpoints
- Ensure all tests pass before committing
