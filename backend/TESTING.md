# 🧪 API Handler Testing Guide

This document describes how to run the automated API handler tests for the Food Ordering System.

## Prerequisites

### 1. PostgreSQL Installation

Before running tests, you need to have PostgreSQL installed and running.

#### macOS
```bash
# Install PostgreSQL using Homebrew
brew install postgresql@14

# Start the PostgreSQL service
brew services start postgresql@14

# Verify installation
psql --version
```

#### Ubuntu/Debian
```bash
# Install PostgreSQL
sudo apt-get update
sudo apt-get install postgresql postgresql-contrib

# Start the service
sudo systemctl start postgresql

# Verify installation
psql --version
```

#### Windows
- Download and install PostgreSQL from https://www.postgresql.org/download/windows/
- Make sure to save the postgres user password
- Add PostgreSQL bin folder to PATH

### 2. Create Test Database

```bash
# Create the test database
createdb food_ordering_test

# Verify the database was created
psql -l | grep food_ordering_test
```

### 3. Go Dependencies

Ensure Go 1.21+ is installed and dependencies are up to date:

```bash
cd backend
go mod tidy
go mod download
```

## Running Tests

### Quick Start

```bash
# Set the test database URL and run tests
export TEST_DATABASE_URL="postgres://postgres:password@localhost/food_ordering_test?sslmode=disable"
cd backend
go test -v ./handlers
```

### Using the test.sh Script

```bash
# Run Go unit tests
./test.sh unit

# Run integration tests (requires running backend server)
./test.sh
```

### Advanced Testing Options

```bash
# Run specific test
go test -v ./handlers -run TestLoginSuccess

# Run with coverage report
go test -cover ./handlers

# Run with detailed coverage
go test -coverprofile=coverage.out ./handlers
go tool cover -html=coverage.out

# Run with race condition detector
go test -race ./handlers

# Run with timeout
go test -timeout 30s ./handlers

# Verbose output with print statements
go test -v ./handlers -run TestCreateOrderSuccess -timeout 30s
```

## Database Connection String

The test suite uses the following environment variable to configure the test database:

```
TEST_DATABASE_URL=postgres://username:password@host:port/database?sslmode=disable
```

Default: `postgres://postgres:password@localhost/food_ordering_test?sslmode=disable`

If your PostgreSQL has no password, use:
```bash
export TEST_DATABASE_URL="postgres://postgres@localhost/food_ordering_test?sslmode=disable"
```

## Test Structure

### TestMain Setup

The `TestMain` function in `handlers_test.go`:
1. Calls `setup()` to initialize the test database
2. Runs all tests via `m.Run()`
3. Calls `teardown()` to clean up

### Setup Steps

1. **Database Connection**: Connects to the test database
2. **Schema Initialization**: Creates all required tables and indexes
3. **Data Seeding**: Inserts fixture data (categories, users, dishes, recommendations, config)

### Cleanup Steps

1. **Table Cleanup**: Drops all tables to ensure a clean state
2. **Connection Close**: Closes the database connection

## Test Coverage

The test suite includes 47+ automated tests covering:

### Authentication (7 tests)
- Login success (admin and regular user)
- Login failures (wrong password, non-existent user, empty credentials)
- Profile access with and without token
- Token validation and expiration

### Dishes (12 tests)
- Get dishes list with pagination
- Get single dish
- Filter by category
- Search functionality
- Seasonal dishes
- Admin CRUD operations
- 404 handling for non-existent dishes

### Categories (5 tests)
- Get all categories
- Admin create/update/delete
- Cascade delete validation

### Orders (8 tests)
- Create orders with single/multiple items
- Get orders list
- Order validation (empty, invalid dish, invalid quantity)
- Database consistency checks
- Total amount calculation

### Favorites (7 tests)
- Add to favorites
- Get favorites list
- Remove from favorites
- Duplicate handling
- 404 for non-existent dishes
- Database consistency

### Configuration (3 tests)
- Get system config
- Update system config
- Admin-only access

