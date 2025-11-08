# Test Documentation

This document outlines the comprehensive testing strategy for the Full-Stack Application.

## Test Coverage Overview

### Backend Tests

#### Unit Tests
- **Authentication Service** (`backend/internal/auth/auth_test.go`)
  - Password hashing and verification
  - JWT token generation and validation
  - Token security and expiration

- **Content Service** (`backend/internal/services/content_service_test.go`)
  - Content validation logic
  - Business logic testing
  - Data transformation

#### Integration Tests
- **API Endpoints** (`backend/tests/integration/api_test.go`)
  - User registration and login flow
  - Content retrieval and pagination
  - Protected route authentication
  - Database integration

### Frontend Tests

#### Unit Tests
- **Component Testing** (`frontend/src/components/__tests__/Layout.test.tsx`)
  - Layout component rendering
  - Navigation functionality
  - Authentication state handling

- **Page Testing** (`frontend/src/pages/__tests__/Login.test.tsx`)
  - Login form validation
  - User interaction handling
  - API integration mocking

#### E2E Tests (Cypress)
- **Authentication Flow** (`frontend/cypress/e2e/auth.cy.ts`)
  - User registration
  - Login functionality
  - Logout process
  - Error handling

- **Content Browsing** (`frontend/cypress/e2e/content.cy.ts`)
  - Content listing
  - Search and filtering
  - Content detail view
  - View count tracking

- **Admin Operations** (`frontend/cypress/e2e/admin.cy.ts`)
  - User CRUD operations
  - Content management
  - Category management
  - Admin dashboard functionality

## Test Execution

### Backend Tests

```bash
# Run all unit tests
cd backend
go test ./...

# Run integration tests
go test -tags=integration ./tests/...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Frontend Tests

```bash
# Run unit tests
cd frontend
npm test

# Run tests with coverage
npm run test:coverage

# Run tests in watch mode
npm run test:watch

# Run tests with UI
npm run test:ui
```

### E2E Tests

```bash
# Run all E2E tests
cd frontend
npm run test:e2e

# Run E2E tests with interactive mode
npm run test:e2e:open

# Run specific test file
npx cypress run --spec "cypress/e2e/auth.cy.ts"
```

## Test Data Management

### Test Database Setup
- Uses PostgreSQL for integration tests
- Automatic database cleanup between tests
- Seed data for consistent testing

### Mock Data
- User accounts for different roles (admin, user)
- Sample content and categories
- Test fixtures for consistent testing

## Test Scenarios

### Authentication Test Matrix

| Scenario | Description | Expected Outcome |
|----------|-------------|------------------|
| Valid Registration | New user with valid data | Account created, user logged in |
| Invalid Registration | Duplicate username/email | Error message, no account created |
| Valid Login | Correct credentials | JWT token returned, user authenticated |
| Invalid Login | Wrong credentials | Error message, no authentication |
| Token Validation | Valid/invalid JWT tokens | Proper access control |
| Session Expiration | Expired token | Redirect to login |

### Content Management Test Matrix

| Feature | Test Cases | Coverage |
|---------|------------|----------|
| Content Listing | Pagination, filtering, search | Full coverage |
| Content Creation | Validation, authorization, publishing | Full coverage |
| Content Updates | Partial updates, permissions | Full coverage |
| Content Deletion | Soft delete, admin only | Full coverage |
| View Tracking | Increment on view, uniqueness | Full coverage |

### Admin Operations Test Matrix

| Operation | Test Cases | Coverage |
|-----------|------------|----------|
| User Management | CRUD, role assignment, activation | Full coverage |
| Content Management | Full lifecycle, publishing workflow | Full coverage |
| Category Management | CRUD, content association | Full coverage |
| Dashboard Analytics | Data aggregation, charts | Full coverage |

## Performance Testing

### Load Testing Scenarios
- Concurrent user authentication
- High-volume content browsing
- Admin panel stress testing
- Database query optimization

### Metrics Monitored
- Response times
- Database query performance
- Memory usage
- Error rates

## Security Testing

### Authentication Security
- JWT token strength
- Password hashing
- Session management
- CSRF protection

### API Security
- Input validation
- SQL injection prevention
- XSS protection
- Rate limiting

### Authorization Testing
- Role-based access control
- Resource ownership
- Admin privilege escalation
- API endpoint protection

## Continuous Integration

### GitHub Actions Pipeline
1. **Backend Tests**
   - Unit tests with Go
   - Integration tests with PostgreSQL
   - Code coverage reporting
   - Swagger documentation generation

2. **Frontend Tests**
   - Unit tests with Vitest
   - Component testing
   - Linting and type checking
   - Build verification

3. **E2E Tests**
   - Full application deployment
   - Cypress test execution
   - Cross-browser testing
   - Performance monitoring

4. **Security Scanning**
   - Dependency vulnerability scanning
   - Code security analysis
   - Static code analysis

### Quality Gates
- Minimum 80% code coverage
- All tests must pass
- No high-severity security issues
- Build must complete within time limits

## Test Environment Configuration

### Development Environment
- Hot reload for both frontend and backend
- Database seeding for consistent test data
- Debug logging enabled
- Mock external services

### Staging Environment
- Production-like configuration
- Real database with test data
- Full CI/CD pipeline
- Performance monitoring

### Production Environment
- Smoke tests after deployment
- Health checks
- Error monitoring
- Performance metrics

## Test Data Examples

### Sample User Data
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "firstName": "Test",
  "lastName": "User",
  "role": "user",
  "isActive": true
}
```

### Sample Content Data
```json
{
  "title": "Sample Article",
  "description": "A test article",
  "body": "Full article content here...",
  "categoryId": 1,
  "tags": "test,sample",
  "isPublished": true
}
```

## Best Practices

### Test Writing Guidelines
- Write descriptive test names
- Use AAA pattern (Arrange, Act, Assert)
- Test one thing per test
- Use meaningful test data
- Avoid test interdependence

### Code Coverage Targets
- Backend: 85% minimum
- Frontend: 80% minimum
- Critical paths: 100% coverage

### Test Maintenance
- Regular test review and cleanup
- Update tests with feature changes
- Monitor test execution times
- Refactor flaky tests

## Troubleshooting

### Common Issues
- Database connection failures in integration tests
- CORS issues in E2E tests
- Mock service misconfiguration
- Test data cleanup failures

### Debugging Strategies
- Enable verbose logging
- Use database transaction rollbacks
- Isolate failing tests
- Check test environment configuration

## Reporting

### Test Reports
- HTML coverage reports
- Test execution summaries
- Performance benchmarks
- Security scan results

### Metrics Tracking
- Test execution time trends
- Coverage percentage changes
- Bug detection rates
- Performance regressions