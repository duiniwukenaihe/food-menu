# Full-Stack Application - Monorepo Bootstrap

A modern full-stack web application with Go backend, Vue 3 frontend, and comprehensive infrastructure support.

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.21+ (for local backend development)
- Node.js 18+ (for local frontend development)
- PostgreSQL 14+ (if running locally without Docker)

### Local Development with Docker

```bash
# Start all services (database, API, frontend)
docker-compose up

# Services will be available at:
# - Frontend: http://localhost:3000
# - Backend API: http://localhost:8080/api/v1
# - API Docs: http://localhost:8080/docs/index.html
# - MinIO Console: http://localhost:9001
```

### Local Backend Development

1. **Setup environment**
   ```bash
   cd backend
   cp .env.example .env
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Start PostgreSQL and MinIO** (via Docker)
   ```bash
   docker-compose up postgres minio -d
   ```

4. **Run the server**
   ```bash
   go run main.go
   ```

The backend will start on `http://localhost:8080`

### Local Frontend Development

1. **Setup environment**
   ```bash
   cd frontend
   cp .env.example .env
   npm install
   ```

2. **Start development server**
   ```bash
   npm run dev
   ```

The frontend will be available at `http://localhost:3000`

## Architecture

### Backend (Go)

**Structure:**
```
backend/
├── cmd/              # Application entry points
├── internal/
│   ├── api/         # HTTP handlers and routes
│   ├── auth/        # Authentication logic
│   ├── config/      # Configuration management (Viper)
│   ├── database/    # Database initialization and management
│   ├── models/      # Data models
│   └── services/    # Business logic
├── pkg/
│   ├── errors/      # Centralized error handling
│   └── logger/      # Structured logging (Zap)
├── main.go          # Application entry point
├── .env.example     # Environment variables template
└── Dockerfile       # Container build configuration
```

**Tech Stack:**
- **Framework:** Gin HTTP framework
- **Database:** PostgreSQL with GORM ORM
- **Authentication:** JWT-based with bcrypt password hashing
- **Configuration:** Viper for environment-based configuration
- **Logging:** Zap structured logging
- **Storage:** MinIO S3-compatible object storage
- **API Documentation:** Swagger/OpenAPI 3.0

**Key Features:**
- Clean architecture with separated concerns
- Centralized configuration management
- Structured logging for debugging and monitoring
- Comprehensive error handling
- JWT-based authentication with role-based access control
- S3-compatible file storage integration

### Frontend (Vue 3)

**Structure:**
```
frontend/
├── src/
│   ├── components/   # Reusable Vue components
│   ├── pages/        # Page-level components
│   ├── router/       # Vue Router configuration
│   ├── stores/       # Pinia state management
│   ├── api/          # API client services
│   ├── types/        # TypeScript type definitions
│   └── App.vue       # Root component
├── cypress/          # End-to-end tests
├── public/           # Static assets
├── .env.example      # Environment variables template
└── vite.config.ts    # Vite configuration
```

**Tech Stack:**
- **Framework:** Vue 3 with Composition API
- **Build Tool:** Vite
- **Language:** TypeScript
- **State Management:** Pinia
- **Routing:** Vue Router 4
- **UI Components:** Element Plus
- **HTTP Client:** Axios
- **Styling:** Tailwind CSS
- **Testing:** Vitest + Cypress
- **Code Quality:** ESLint + Prettier

**Key Features:**
- Type-safe development with TypeScript
- Reactive state management with Pinia
- Client-side routing with Vue Router
- Pre-built UI components with Element Plus
- Responsive design with Tailwind CSS
- Comprehensive test coverage

## Configuration

### Environment Variables

**Backend (.env)**
```env
# Server Configuration
PORT=8080
ENVIRONMENT=development

# Database Configuration
DATABASE_URL=postgres://postgres:password@localhost:5432/appdb?sslmode=disable

# Authentication
JWT_SECRET=your-secret-key-change-in-production

# Storage Configuration (MinIO)
STORAGE_TYPE=minio
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=media
MINIO_USE_SSL=false
```

**Frontend (.env)**
```env
# API Configuration
VITE_API_URL=http://localhost:8080/api/v1

# Environment
VITE_ENVIRONMENT=development

# Features
VITE_ENABLE_DEVTOOLS=true
```

Copy `.env.example` to `.env` in both directories and update values for your environment.

## Docker Compose Services

- **PostgreSQL:** Database server (port 5432)
- **MinIO:** S3-compatible object storage (ports 9000, 9001)
- **Backend:** Go API server (port 8080)
- **Frontend:** Vue.js development server (port 3000)

### Docker Commands

```bash
# Start all services
docker-compose up -d

# Stop all services
docker-compose down

# View logs
docker-compose logs -f

# Reset database (warning: deletes data)
docker-compose down -v
```

## Development

### Makefile Commands

