.PHONY: help setup dev test clean build docker-up docker-down

# Default target
help:
    @echo "Available commands:"
    @echo "  setup     - Install dependencies and setup project"
    @echo "  dev       - Start development servers (backend + frontend)"
    @echo "  test      - Run all tests (backend + frontend)"
    @echo "  lint      - Run linting (backend + frontend)"
    @echo "  build     - Build production binaries"
    @echo "  docker-up - Start services with Docker Compose"
    @echo "  docker-down - Stop Docker Compose services"
    @echo "  clean     - Clean build artifacts and dependencies"

# Setup project
setup:
    @echo "Setting up project..."
    cd backend && go mod tidy && go install github.com/swaggo/swag/cmd/swag@latest
    cd frontend && npm install
    @echo "Setup complete!"

# Development
dev:
    @echo "Starting development servers..."
    @echo "Backend will be available at http://localhost:8080"
    @echo "Frontend will be available at http://localhost:3000"
    @echo "API docs will be available at http://localhost:8080/docs/index.html"
    docker-compose up -d postgres
    cd backend && swag init && go run main.go &
    cd frontend && npm run dev

# Testing
test:
    @echo "Running backend tests..."
    cd backend && go test -v ./...
    @echo "Running frontend tests..."
    cd frontend && npm test

test-integration:
    @echo "Running integration tests..."
    cd backend && go test -v -tags=integration ./tests/...

test-e2e:
    @echo "Running E2E tests..."
    cd frontend && npm run test:e2e

test-coverage:
    @echo "Running tests with coverage..."
    cd backend && go test -cover ./...
    cd frontend && npm run test:coverage

# Linting
lint:
    @echo "Running backend linting..."
    cd backend && golangci-lint run
    @echo "Running frontend linting..."
    cd frontend && npm run lint

# Building
build:
    @echo "Building backend..."
    cd backend && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .
    @echo "Building frontend..."
    cd frontend && npm run build

# Docker commands
docker-up:
    docker-compose up -d

docker-down:
    docker-compose down

docker-logs:
    docker-compose logs -f

# Database operations
db-migrate:
    @echo "Running database migrations..."
    cd backend && go run main.go

db-reset:
    @echo "Resetting database..."
    docker-compose exec postgres psql -U postgres -d appdb -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
    cd backend && go run main.go

# API Documentation
docs:
    @echo "Generating API documentation..."
    cd backend && swag init
    @echo "API docs available at http://localhost:8080/docs/index.html"

# Clean
clean:
    @echo "Cleaning up..."
    cd backend && rm -f app main docs/swagger.json docs/swagger.yaml
    cd frontend && rm -rf dist node_modules/.cache
    docker-compose down -v

# Install development tools
install-tools:
    @echo "Installing development tools..."
    go install github.com/swaggo/swag/cmd/swag@latest
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    cd frontend && npm install -g @vitest/ui

# Production deployment
deploy-staging:
    @echo "Deploying to staging..."
    # Add staging deployment commands here

deploy-production:
    @echo "Deploying to production..."
    # Add production deployment commands here

# Security scanning
security-scan:
    @echo "Running security scans..."
    cd backend && gosec ./...
    cd frontend && npm audit --audit-level high

# Performance testing
perf-test:
    @echo "Running performance tests..."
    # Add performance testing commands here

# Health check
health:
    @echo "Checking service health..."
    curl -f http://localhost:8080/health || echo "Backend health check failed"
    curl -f http://localhost:3000 || echo "Frontend health check failed"