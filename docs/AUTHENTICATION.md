# Authentication API Documentation

## Overview

This document describes the authentication endpoints implemented in the backend API. The authentication system uses JWT (JSON Web Tokens) for stateless authentication with bcrypt for password hashing.

## Security Features

- **Password Hashing**: Uses bcrypt with default cost for secure password storage
- **JWT Tokens**: Stateless authentication with configurable expiry
- **Strong Password Validation**: Enforces minimum security requirements for passwords
- **Request Validation**: Comprehensive input validation and sanitization
- **Security Headers**: Added security headers for XSS and clickjacking protection

## Environment Variables

### Required for Production
- `JWT_SECRET`: Secret key for signing JWT tokens (must be set in production)

### Optional Configuration
- `JWT_EXPIRY_HOURS`: Token expiry time in hours (default: 24)
- `DATABASE_URL`: PostgreSQL connection string
- `PORT`: Server port (default: 8080)
- `ENVIRONMENT`: Environment mode (development/production)

## Password Requirements

Passwords must meet the following criteria:
- Minimum 8 characters
- At least one uppercase letter
- At least one lowercase letter
- At least one digit
- At least one special character

## API Endpoints

### 1. Register User

**Endpoint**: `POST /api/v1/auth/register`

**Description**: Creates a new user account with the provided credentials.

**Request Body**:
```json
{
  "username": "string (3-50 characters)",
  "email": "string (valid email)",
  "password": "string (meets strong password requirements)",
  "firstName": "string (required)",
  "lastName": "string (required)"
}
```

**Response** (201 Created):
```json
{
  "token": "jwt_token_string",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "firstName": "Test",
    "lastName": "User",
    "role": "user",
    "avatar": null,
    "isActive": true,
    "createdAt": "2023-12-07T10:30:00Z",
    "updatedAt": "2023-12-07T10:30:00Z"
  }
}
```

**Error Responses**:
- `400 Bad Request`: Invalid input data
- `409 Conflict`: User already exists

### 2. Login User

**Endpoint**: `POST /api/v1/auth/login`

**Description**: Authenticates a user and returns a JWT token.

**Request Body**:
```json
{
  "username": "string",
  "password": "string"
}
```

**Response** (200 OK):
```json
{
  "token": "jwt_token_string",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "firstName": "Test",
    "lastName": "User",
    "role": "user",
    "avatar": null,
    "isActive": true,
    "createdAt": "2023-12-07T10:30:00Z",
    "updatedAt": "2023-12-07T10:30:00Z"
  }
}
```

**Error Responses**:
- `400 Bad Request`: Invalid input data
- `401 Unauthorized`: Invalid credentials

### 3. Get User Profile

**Endpoint**: `GET /api/v1/auth/profile`

**Description**: Retrieves the current user's profile information.

**Authentication**: Requires valid JWT token in `Authorization: Bearer <token>` header.

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "firstName": "Test",
    "lastName": "User",
    "role": "user",
    "avatar": null,
    "isActive": true,
    "createdAt": "2023-12-07T10:30:00Z",
    "updatedAt": "2023-12-07T10:30:00Z"
  }
}
```

**Error Responses**:
- `401 Unauthorized`: Invalid or missing token

### 4. Update User Profile

**Endpoint**: `PUT /api/v1/auth/profile`

**Description**: Updates the current user's profile information.

**Authentication**: Requires valid JWT token in `Authorization: Bearer <token>` header.

**Request Body**:
```json
{
  "firstName": "string (optional)",
  "lastName": "string (optional)",
  "avatar": "string (optional, URL)"
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "firstName": "Updated First Name",
    "lastName": "Updated Last Name",
    "role": "user",
    "avatar": "https://example.com/avatar.jpg",
    "isActive": true,
    "createdAt": "2023-12-07T10:30:00Z",
    "updatedAt": "2023-12-07T11:00:00Z"
  }
}
```

**Error Responses**:
- `400 Bad Request`: Invalid input data
- `401 Unauthorized`: Invalid or missing token

### 5. Logout User

**Endpoint**: `POST /api/v1/auth/logout`

**Description**: Logs out the current user. In a stateless JWT implementation, this is primarily handled client-side by removing the token.

**Authentication**: Requires valid JWT token in `Authorization: Bearer <token>` header.

**Response** (200 OK):
```json
{
  "success": true,
  "message": "Successfully logged out"
}
```

**Error Responses**:
- `401 Unauthorized`: Invalid or missing token

## Authentication Middleware

All protected routes use the `AuthMiddleware()` which:

1. Validates the JWT token from the `Authorization` header
2. Checks if the token is expired
3. Retrieves the user from the database
4. Sets user context for downstream handlers

## Authorization

### Role-Based Access Control

- **user**: Default role for regular users
- **admin**: Administrative access with additional permissions

### Admin Middleware

Protected admin routes use `AdminMiddleware()` which ensures the user has the "admin" role.

## Security Considerations

### Token Storage
- Tokens should be stored securely on the client side (e.g., httpOnly cookies or secure storage)
- Consider implementing refresh tokens for enhanced security

### Password Security
- All passwords are hashed using bcrypt before storage
- Strong password requirements are enforced during registration

### Session Management
- JWT tokens are stateless and contain all necessary user information
- Token expiry is configurable via environment variables

### Rate Limiting
- Basic rate limiting middleware is included (placeholder for production implementation)
- Consider implementing more sophisticated rate limiting in production

## Testing

### Unit Tests
- Authentication service tests: `internal/auth/auth_test.go`
- API handler tests: `internal/api/auth_test.go`

### Integration Tests
- End-to-end authentication flow tests: `tests/integration/api_test.go`

### Running Tests
```bash
# Run all backend tests
make test

# Run tests with coverage
make test-coverage

# Run integration tests
make test-integration
```

## Future Enhancements

1. **Refresh Token Strategy**: Implement refresh tokens for better security
2. **Token Blacklisting**: Add token invalidation for immediate logout
3. **Multi-Factor Authentication**: Add 2FA support
4. **OAuth Integration**: Support third-party authentication providers
5. **Rate Limiting**: Implement advanced rate limiting
6. **Account Lockout**: Add account lockout after failed attempts
7. **Password Reset**: Implement secure password reset functionality

## Error Response Format

All error responses follow this consistent format:

```json
{
  "success": false,
  "message": "Human-readable error message",
  "error": "Detailed error information (optional)"
}
```

## Success Response Format

Success responses follow this format:

```json
{
  "success": true,
  "message": "Success message (optional)",
  "data": {} // Response data (optional)
}
```