# Monorepo Bootstrap Guide

This document describes how the monorepo has been bootstrapped and provides guidance for development.

## Project Structure

The monorepo follows a clean architecture pattern with clear separation of concerns:

```
.
├── backend/                    # Go backend application
│   ├── cmd/                   # Application entry points (reserved for cmd/server)
│   ├── internal/              # Internal packages
│   │   ├── api/              # HTTP handlers with OpenAPI annotations
│   │   ├── auth/             # Authentication logic
│   │   ├── config/           # Configuration management (Viper)
│   │   ├── database/         # Database initialization
│   │   ├── models/           # Data models
│   │   └── services/         # Business logic services
│   ├── pkg/                   # Shared utilities across the app
│   │   ├── errors/           # Centralized error handling
│   │   └── logger/           # Structured logging (Zap)
│   ├── main.go               # Application entry point
│   ├── go.mod                # Go module definition
│   ├── go.sum                # Go dependencies lock file
│   ├── .env.example          # Environment variables template
│   ├── .air.toml             # Hot reload configuration
│   └── Dockerfile            # Container build
│
├── frontend/                   # Vue 3 frontend application
│   ├── src/                  # Source code
│   │   ├── components/       # Reusable Vue components
│   │   ├── pages/            # Page components
│   │   ├── router/           # Vue Router configuration
│   │   ├── stores/           # Pinia state management
│   │   ├── api/              # API client services
│   │   ├── types/            # TypeScript types
│   │   └── App.vue           # Root component
│   ├── cypress/              # E2E tests
│   ├── public/               # Static assets
│   ├── package.json          # Dependencies
│   ├── tsconfig.json         # TypeScript configuration
│   ├── vite.config.ts        # Vite configuration
│   ├── .env.example          # Environment variables template
│   ├── .prettierrc            # Prettier formatting rules
│   ├── .prettierignore        # Prettier ignore patterns
│   ├── .eslintrc.cjs         # ESLint configuration
│   └── Dockerfile            # Container build
│
├── docker-compose.yml         # Service orchestration
├── Makefile                   # Development commands
├── .prettierrc                # Root Prettier configuration
├── .gitignore                 # Git ignore rules
├── README.md                  # Project overview
└── docs/                      # Documentation
    └── BOOTSTRAP.md          # This file
```

## Backend Architecture

### Clean Architecture Layers

The backend follows clean architecture principles with distinct layers:

1. **HTTP Layer (internal/api)**: Handles incoming requests, validation, and response formatting
2. **Service Layer (internal/services)**: Contains business logic
3. **Data Layer (internal/models)**: Database models and data access
4. **Authentication Layer (internal/auth)**: JWT validation and user authentication

### Configuration Management

**Viper Configuration** (`internal/config/config.go`):
- Loads configuration from `.env` file or environment variables
- Supports environment-based settings (development, production)
- Provides type-safe config access
- Validates required configuration at startup

**Environment Variables**:
```env
PORT=8080
ENVIRONMENT=development
DATABASE_URL=postgres://postgres:password@localhost:5432/appdb?sslmode=disable
JWT_SECRET=your-secret-key-change-in-production
STORAGE_TYPE=minio
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=media
MINIO_USE_SSL=false
```

### Structured Logging

**Zap Logger** (`pkg/logger/logger.go`):
- Production-grade structured logging
- Different output formats for development vs. production
- Colored output in development mode
- Convenient functions: `Info()`, `Error()`, `Warn()`, `Debug()`, `Fatal()`
- Proper log levels and field support

Usage:
```go
import "example.com/app/pkg/logger"
import "go.uber.org/zap"

logger.Info("User created", 
    zap.String("username", user.Username),
    zap.String("email", user.Email),
)
```

### Centralized Error Handling

**Error Package** (`pkg/errors/errors.go`):
- Consistent error responses across the API
- HTTP status code mapping
- Structured error creation
- Convenience methods for common errors:
  - `BadRequest(msg)`
  - `Unauthorized(msg)`
  - `Forbidden(msg)`
  - `NotFound(msg)`
  - `Conflict(msg)`
  - `InternalServerError(msg)`

Usage:
```go
import "example.com/app/pkg/errors"

if !user.IsAdmin {
    return errors.Forbidden("Only admins can perform this action")
}
```

## Frontend Architecture

### Vue 3 + Vite Stack

**Key Technologies**:
- **Vue 3**: Modern reactive framework with Composition API
- **Vite**: Fast build tool and dev server
- **TypeScript**: Type-safe development
- **Pinia**: State management (simpler than Vuex)
- **Vue Router**: Client-side routing
- **Element Plus**: Enterprise UI component library
- **Tailwind CSS**: Utility-first CSS framework
- **Axios**: HTTP client

### Project Organization

- **Components**: Reusable UI building blocks
- **Pages**: Full page components
- **Router**: Route definitions and navigation
- **Stores**: Pinia stores for state management
- **API**: Centralized API client services
- **Types**: Shared TypeScript interfaces

## Development Workflow

### Local Development Setup