```bash
make help          # Show available commands
make setup         # Install dependencies
make dev           # Start development servers
make test          # Run all tests
make lint          # Run linting
make build         # Build production binaries
make docker-up     # Start Docker Compose
make docker-down   # Stop Docker Compose
make docs          # Generate API documentation
```

### Project Scripts

**Backend:**
```bash
cd backend
go run main.go              # Run server
go test ./...              # Run tests
go test -v -tags=integration ./tests/... # Integration tests
```

**Frontend:**
```bash
cd frontend
npm run dev                 # Development server
npm run build              # Production build
npm run test               # Unit tests
npm run test:e2e           # E2E tests
npm run lint               # Linting
```

## API Documentation

Once the backend is running, visit:
- **Swagger UI:** http://localhost:8080/docs/index.html
- **OpenAPI JSON:** http://localhost:8080/docs/doc.json
- **Health Check:** http://localhost:8080/health

## Testing

### Backend Testing
```bash
cd backend
go test ./...                           # Unit tests
go test -v -tags=integration ./tests/...  # Integration tests
go test -cover ./...                   # With coverage
```

### Frontend Testing
```bash
cd frontend
npm run test                # Unit/component tests
npm run test:ui            # Test UI
npm run test:e2e           # E2E tests with Cypress
npm run test:coverage      # Coverage report
```

## Deployment

### Docker Build
```bash
docker-compose up -d
```

### Environment Variables for Production
- Set `ENVIRONMENT=production`
- Update `JWT_SECRET` to a strong, secure value
- Use actual PostgreSQL credentials
- Configure MinIO with proper credentials
- Update `VITE_API_URL` to production backend URL

## Project Structure

```
.
├── backend/                 # Go backend application
│   ├── cmd/                # Entry points
│   ├── internal/           # Internal packages
│   ├── pkg/                # Shared utilities
│   ├── go.mod             # Go module definition
│   └── Dockerfile         # Backend container
├── frontend/               # Vue 3 frontend application
│   ├── src/               # Source code
│   ├── cypress/           # E2E tests
│   ├── package.json       # Dependencies
│   └── Dockerfile         # Frontend container
├── docker-compose.yml      # Service orchestration
├── Makefile               # Development commands
├── .env.example           # Template env vars
└── README.md              # This file
```

## Key Files and Their Purpose

- **backend/.env.example** - Template for backend environment variables
- **frontend/.env.example** - Template for frontend environment variables
- **backend/main.go** - Backend entry point with config and logger initialization
- **backend/internal/config/config.go** - Configuration management with Viper
- **backend/pkg/logger/logger.go** - Structured logging setup with Zap
- **backend/pkg/errors/errors.go** - Centralized error handling utilities
- **docker-compose.yml** - Services definition (PostgreSQL, MinIO, Backend, Frontend)
- **Makefile** - Development workflow automation

## Development Guidelines

### Backend Development
- Follow Go best practices and conventions
- Use the config package for environment variables
- Use the logger package for structured logging
- Use the errors package for consistent error handling
- Write tests for all business logic
- Keep handlers in `internal/api`, services in `internal/services`

### Frontend Development
- Use TypeScript strict mode
- Follow Vue 3 Composition API best practices
- Use Pinia for state management
- Implement components with Element Plus for consistency
- Write tests for complex logic and components
- Format code with Prettier before committing

### General Guidelines
- Commit messages should be descriptive
- Keep dependencies updated
- Document new APIs and features
- Ensure tests pass before pushing

## Troubleshooting

### Database Connection Issues
- Verify PostgreSQL is running: `docker-compose ps postgres`
- Check DATABASE_URL environment variable is correct
- Run migrations: `docker-compose exec backend go run main.go`

### Port Conflicts
- If ports are already in use, modify `docker-compose.yml` port mappings
- Or kill existing processes: `lsof -ti:8080 | xargs kill -9`

### MinIO Connection Issues
- Verify MinIO is running: `docker-compose ps minio`
- Check MinIO credentials in environment variables
- Access MinIO console at http://localhost:9001

### Frontend Build Issues
- Clear cache: `npm run clean` (if script exists) or `rm -rf node_modules package-lock.json`
- Reinstall dependencies: `npm install`
- Clear Vite cache: `rm -rf .vite`

## Next Steps

1. Configure environment variables in `.env` files
2. Start services with Docker Compose
3. Access the application at http://localhost:3000
4. Review API documentation at http://localhost:8080/docs
5. Start developing!

## Resources

- [Go Best Practices](https://golang.org/doc/effective_go)
- [Gin Framework](https://github.com/gin-gonic/gin)
- [Vue 3 Documentation](https://vuejs.org/)
- [Element Plus Components](https://element-plus.org/)
- [Viper Configuration](https://github.com/spf13/viper)
- [Zap Logging](https://github.com/uber-go/zap)

## Support

For issues or questions, please refer to the documentation or create an issue in the repository.
