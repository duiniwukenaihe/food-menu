# What's New

## Latest Updates - Food Ordering System

### 🎉 New Features

#### 1. Media Upload System
Upload dish images and videos directly to S3 storage through the API.

**Endpoint:** `POST /api/v1/admin/upload/media`

**Benefits:**
- Centralized media management
- Automatic S3 storage integration
- Supports multiple cloud providers (AWS S3, Aliyun OSS, Tencent COS, MinIO)
- Built-in file validation and size limits

**Quick Start:**
```bash
# Upload an image
curl -X POST http://localhost:8080/api/v1/admin/upload/media \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@dish.jpg"
```

See: [API Documentation - Media Upload](API.md#媒体上传)

---

#### 2. Random Combo Generator
Automatically generate balanced dish combinations based on meat and vegetable preferences.

**Endpoint:** `POST /api/v1/combo/generate`

**Features:**
- Configurable meat/vegetable ratios
- Exclude specific dishes
- Returns total price calculation
- Smart category-based selection

**Quick Start:**
```bash
# Generate a combo: 1 meat + 2 vegetables
curl -X POST http://localhost:8080/api/v1/combo/generate \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"meat_count":1,"vegetable_count":2}'
```

See: [API Documentation - Random Combo](API.md#随机搭配)

---

#### 3. Public Configuration API
Access system configuration without authentication for frontend initialization.

**Endpoint:** `GET /api/v1/config/public`

**Features:**
- No authentication required
- Returns public settings only
- Feature flags support
- Safe for client-side use

**Quick Start:**
```bash
# Get public config
curl http://localhost:8080/api/v1/config/public
```

See: [API Documentation - Public Config](API.md#公共配置)

---

#### 4. API Testing Tool
Built-in web interface for testing all API endpoints.

**Access:** `http://localhost:8080/api-tester.html`

**Features:**
- Visual interface for all endpoints
- Automatic JWT token management
- File upload testing support
- Request/response inspection
- Save frequent requests

**Quick Start:**
1. Start the backend server
2. Visit http://localhost:8080/api-tester.html
3. Login with admin/admin123
4. Test any endpoint

See: [README - API Testing Tool](../README.md#-api测试工具)

---

### 🔧 Configuration Enhancements

#### New Environment Variables

**JWT_EXPIRATION**
```bash
JWT_EXPIRATION=24h  # Token validity period
```
Supported formats: `24h`, `30m`, `3600s`

**S3_PATH_STYLE**
```bash
S3_PATH_STYLE=false  # Set to true for MinIO
```
Required for MinIO compatibility.

**TEST_DATABASE_URL**
```bash
TEST_DATABASE_URL=postgres://user:pass@localhost/test_db?sslmode=disable
```
Separate database for running tests.

See: [README - Configuration](../README.md#-配置说明)

---

### 📚 Documentation Improvements

#### Comprehensive Troubleshooting Guide
New troubleshooting section covering:
- S3 storage issues (connection, authentication, permissions)
- Database migration problems
- JWT authentication errors
- CORS and frontend issues
- Common startup errors

See: [README - Troubleshooting](../README.md#-故障排除)

#### Migration Instructions
Detailed database setup and migration guide:
- Initial database creation
- Running migrations manually
- Test database setup
- Migration file organization

See: [README - Database Migration](../README.md#-数据库迁移)

#### S3 Provider Setup Guides
Step-by-step instructions for:
- AWS S3
- Aliyun OSS
- Tencent COS
- MinIO (local development)

See: [README - S3 Configuration](../README.md#s3对象存储配置)

---

### 🔄 Workflow Changes

#### ✅ New Recommended Workflow

**Creating a dish with image:**
1. Upload image via `/admin/upload/media`
2. Receive S3 URL in response
3. Create dish with the returned URL
4. Image is stored securely in S3

**Old workflow (deprecated):**
~~Manually enter external image URLs~~

**Benefits:**
- Centralized storage
- Better control and security
- Consistent file naming
- Automatic cleanup possible
- CDN integration ready

---

### 📖 Quick Links

- [Complete API Documentation](API.md)
- [Setup Guide](../README.md#-快速开始)
- [Troubleshooting](../README.md#-故障排除)
- [Configuration Reference](../README.md#-配置说明)
- [Deployment Guide](DEPLOYMENT.md)
- [Changelog](../CHANGELOG.md)

---

### 🆘 Need Help?

**Common Issues:**

1. **S3 upload fails**
   - Check S3 credentials in `.env`
   - Verify bucket permissions
   - See [Troubleshooting - S3](../README.md#s3存储相关问题)

2. **Database connection error**
   - Verify PostgreSQL is running
   - Check DATABASE_URL format
   - See [Troubleshooting - Database](../README.md#数据库迁移问题)

3. **Token expired quickly**
   - Adjust JWT_EXPIRATION value
   - See [Troubleshooting - JWT](../README.md#jwt认证问题)

**Still stuck?**
- Check the [Troubleshooting Guide](../README.md#-故障排除)
- Review [API Documentation](API.md)
- Open an issue on GitHub

---

### 🎓 Learning Resources

**For Developers:**
- [API Quick Reference](API.md#快速参考)
- [Common Development Scenarios](API.md#常见开发场景)
- [Environment Variables List](API.md#环境变量完整列表)

**For Administrators:**
- [S3 Setup Guide](../README.md#s3对象存储配置)
- [Migration Guide](../README.md#-数据库迁移)
- [System Configuration](API.md#管理员接口)

**For Users:**
- [Quick Start Guide](../README.md#-快速开始)
- [Feature Overview](../README.md#-功能特性)
- [Demo Accounts](../README.md#-测试账号)