### Other Features (5+ tests)
- Recommendations
- User list and search
- Admin access control
- Response time validation

## Troubleshooting

### "connection refused" Error

```bash
# Check if PostgreSQL is running
# macOS
brew services list | grep postgresql

# Ubuntu
sudo systemctl status postgresql

# Start PostgreSQL if not running
# macOS
brew services start postgresql@14

# Ubuntu
sudo systemctl start postgresql
```

### "database does not exist" Error

```bash
# Create the test database
createdb food_ordering_test
```

### "permission denied" Error

```bash
# Check PostgreSQL user permissions
# macOS - might need to use -U postgres
psql -U postgres -c "CREATE DATABASE food_ordering_test;"

# Or run with sudo
sudo -u postgres psql -c "CREATE DATABASE food_ordering_test;"
```

### Port Already in Use

If PostgreSQL is using a different port:

```bash
# Find PostgreSQL port (usually 5432)
sudo netstat -tlnp | grep postgres

# Connect to specific port
export TEST_DATABASE_URL="postgres://postgres:password@localhost:5433/food_ordering_test?sslmode=disable"
go test -v ./handlers
```

### Tests Hang or Timeout

```bash
# Run with explicit timeout
go test -timeout 60s -v ./handlers

# Run single test with debugging
go test -v ./handlers -run TestLoginSuccess -timeout 30s
```

## CI/CD Integration

### GitHub Actions

Add to `.github/workflows/test.yml`:

```yaml
name: Go Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:14
        env:
          POSTGRES_PASSWORD: password
          POSTGRES_DB: food_ordering_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

    steps:
      - uses: actions/checkout@v2
      
      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.21
      
      - name: Run tests
        working-directory: ./backend
        env:
          TEST_DATABASE_URL: postgres://postgres:password@localhost:5432/food_ordering_test?sslmode=disable
        run: go test -v ./handlers
```

## Best Practices

1. **Isolate Test Data**: Each test uses unique data; modifications in one test don't affect others
2. **Clean Database State**: Database is cleaned before and after tests
3. **Use Helper Functions**: `getAuthToken()`, `setupRouter()` simplify test writing
4. **Test Edge Cases**: Include tests for invalid inputs and boundary conditions
5. **Verify Side Effects**: Check that HTTP operations correctly update the database
6. **Use Table-Driven Tests**: For parametrized testing of similar scenarios

## Writing New Tests

### Template

```go
func TestNewFeature(t *testing.T) {
    r := setupRouter()
    
    // Arrange
    token, _ := getAuthToken("admin", "admin123")
    payload := []byte(`{"key":"value"}`)
    
    // Act
    req := httptest.NewRequest("POST", "/api/v1/endpoint", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)
    
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    
    // Assert
    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }
    
    var response models.ResponseType
    json.Unmarshal(w.Body.Bytes(), &response)
    
    if response.Field != expectedValue {
        t.Errorf("Expected %v, got %v", expectedValue, response.Field)
    }
}
```

### Key Points

- Use descriptive test names: `TestFeatureSuccess`, `TestFeatureInvalidInput`
- Include both happy path and error cases
- Verify both HTTP response and database state
- Use `getAuthToken()` for authenticated requests
- Use `setupRouter()` to initialize the test router
- Clean up test data if needed (optional, auto-cleanup happens)

## Performance

- Each test takes ~10-50ms depending on database operations
- Full test suite runs in ~30-60 seconds
- Use `go test -parallel` flag for faster execution on multi-core systems

## Continuous Integration

Tests are designed to run in isolated CI/CD environments:
- No external dependencies required (database is local or containerized)
- Tests clean up after themselves
- Deterministic results (no flakiness)
- Can run in parallel safely

## Additional Resources

- [Go testing documentation](https://golang.org/pkg/testing/)
- [Gin Web Framework Testing](https://gin-gonic.com/docs/testing/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [httptest Package](https://golang.org/pkg/net/http/httptest/)
