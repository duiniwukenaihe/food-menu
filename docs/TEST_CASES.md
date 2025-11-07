# Test Cases Documentation

This document provides a comprehensive overview of all test cases for the Full-Stack Application.

## Table of Contents
1. [Authentication Tests](#authentication-tests)
2. [Content Management Tests](#content-management-tests)
3. [User Management Tests](#user-management-tests)
4. [Admin Operations Tests](#admin-operations-tests)
5. [Recommendation System Tests](#recommendation-system-tests)
6. [API Integration Tests](#api-integration-tests)
7. [Frontend Component Tests](#frontend-component-tests)
8. [End-to-End Tests](#end-to-end-tests)

---

## Authentication Tests

### User Registration

| Test Case ID | Description | Test Data | Expected Result | Priority |
|--------------|-------------|------------|----------------|----------|
| AUTH-001 | Successful user registration | Valid: username, email, password, firstName, lastName | User created, JWT token returned, user redirected to dashboard | High |
| AUTH-002 | Registration with duplicate username | Existing username | Error message "User already exists" | High |
| AUTH-003 | Registration with duplicate email | Existing email | Error message "User already exists" | High |
| AUTH-004 | Registration with invalid email | Invalid email format | Validation error "Invalid email format" | Medium |
| AUTH-005 | Registration with weak password | Password < 6 characters | Validation error "Password must be at least 6 characters" | Medium |
| AUTH-006 | Registration with missing fields | Empty required fields | Validation errors for missing fields | High |
| AUTH-007 | Registration with SQL injection attempt | Malicious input in username/email | Input sanitized, no SQL errors | High |

### User Login

| Test Case ID | Description | Test Data | Expected Result | Priority |
|--------------|-------------|------------|----------------|----------|
| AUTH-008 | Successful login | Valid credentials | JWT token returned, user authenticated | High |
| AUTH-009 | Login with wrong password | Invalid password | Error message "Invalid credentials" | High |
| AUTH-010 | Login with non-existent username | Invalid username | Error message "Invalid credentials" | High |
| AUTH-011 | Login with inactive account | Inactive user credentials | Error message "Account is inactive" | Medium |
| AUTH-012 | Login with empty fields | Empty username/password | Validation errors | Medium |
| AUTH-013 | Login with SQL injection attempt | Malicious input | Input sanitized, authentication fails properly | High |

### Token Management

| Test Case ID | Description | Test Data | Expected Result | Priority |
|--------------|-------------|------------|----------------|----------|
| AUTH-014 | Valid JWT token | Properly formatted token | API access granted | High |
| AUTH-015 | Expired JWT token | Expired token | Error "Token expired", redirect to login | High |
| AUTH-016 | Invalid JWT token | Malformed token | Error "Invalid token", access denied | High |
| AUTH-017 | Missing JWT token | No Authorization header | Error "Authorization header required" | High |
| AUTH-018 | Token with wrong signature | Tampered token | Error "Invalid token" | High |

---

## Content Management Tests

### Content Creation

| Test Case ID | Description | Test Data | Expected Result | Priority |
|--------------|-------------|------------|----------------|----------|
| CONT-001 | Create published content | Valid content with isPublished=true | Content created and visible to all users | High |
| CONT-002 | Create draft content | Valid content with isPublished=false | Content created but not visible publicly | Medium |
| CONT-003 | Create content with missing title | Empty title | Validation error "Title is required" | High |
| CONT-004 | Create content with missing body | Empty body | Validation error "Body is required" | High |
| CONT-005 | Create content with invalid category | Non-existent category ID | Error "Category not found" | Medium |
| CONT-006 | Create content with HTML in body | HTML tags | HTML properly sanitized or escaped | High |
| CONT-007 | Create content with very long title | Title > 200 characters | Validation error "Title too long" | Medium |

### Content Retrieval

| Test Case ID | Description | Test Data | Expected Result | Priority |
|--------------|-------------|------------|----------------|----------|
| CONT-008 | Get all published content | No filters | List of published content with pagination | High |
| CONT-009 | Get content by category | Category filter | Content filtered by category | High |
| CONT-010 | Search content by title | Search term | Content matching search term returned | High |
| CONT-011 | Search content by description | Search term | Content matching description returned | Medium |
| CONT-012 | Pagination - first page | page=1, limit=10 | First 10 items returned | High |
| CONT-013 | Pagination - last page | page=last, limit=10 | Remaining items returned | Medium |
| CONT-014 | Get content by ID | Valid content ID | Specific content details returned | High |
| CONT-015 | Get non-existent content by ID | Invalid content ID | Error "Content not found" | Medium |
| CONT-016 | View count increment | Multiple views to same content | View count increases appropriately | High |

### Content Update

| Test Case ID | Description | Test Data | Expected Result | Priority |
|--------------|-------------|------------|----------------|----------|
| CONT-017 | Update content title | New title | Title updated successfully | High |
| CONT-018 | Update content body | New body content | Body updated successfully | High |
| CONT-019 | Publish draft content | Set isPublished=true | Content becomes visible publicly | High |
| CONT-020 | Unpublish content | Set isPublished=false | Content no longer visible publicly | Medium |
| CONT-021 | Update content category | New category ID | Category updated successfully | Medium |
| CONT-022 | Update non-existent content | Invalid content ID | Error "Content not found" | Medium |

### Content Deletion

| Test Case ID | Description | Test Data | Expected Result | Priority |
|--------------|-------------|------------|----------------|----------|
| CONT-023 | Delete content | Valid content ID | Content soft deleted, no longer visible | High |
| CONT-024 | Delete non-existent content | Invalid content ID | Error "Content not found" | Medium |
| CONT-025 | Delete already deleted content | Deleted content ID | Error "Content not found" | Low |

---

## User Management Tests

### User Profile

| Test Case ID | Description | Test Data | Expected Result | Priority |
|--------------|-------------|------------|----------------|----------|
| USER-001 | View user profile | Authenticated user | User profile data returned | High |
| USER-002 | Update user first name | New first name | First name updated successfully | High |
| USER-003 | Update user last name | New last name | Last name updated successfully | High |
| USER-004 | Update user avatar | New avatar URL | Avatar updated successfully | Medium |
| USER-005 | Update profile with invalid data | Invalid email format | Validation error | Medium |
| USER-006 | View profile without authentication | No token | Error "Authentication required" | High |

### User Administration

| Test Case ID | Description | Test Data | Expected Result | Priority |
|--------------|-------------|------------|----------------|----------|
| USER-007 | List all users (admin) | Admin user | Paginated list of all users | High |
| USER-008 | Create new user (admin) | Valid user data | User created successfully | High |
| USER-009 | Update user role (admin) | New role assignment | User role updated | High |
| USER-010 | Deactivate user (admin) | Set isActive=false | User account deactivated | High |
| USER-011 | Delete user (admin) | Valid user ID | User soft deleted | Medium |
| USER-012 | Access admin routes as regular user | Regular user token | Error "Admin access required" | High |

---

## Admin Operations Tests

### Category Management

| Test Case ID | Description | Test Data | Expected Result | Priority |
|--------------|-------------|------------|----------------|----------|
| ADMIN-001 | Create category | Valid category data | Category created successfully | High |
| ADMIN-002 | Create duplicate category | Existing category name | Error "Category already exists" | High |
| ADMIN-003 | Update category name | New category name | Name updated successfully | Medium |
| ADMIN-004 | Update category color | New color hex code | Color updated successfully | Low |
| ADMIN-005 | Delete category with content | Category has associated content | Error "Cannot delete category with content" | Medium |
| ADMIN-006 | List all categories | Admin user | All categories including inactive | High |

### Dashboard Analytics

| Test Case ID | Description | Test Data | Expected Result | Priority |
|--------------|-------------|------------|----------------|----------|
| ADMIN-007 | View user statistics | Admin user | User count, growth metrics displayed | Medium |
| ADMIN-008 | View content statistics | Admin user | Content count, view metrics displayed | Medium |
| ADMIN-009 | View system performance | Admin user | Performance metrics available | Low |
| ADMIN-010 | Export data reports | Admin user | CSV/PDF reports generated | Low |

---

## Recommendation System Tests

### Content Recommendations

| Test Case ID | Description | Test Data | Expected Result | Priority |
|--------------|-------------|------------|----------------|----------|
| REC-001 | Get recommendations for new user | New authenticated user | Popular content recommendations | High |
| REC-002 | Get recommendations for active user | User with view history | Personalized recommendations | High |
| REC-003 | Get recommendations with limit | Limit parameter | Limited number of recommendations | Medium |
| REC-004 | Mark recommendation as viewed | Recommendation ID | Recommendation marked as viewed | Medium |
| REC-005 | Get recommendations without authentication | No auth token | Error "Authentication required" | High |
| REC-006 | Recommendation algorithm accuracy | Test data | Relevant content prioritized | Medium |

---

## API Integration Tests

### Error Handling

| Test Case ID | Description | Test Data | Expected Result | Priority |
|--------------|-------------|------------|----------------|----------|
| API-001 | Handle malformed JSON | Invalid JSON in request body | Error "Invalid JSON format" | High |
| API-002 | Handle missing fields | Required field absent | Validation error | High |
| API-003 | Handle rate limiting | Excessive requests | Rate limit error response | Medium |
| API-004 | Handle database connection errors | Database unavailable | Error "Service unavailable" | High |
| API-005 | Handle CORS requests | Cross-origin request | Proper CORS headers | High |

### Response Format

| Test Case ID | Description | Test Data | Expected Result | Priority |
|--------------|-------------|------------|----------------|----------|
| API-006 | Success response format | Successful operation | Consistent success response structure | High |
| API-007 | Error response format | Failed operation | Consistent error response structure | High |
| API-008 | Pagination response | List endpoint | Proper pagination metadata | High |
| API-009 | Content-Type headers | All endpoints | Correct Content-Type headers | Medium |
| API-010 | HTTP status codes | Various operations | Appropriate HTTP status codes | High |

---

## Frontend Component Tests

### Layout Component

| Test Case ID | Description | Test Data | Expected Result | Priority |
|--------------|-------------|------------|----------------|----------|
| UI-001 | Render navigation bar | Default state | Navigation links displayed correctly | High |
| UI-002 | Show user menu when authenticated | Valid user token | User name and logout button visible | High |
| UI-003 | Show login/register when not authenticated | No token | Login and register links visible | High |
| UI-004 | Mobile menu toggle | Small screen | Mobile menu opens/closes properly | Medium |
| UI-005 | Admin link for admin users | Admin role user | Admin navigation link visible | High |

### Form Components

| Test Case ID | Description | Test Data | Expected Result | Priority |
|--------------|-------------|------------|----------------|----------|
| UI-006 | Login form validation | Empty fields | Validation errors displayed | High |
| UI-007 | Registration form validation | Invalid data | Appropriate validation messages | High |
| UI-008 | Form submission | Valid data | Form submits successfully | High |
| UI-009 | Loading state during submission | Form submission | Loading indicator displayed | Medium |
| UI-010 | Error message display | API error | Error messages shown to user | High |

### Content Display

| Test Case ID | Description | Test Data | Expected Result | Priority |
|--------------|-------------|------------|----------------|----------|
| UI-011 | Content list rendering | Content data | Content cards displayed correctly | High |
| UI-012 | Content detail view | Single content | Full content with metadata displayed | High |
| UI-013 | Search functionality | Search term | Results filtered properly | High |
| UI-014 | Category filtering | Category selection | Content filtered by category | High |
| UI-015 | Pagination navigation | Multiple pages | Page navigation works correctly | Medium |

---

## End-to-End Tests

### User Journey Tests

| Test Case ID | Description | Test Steps | Expected Result | Priority |
|--------------|-------------|------------|----------------|----------|
| E2E-001 | Complete user registration flow | 1. Visit register page<br>2. Fill valid form<br>3. Submit form<br>4. Verify dashboard | User registered and logged in successfully | High |
| E2E-002 | Complete user login flow | 1. Visit login page<br>2. Enter credentials<br>3. Submit form<br>4. Verify dashboard | User authenticated and redirected | High |
| E2E-003 | Content browsing journey | 1. Browse content list<br>2. Filter by category<br>3. Search content<br>4. View content details | All browsing features work correctly | High |
| E2E-004 | Admin user management | 1. Login as admin<br>2. Navigate to admin panel<br>3. Create new user<br>4. Verify user created | Admin operations work correctly | High |
| E2E-005 | Content creation workflow | 1. Login as admin<br>2. Create new content<br>3. Publish content<br>4. Verify content visible | Content creation and publishing works | High |

### Cross-Browser Tests

| Test Case ID | Description | Browsers | Expected Result | Priority |
|--------------|-------------|-----------|----------------|----------|
| E2E-006 | Chrome compatibility | Chrome latest | All features work correctly | High |
| E2E-007 | Firefox compatibility | Firefox latest | All features work correctly | High |
| E2E-008 | Safari compatibility | Safari latest | All features work correctly | Medium |
| E2E-009 | Mobile responsiveness | Mobile viewports | Responsive design works properly | High |

### Performance Tests

| Test Case ID | Description | Test Conditions | Expected Result | Priority |
|--------------|-------------|-----------------|----------------|----------|
| E2E-010 | Page load performance | Standard connection | Pages load within 3 seconds | High |
| E2E-011 | Large content list | 1000+ items | Pagination handles large datasets | Medium |
| E2E-012 | Concurrent users | 100 simultaneous users | System remains responsive | Medium |

---

## Test Data Examples

### User Test Data
```json
{
  "validUser": {
    "username": "testuser123",
    "email": "test@example.com",
    "password": "SecurePass123!",
    "firstName": "Test",
    "lastName": "User"
  },
  "adminUser": {
    "username": "admin",
    "password": "AdminPass123!",
    "role": "admin"
  }
}
```

### Content Test Data
```json
{
  "validContent": {
    "title": "Test Article Title",
    "description": "A test article description",
    "body": "This is the full body content of the test article with sufficient length.",
    "categoryId": 1,
    "tags": "test,article,example",
    "isPublished": true
  }
}
```

### Category Test Data
```json
{
  "validCategory": {
    "name": "Technology",
    "description": "Technology related content",
    "color": "#3B82F6",
    "icon": "tech"
  }
}
```

---

## Test Execution Priority

### Critical (High Priority)
- Authentication flows
- Authorization checks
- Core CRUD operations
- Data validation
- Security vulnerabilities

### Important (Medium Priority)
- UI/UX interactions
- Performance optimization
- Error handling
- Edge cases

### Nice to Have (Low Priority)
- Advanced features
- Optional functionality
- Minor UI improvements
- Administrative tools

---

## Test Environment Setup

### Required Services
- PostgreSQL database
- Redis (for caching)
- Email service (for notifications)
- File storage (for images)

### Test Data Seeding
- Default admin user
- Sample categories
- Test content items
- Sample regular users

### Configuration
- Test database connection
- JWT secret keys
- CORS settings
- Rate limiting configuration

---

## Automation Strategy

### Continuous Integration
- Run unit tests on every commit
- Run integration tests on pull requests
- Run E2E tests on main branch merges
- Generate coverage reports

### Test Scheduling
- Daily full test suite execution
- Weekly performance testing
- Monthly security scanning
- Quarterly load testing

### Reporting
- Test execution dashboards
- Coverage trend analysis
- Performance metrics tracking
- Bug detection rates