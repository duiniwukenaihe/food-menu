# Admin Dishes and Media Management API

This document describes the admin-only endpoints for managing dishes and media uploads.

## Authentication

All admin endpoints require:
1. Valid JWT token in Authorization header: `Bearer <token>`
2. User account with `admin` role

## Dish Management

### List All Dishes (Admin)

**GET** `/api/v1/admin/dishes`

Lists all dishes including inactive ones with pagination and filtering.

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `limit` (integer, optional): Items per page (default: 10)
- `search` (string, optional): Search by name or description
- `isSeasonal` (boolean, optional): Filter by seasonal dishes
- `isActive` (boolean, optional): Filter by active status
- `tags` (string, optional): Filter by tags (comma-separated)

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Summer Salad",
      "description": "Fresh seasonal salad",
      "tags": "salad,vegetarian,healthy",
      "isActive": true,
      "isSeasonal": true,
      "availableMonths": "6,7,8",
      "seasonalNote": "Available during summer months",
      "imageUrl": "http://localhost:8080/api/v1/media/dishes/2024/01/abc123.jpg",
      "thumbnailUrl": "http://localhost:8080/api/v1/media/dishes/2024/01/abc123_thumb.jpg",
      "galleryUrls": "[\"url1\", \"url2\"]",
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 50,
  "page": 1,
  "limit": 10
}
```

### Create Dish

**POST** `/api/v1/admin/dishes`

Creates a new dish.

**Request Body:**
```json
{
  "name": "Summer Salad",
  "description": "Fresh seasonal salad",
  "tags": "salad,vegetarian,healthy",
  "isActive": true,
  "isSeasonal": true,
  "availableMonths": "6,7,8",
  "seasonalNote": "Available during summer months",
  "imageUrl": "http://localhost:8080/api/v1/media/dishes/2024/01/abc123.jpg",
  "thumbnailUrl": "http://localhost:8080/api/v1/media/dishes/2024/01/abc123_thumb.jpg",
  "galleryUrls": "[\"url1\", \"url2\"]"
}
```

**Validation:**
- `name`: Required, 1-255 characters
- All fields except `name` are optional

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Dish created successfully",
  "data": {
    "id": 1,
    "name": "Summer Salad",
    ...
  }
}
```

### Update Dish

**PUT** `/api/v1/admin/dishes/:id`

Updates an existing dish. Only provided fields will be updated.

**Request Body:**
```json
{
  "name": "Updated Summer Salad",
  "imageUrl": "http://localhost:8080/api/v1/media/dishes/2024/01/new-image.jpg"
}
```

**Validation:**
- `name`: If provided, 1-255 characters
- When updating media URLs, old media is automatically cleaned up

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Dish updated successfully",
  "data": {
    "id": 1,
    "name": "Updated Summer Salad",
    ...
  }
}
```

### Delete Dish

**DELETE** `/api/v1/admin/dishes/:id`

Deletes a dish and cleans up associated media files.

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Dish deleted successfully"
}
```

## Media Management

### Get Upload URL

**POST** `/api/v1/admin/media/upload-url`

Generates a presigned URL for direct file upload to object storage.

**Request Body:**
```json
{
  "fileName": "dish-image.jpg",
  "contentType": "image/jpeg"
}
```