1. **Backend Setup**:
   ```bash
   cd backend
   cp .env.example .env
   go mod tidy
   docker-compose up postgres minio -d
   go run main.go
   ```

2. **Frontend Setup**:
   ```bash
   cd frontend
   cp .env.example .env
   npm install
   npm run dev
   ```

3. **Full Stack with Docker**:
   ```bash
   docker-compose up
   ```

### Development Commands

**Makefile** provides shortcuts for common tasks:

```bash
make setup          # Install all dependencies
make dev            # Start dev servers (backend + frontend)
make test           # Run all tests
make lint           # Run linting
make build          # Build production binaries
make docker-up      # Start Docker Compose
make docker-down    # Stop Docker Compose
make clean          # Clean build artifacts
make docs           # Generate API docs
```

### Code Quality

**Formatting**:
- `prettier` for both frontend and backend code formatting
- Configuration: `.prettierrc`
- Run: `npm run lint` (frontend already includes prettier)

**Linting**:
- ESLint for frontend TypeScript
- golangci-lint for backend Go

## Docker Compose Services

- **PostgreSQL** (postgres:15): Relational database on port 5432
- **MinIO** (minio/minio): S3-compatible storage on ports 9000, 9001
- **Backend**: Go API server on port 8080
- **Frontend**: Vue dev server on port 3000

Services use health checks and proper dependencies to ensure startup order.

## Configuration Files

### Backend Configuration Files

- **main.go**: Entry point that initializes config, logger, and database
- **internal/config/config.go**: Viper-based configuration loader
- **pkg/logger/logger.go**: Zap structured logging setup
- **pkg/errors/errors.go**: Error handling utilities
- **go.mod / go.sum**: Dependency management
- **.env.example**: Environment variables template
- **.air.toml**: Hot reload configuration for development
- **Dockerfile**: Multi-stage build for production

### Frontend Configuration Files

- **vite.config.ts**: Vite build configuration
- **tsconfig.json**: TypeScript compiler options
- **package.json**: Dependencies and scripts
- **.env.example**: Environment variables template
- **.eslintrc.cjs**: ESLint rules
- **.prettierrc**: Prettier formatting rules
- **.prettierignore**: Files to ignore during formatting
- **Dockerfile**: Production-ready container build

### Root Configuration

- **.gitignore**: Version control exclusions
- **.prettierrc**: Root-level formatting rules
- **docker-compose.yml**: Service definitions
- **Makefile**: Development commands
- **README.md**: Project documentation

## Adding New Features

### Backend: Adding a New Service

1. Create models in `internal/models/`
2. Create service logic in `internal/services/`
3. Create API handlers in `internal/api/`
4. Add routes in `main.go`
5. Add migration code if needed

### Frontend: Adding a New Page

1. Create page component in `src/pages/`
2. Define types in `src/types/` if needed
3. Create API client in `src/api/` if needed
4. Add Pinia store in `src/stores/` if needed
5. Add route in `src/router/`
6. Use Element Plus components for UI

## Deployment

### Environment-Based Configuration

**Development** (`ENVIRONMENT=development`):
- Colored console logging
- Full debug information
- Hot reload enabled
- CORS disabled for all origins

**Production** (`ENVIRONMENT=production`):
- JSON structured logging
- Error tracking
- Performance optimizations
- Strict CORS policy required

### Building for Production

**Backend**:
```bash
cd backend
CGO_ENABLED=0 GOOS=linux go build -o app .
```

**Frontend**:
```bash
cd frontend
npm run build
```

### Docker Deployment

```bash
docker-compose up -d
```

All services start with proper health checks and dependencies.

## Troubleshooting

### Database Connection Issues
- Ensure PostgreSQL is running: `docker-compose ps postgres`
- Check DATABASE_URL in `.env`
- Verify port 5432 is not already in use

### Go Module Issues
- Clear cache: `go clean -modcache`
- Reinstall: `go mod tidy`
- Update: `go get -u ./...`

### Frontend Build Issues
- Clear node_modules: `rm -rf node_modules && npm install`
- Clear cache: `rm -rf .vite dist`
- Update npm: `npm update`

### Port Conflicts
- Find process: `lsof -ti:8080`
- Kill process: `kill -9 <PID>`
- Or change port in .env or docker-compose.yml

## Next Steps

1. Review the main README.md for quick start
2. Explore the backend API at http://localhost:8080/docs/index.html
3. Start the frontend at http://localhost:3000
4. Create your first feature following the patterns established
5. Keep configuration in environment variables, never commit secrets

## Resources

- [Go Best Practices](https://golang.org/doc/effective_go)
- [Gin Web Framework](https://gin-gonic.com/)
- [Vue 3 Guide](https://vuejs.org/guide/)
- [Vite Documentation](https://vitejs.dev/)
- [Viper Configuration](https://github.com/spf13/viper)
- [Zap Logger](https://pkg.go.dev/go.uber.org/zap)
- [Element Plus](https://element-plus.org/)
- [Pinia Documentation](https://pinia.vuejs.org/)