**Validation:**
- `fileName`: Required, max 255 characters
- `contentType`: Must be allowed image or video type
  - Images: `image/jpeg`, `image/jpg`, `image/png`, `image/gif`, `image/webp`
  - Videos: `video/mp4`, `video/webm`, `video/ogg`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "uploadUrl": "http://localhost:8080/api/v1/admin/media/upload?key=dishes/2024/01/abc123.jpg",
    "fileUrl": "http://localhost:8080/api/v1/media/dishes/2024/01/abc123.jpg",
    "key": "dishes/2024/01/abc123.jpg"
  }
}
```

**Usage:**
1. Get presigned URL from this endpoint
2. Upload file to `uploadUrl` using PUT or POST
3. Use `fileUrl` in dish creation/update requests

### Upload File

**POST** `/api/v1/admin/media/upload`

Uploads a file directly through the API.

**Request:**
- Content-Type: `multipart/form-data`
- Form field `file`: The file to upload
- Query parameter `key` (optional): Custom file key

**Validation:**
- File size limits:
  - Images: 10MB
  - Videos: 100MB
- Content type: Must be allowed image or video type
- File is saved to storage and metadata is recorded in database

**Response (200 OK):**
```json
{
  "success": true,
  "message": "File uploaded successfully",
  "data": {
    "url": "http://localhost:8080/api/v1/media/dishes/2024/01/abc123.jpg",
    "fileName": "dish-image.jpg",
    "contentType": "image/jpeg",
    "size": 102400
  }
}
```

### Delete Media

**DELETE** `/api/v1/admin/media`

Deletes an uploaded media file and its database record.

**Request Body:**
```json
{
  "key": "dishes/2024/01/abc123.jpg"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "File deleted successfully"
}
```

## Public Endpoints

### Get Dishes (Public)

**GET** `/api/v1/dishes`

Lists active dishes only.

### Get Dish by ID (Public)

**GET** `/api/v1/dishes/:id`

Gets a single active dish.

### Get Media File

**GET** `/api/v1/media/*filepath`

Serves uploaded media files.

## Storage Configuration

The system supports two storage backends:

### MinIO (S3-Compatible)

Set environment variables:
```bash
STORAGE_TYPE=minio
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=media
MINIO_USE_SSL=false
```

### Local File System

Set environment variables:
```bash
STORAGE_TYPE=local
UPLOAD_DIR=./uploads
API_BASE_URL=http://localhost:8080
```

## Error Responses

All endpoints return consistent error responses:

**400 Bad Request:**
```json
{
  "success": false,
  "message": "Invalid request body",
  "error": "detailed error message"
}
```

**401 Unauthorized:**
```json
{
  "success": false,
  "message": "Authorization header required"
}
```

**403 Forbidden:**
```json
{
  "success": false,
  "message": "Admin access required"
}
```

**404 Not Found:**
```json
{
  "success": false,
  "message": "Dish not found"
}
```

**500 Internal Server Error:**
```json
{
  "success": false,
  "message": "Failed to create dish",
  "error": "detailed error message"
}
```

## Testing

### Running Integration Tests

```bash
# Run all integration tests
cd backend
go test -tags=integration ./tests/integration/...

# Run specific test file
go test -tags=integration ./tests/integration/ -run TestAdminDishCRUD

# With verbose output
go test -v -tags=integration ./tests/integration/...
```

### Test Environment Setup

Tests require:
1. PostgreSQL database (connection via DATABASE_URL env var)
2. Optional: MinIO instance (set MINIO_ENDPOINT for storage tests)

The test suite will:
- Create test users with admin privileges
- Test CRUD operations on dishes
- Test media upload and deletion
- Test validation and authorization
- Clean up test data automatically

## Media Cleanup

The system automatically cleans up media files when:
1. A dish is deleted (removes associated image, thumbnail)
2. A dish's media URL is updated (removes old media file)
3. Media is explicitly deleted via DELETE endpoint

This prevents orphaned files and conserves storage space.

## Best Practices

1. **Upload Flow:**
   - Upload media first using `/admin/media/upload`
   - Get the returned URL
   - Use URL in dish creation/update

2. **Image Optimization:**
   - Upload appropriately sized images
   - Use thumbnails for list views
   - Full images for detail views

3. **Validation:**
   - Validate file types on client side first
   - Check file sizes before upload
   - Handle validation errors gracefully

4. **Security:**
   - Always verify admin role before accessing endpoints
   - Prevent directory traversal in file paths
   - Sanitize file names

5. **Performance:**
   - Use pagination for large lists
   - Implement caching for media files
   - Consider CDN for production media delivery
