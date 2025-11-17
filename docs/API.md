# 食物点餐系统 API 文档

## 概述

本文档描述了食物点餐系统的后端API接口。所有API都基于RESTful风格，使用JSON格式进行数据交换。

**基础URL:** `http://localhost:8080/api/v1`

**认证方式:** Bearer Token (JWT)

## 认证

### 用户登录

**POST** `/login`

登录系统获取访问令牌。

**请求体:**
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**响应:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "role": "admin",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

### 获取用户信息

**GET** `/profile`

获取当前登录用户的详细信息。

**Headers:**
```
Authorization: Bearer {token}
```

**响应:**
```json
{
  "id": 1,
  "username": "admin",
  "email": "admin@example.com",
  "role": "admin",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

## 菜品管理

### 获取菜品列表

**GET** `/dishes`

获取菜品列表，支持分页和筛选。

**查询参数:**
- `page` (int, optional): 页码，默认1
- `limit` (int, optional): 每页数量，默认20
- `category_id` (int, optional): 分类ID
- `search` (string, optional): 搜索关键词

**响应:**
```json
{
  "dishes": [
    {
      "id": 1,
      "name": "宫保鸡丁",
      "description": "经典川菜，麻辣鲜香",
      "category_id": 1,
      "category": {
        "id": 1,
        "name": "肉类"
      },
      "price": 28.00,
      "image_url": "https://example.com/dish1.jpg",
      "video_url": "https://example.com/dish1.mp4",
      "cooking_steps": "1. 切鸡丁\n2. 准备配料\n3. 爆炒",
      "is_seasonal": false,
      "is_active": true,
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  ],
  "total": 50,
  "page": 1,
  "limit": 20
}
```

### 获取单个菜品

**GET** `/dishes/{id}`

获取指定菜品的详细信息。

**路径参数:**
- `id` (int): 菜品ID

**响应:**
```json
{
  "id": 1,
  "name": "宫保鸡丁",
  "description": "经典川菜，麻辣鲜香",
  "category_id": 1,
  "category": {
    "id": 1,
    "name": "肉类"
  },
  "price": 28.00,
  "image_url": "https://example.com/dish1.jpg",
  "video_url": "https://example.com/dish1.mp4",
  "cooking_steps": "1. 切鸡丁\n2. 准备配料\n3. 爆炒",
  "is_seasonal": false,
  "is_active": true,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### 创建菜品 (管理员)

**POST** `/admin/dishes`

创建新菜品。

**Headers:**
```
Authorization: Bearer {admin_token}
```

**请求体:**
```json
{
  "name": "新菜品",
  "description": "菜品描述",
  "category_id": 1,
  "price": 25.00,
  "image_url": "https://your-bucket.s3.region.amazonaws.com/dishes/2023/12/image-uuid.jpg",
  "video_url": "https://your-bucket.s3.region.amazonaws.com/dishes/2023/12/video-uuid.mp4",
  "cooking_steps": "制作步骤",
  "is_seasonal": false
}
```

**注意：**
- `image_url` 和 `video_url` 应该使用 `/admin/upload/media` 接口上传文件后返回的URL
- 推荐工作流程：先上传媒体文件获取URL，再创建或更新菜品

### 更新菜品 (管理员)

**PUT** `/admin/dishes/{id}`

更新指定菜品信息。

**Headers:**
```
Authorization: Bearer {admin_token}
```

**请求体:**
```json
{
  "name": "更新的菜品名",
  "price": 30.00,
  "is_active": false
}
```

### 删除菜品 (管理员)

**DELETE** `/admin/dishes/{id}`

软删除指定菜品（设置为不活跃状态）。

**Headers:**
```
Authorization: Bearer {admin_token}
```

## 分类管理

### 获取分类列表

**GET** `/categories`

获取所有菜品分类。

**响应:**
```json
[
  {
    "id": 1,
    "name": "肉类",
    "description": "各种肉类菜品",
    "created_at": "2023-01-01T00:00:00Z"
  }
]
```

### 创建分类 (管理员)

**POST** `/admin/categories`

创建新的菜品分类。

**Headers:**
```
Authorization: Bearer {admin_token}
```

**请求体:**
```json
{
  "name": "新分类",
  "description": "分类描述"
}
```

## 推荐管理

### 获取推荐配置

**GET** `/recommendations`

获取所有推荐配置。

**响应:**
```json
[
  {
    "id": 1,
    "name": "经典搭配",
    "description": "一荤两素的经典搭配",
    "meat_count": 1,
    "vegetable_count": 2,
    "is_active": true,
    "created_at": "2023-01-01T00:00:00Z"
  }
]
```

### 获取应季菜品

**GET** `/seasonal-dishes`

获取应季推荐菜品列表。

**响应:**
```json
[
  {
    "id": 1,
    "name": "春季时蔬",
    "description": "春季新鲜蔬菜",
    "price": 18.00,
    "image_url": "https://example.com/seasonal.jpg",
    "is_seasonal": true,
    "created_at": "2023-01-01T00:00:00Z"
  }
]
```

## 订单管理

### 创建订单

**POST** `/orders`

创建新订单。

**Headers:**
```
Authorization: Bearer {token}
```

**请求体:**
```json
{
  "items": [
    {
      "dish_id": 1,
      "quantity": 2
    },
    {
      "dish_id": 2,
      "quantity": 1
    }
  ]
}
```

**响应:**
```json
{
  "id": 1,
  "user_id": 1,
  "total_amount": 74.00,
  "status": "pending",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z",
  "items": [
    {
      "id": 1,
      "order_id": 1,
      "dish_id": 1,
      "dish": {
        "id": 1,
        "name": "宫保鸡丁",
        "price": 28.00
      },
      "quantity": 2,
      "price": 28.00,
      "created_at": "2023-01-01T00:00:00Z"
    }
  ]
}
```

### 获取用户订单

**GET** `/orders`

获取当前用户的订单列表。

**Headers:**
```
Authorization: Bearer {token}
```

**查询参数:**
- `page` (int, optional): 页码，默认1
- `limit` (int, optional): 每页数量，默认20

**响应:**
```json
{
  "orders": [
    {
      "id": 1,
      "user_id": 1,
      "total_amount": 74.00,
      "status": "pending",
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  ],
  "total": 5,
  "page": 1,
  "limit": 20
}
```

## 收藏管理

### 添加到收藏

**POST** `/favorites/{dishId}`

将菜品添加到用户收藏。

**Headers:**
```
Authorization: Bearer {token}
```

**路径参数:**
- `dishId` (int): 菜品ID

### 从收藏中移除

**DELETE** `/favorites/{dishId}`

将菜品从用户收藏中移除。

**Headers:**
```
Authorization: Bearer {token}
```

**路径参数:**
- `dishId` (int): 菜品ID

### 获取收藏列表

**GET** `/favorites`

获取用户收藏的菜品列表。

**Headers:**
```
Authorization: Bearer {token}
```

**查询参数:**
- `page` (int, optional): 页码，默认1
- `limit` (int, optional): 每页数量，默认20

**响应:**
```json
{
  "favorites": [
    {
      "id": 1,
      "user_id": 1,
      "dish_id": 1,
      "dish": {
        "id": 1,
        "name": "宫保鸡丁",
        "price": 28.00,
        "image_url": "https://example.com/dish1.jpg"
      },
      "created_at": "2023-01-01T00:00:00Z"
    }
  ],
  "total": 10,
  "page": 1,
  "limit": 20
}
```

## 媒体上传

### 上传图片或视频 (管理员)

**POST** `/admin/upload/media`

上传菜品图片或视频到S3存储。支持的格式：图片(jpg, jpeg, png, gif, webp)，视频(mp4, avi, mov, webm)。

**Headers:**
```
Authorization: Bearer {admin_token}
Content-Type: multipart/form-data
```

**请求体:**
- `file` (file): 要上传的文件
- `type` (string, optional): 文件类型，可选值：image, video，默认自动检测

**请求示例（使用 cURL）:**
```bash
curl -X POST http://localhost:8080/api/v1/admin/upload/media \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -F "file=@/path/to/image.jpg" \
  -F "type=image"
```

**成功响应:**
```json
{
  "url": "https://your-bucket.s3.region.amazonaws.com/dishes/2023/12/image-uuid.jpg",
  "filename": "image-uuid.jpg",
  "size": 245678,
  "content_type": "image/jpeg"
}
```

**错误响应:**

文件过大 (413):
```json
{
  "error": "File size exceeds maximum limit of 10MB"
}
```

不支持的文件类型 (400):
```json
{
  "error": "Unsupported file type. Allowed types: jpg, jpeg, png, gif, webp, mp4, avi, mov, webm"
}
```

S3配置错误 (500):
```json
{
  "error": "S3 storage is not configured. Please set S3_ENDPOINT, S3_BUCKET, and credentials"
}
```

上传失败 (500):
```json
{
  "error": "Failed to upload file to S3: connection timeout"
}
```

**使用说明:**
1. 上传成功后，返回的 URL 可直接用于创建或更新菜品的 `image_url` 或 `video_url` 字段
2. 最大文件大小：10MB (图片)，50MB (视频)
3. 需要先配置 S3 存储凭据，详见 [S3配置说明](#s3存储配置)

**推荐工作流程：**

创建新菜品时：
1. 使用 `POST /admin/upload/media` 上传图片，获取 `image_url`
2. 使用 `POST /admin/upload/media` 上传视频（可选），获取 `video_url`
3. 使用 `POST /admin/dishes` 创建菜品，传入上述URL

更新菜品媒体时：
1. 使用 `POST /admin/upload/media` 上传新的图片或视频
2. 使用 `PUT /admin/dishes/{id}` 更新菜品的 `image_url` 或 `video_url`

**已弃用的方式：**
- ❌ 不再支持手动输入外部图片URL（不推荐）
- ✅ 所有媒体文件应通过上传接口上传到S3存储

## 随机搭配

### 生成随机菜品组合

**POST** `/combo/generate`

根据配置生成随机菜品组合，智能匹配荤素搭配。

**Headers:**
```
Authorization: Bearer {token}
```

**请求体:**
```json
{
  "meat_count": 1,
  "vegetable_count": 2,
  "exclude_dish_ids": [1, 5, 8]
}
```

**参数说明:**
- `meat_count` (int, optional): 荤菜数量，默认1
- `vegetable_count` (int, optional): 素菜数量，默认2
- `exclude_dish_ids` (array, optional): 要排除的菜品ID列表

**成功响应:**
```json
{
  "combo_id": "combo-uuid-12345",
  "dishes": [
    {
      "id": 1,
      "name": "宫保鸡丁",
      "description": "经典川菜，麻辣鲜香",
      "category_id": 1,
      "category": {
        "id": 1,
        "name": "肉类"
      },
      "price": 28.00,
      "image_url": "https://example.com/dish1.jpg",
      "is_seasonal": false
    },
    {
      "id": 12,
      "name": "清炒时蔬",
      "price": 15.00,
      "category": {
        "id": 2,
        "name": "蔬菜类"
      }
    }
  ],
  "total_price": 58.00,
  "generated_at": "2023-12-01T10:30:00Z"
}
```

**错误响应:**

参数无效 (400):
```json
{
  "error": "Invalid parameter: meat_count must be between 0 and 5"
}
```

可用菜品不足 (404):
```json
{
  "error": "Not enough dishes available to generate combo with requested criteria"
}
```

## 公共配置

### 获取公共配置

**GET** `/config/public`

获取系统公共配置信息，无需认证。

**响应:**
```json
{
  "default_meat_count": 1,
  "default_vegetable_count": 2,
  "max_dish_count": 6,
  "system_name": "食物点餐系统",
  "features": {
    "seasonal_dishes": true,
    "combo_generation": true,
    "favorites": true
  }
}
```

**说明:**
- 此接口不需要认证，可用于前端初始化配置
- 不包含敏感信息（如S3凭据、数据库连接等）
- 可用于动态调整前端功能显示

## 管理员接口

### 获取用户列表 (管理员)

**GET** `/admin/users`

获取系统所有用户列表。

**Headers:**
```
Authorization: Bearer {admin_token}
```

**查询参数:**
- `page` (int, optional): 页码，默认1
- `limit` (int, optional): 每页数量，默认20
- `search` (string, optional): 搜索关键词

**响应:**
```json
{
  "users": [
    {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "role": "admin",
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  ],
  "total": 10,
  "page": 1,
  "limit": 20
}
```

### 获取系统配置 (管理员)

**GET** `/admin/config`

获取系统配置信息。

**Headers:**
```
Authorization: Bearer {admin_token}
```

**响应:**
```json
[
  {
    "id": 1,
    "config_key": "default_meat_count",
    "config_value": "1",
    "description": "默认荤菜数量",
    "updated_at": "2023-01-01T00:00:00Z"
  },
  {
    "id": 2,
    "config_key": "s3_endpoint",
    "config_value": "https://s3.amazonaws.com",
    "description": "S3端点地址",
    "updated_at": "2023-01-01T00:00:00Z"
  }
]
```

### 更新系统配置 (管理员)

**PUT** `/admin/config`

更新系统配置。

**Headers:**
```
Authorization: Bearer {admin_token}
```

**请求体:**
```json
{
  "default_meat_count": "2",
  "default_vegetable_count": "2",
  "max_dish_count": "8"
}
```

**成功响应:**
```json
{
  "message": "Config updated successfully"
}
```

## S3存储配置

系统支持多种S3兼容存储服务，用于存储菜品图片和视频。

### 配置环境变量

在 `backend/.env` 文件中配置以下变量：

```bash
# S3存储配置
S3_ENDPOINT=https://s3.us-west-2.amazonaws.com  # S3端点地址
S3_ACCESS_KEY=your_access_key_id                # 访问密钥ID
S3_SECRET_KEY=your_secret_access_key            # 密钥
S3_BUCKET=your-bucket-name                      # 存储桶名称
S3_REGION=us-west-2                             # 区域
S3_PATH_STYLE=false                             # 路径样式（MinIO需要设为true）
```

### 支持的存储服务

**AWS S3:**
```bash
S3_ENDPOINT=https://s3.us-west-2.amazonaws.com
S3_REGION=us-west-2
S3_PATH_STYLE=false
```

**阿里云 OSS:**
```bash
S3_ENDPOINT=https://oss-cn-beijing.aliyuncs.com
S3_REGION=oss-cn-beijing
S3_PATH_STYLE=false
```

**腾讯云 COS:**
```bash
S3_ENDPOINT=https://cos.ap-beijing.myqcloud.com
S3_REGION=ap-beijing
S3_PATH_STYLE=false
```

**MinIO:**
```bash
S3_ENDPOINT=http://localhost:9000
S3_REGION=us-east-1
S3_PATH_STYLE=true  # MinIO需要使用路径样式
```

### 测试S3连接

使用以下命令测试S3配置是否正确：

```bash
cd backend
go run main.go --test-s3
```

## JWT认证配置

系统使用JWT (JSON Web Token) 进行用户认证。

### 环境变量

```bash
# JWT配置
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRATION=24h  # Token过期时间：24小时
```

### Token过期时间格式

支持的时间单位：
- `h` - 小时（如：24h）
- `m` - 分钟（如：30m）
- `s` - 秒（如：3600s）

示例：
- `24h` - 24小时
- `7d` - 7天（使用 168h）
- `30m` - 30分钟

### 使用JWT Token

在每个需要认证的请求中，在Header中添加：

```
Authorization: Bearer YOUR_JWT_TOKEN
```

## API测试工具

### Web测试页面

访问 `http://localhost:8080/api-tester.html` 可以使用内置的API测试工具。

功能特性：
- 可视化测试所有API接口
- 自动管理JWT Token
- 支持文件上传测试
- 显示请求/响应详情
- 保存常用请求

### 使用步骤

1. 启动后端服务
2. 浏览器访问 http://localhost:8080/api-tester.html
3. 使用测试账号登录获取Token
4. 选择要测试的接口进行调用

## 错误响应

所有API在出错时都会返回统一的错误格式：

```json
{
  "error": "错误描述信息"
}
```

常见HTTP状态码：
- `200` - 成功
- `201` - 创建成功
- `400` - 请求参数错误
- `401` - 未授权（未登录或Token过期）
- `403` - 权限不足（非管理员）
- `404` - 资源不存在
- `413` - 请求体过大（文件上传超限）
- `500` - 服务器内部错误

### 常见错误及解决方案

**Token过期 (401):**
```json
{
  "error": "Token has expired"
}
```
解决方案：重新登录获取新Token

**权限不足 (403):**
```json
{
  "error": "Admin access required"
}
```
解决方案：使用管理员账号登录

**S3存储未配置 (500):**
```json
{
  "error": "S3 storage is not configured"
}
```
解决方案：检查 `.env` 文件中的S3配置

## 测试账号

系统提供以下测试账号：

**管理员账号:**
- 用户名: `admin`
- 密码: `admin123`

**普通用户账号:**
- 用户名: `user`
- 密码: `user123`

## 使用示例

### JavaScript/TypeScript 示例

```javascript
// 登录
const loginResponse = await fetch('/api/v1/login', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    username: 'admin',
    password: 'admin123'
  })
})
const { token } = await loginResponse.json()

// 获取菜品列表
const dishesResponse = await fetch('/api/v1/dishes?page=1&limit=10', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
})
const dishesData = await dishesResponse.json()
```

### cURL 示例

```bash
# 登录
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 获取菜品列表
curl -X GET "http://localhost:8080/api/v1/dishes?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# 上传媒体文件
curl -X POST http://localhost:8080/api/v1/admin/upload/media \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -F "file=@/path/to/image.jpg" \
  -F "type=image"

# 生成随机搭配
curl -X POST http://localhost:8080/api/v1/combo/generate \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"meat_count":1,"vegetable_count":2}'

# 获取公共配置
curl -X GET http://localhost:8080/api/v1/config/public
```

## 快速参考

### API端点总览

#### 公开接口（无需认证）
| 方法 | 端点 | 说明 |
|------|------|------|
| POST | `/api/v1/login` | 用户登录 |
| GET | `/api/v1/dishes` | 获取菜品列表 |
| GET | `/api/v1/dishes/:id` | 获取单个菜品 |
| GET | `/api/v1/categories` | 获取分类列表 |
| GET | `/api/v1/recommendations` | 获取推荐配置 |
| GET | `/api/v1/seasonal-dishes` | 获取应季菜品 |
| GET | `/api/v1/config/public` | 获取公共配置（新） |

#### 认证接口（需要登录）
| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/api/v1/profile` | 获取用户信息 |
| POST | `/api/v1/orders` | 创建订单 |
| GET | `/api/v1/orders` | 获取订单列表 |
| POST | `/api/v1/favorites/:dishId` | 添加收藏 |
| DELETE | `/api/v1/favorites/:dishId` | 取消收藏 |
| GET | `/api/v1/favorites` | 获取收藏列表 |
| POST | `/api/v1/combo/generate` | 生成随机搭配（新） |

#### 管理员接口（需要管理员权限）
| 方法 | 端点 | 说明 |
|------|------|------|
| POST | `/api/v1/admin/upload/media` | **上传媒体文件（新）** |
| GET | `/api/v1/admin/users` | 获取用户列表 |
| POST | `/api/v1/admin/dishes` | 创建菜品 |
| PUT | `/api/v1/admin/dishes/:id` | 更新菜品 |
| DELETE | `/api/v1/admin/dishes/:id` | 删除菜品 |
| POST | `/api/v1/admin/categories` | 创建分类 |
| PUT | `/api/v1/admin/categories/:id` | 更新分类 |
| DELETE | `/api/v1/admin/categories/:id` | 删除分类 |
| GET | `/api/v1/admin/config` | 获取系统配置 |
| PUT | `/api/v1/admin/config` | 更新系统配置 |

### 新增功能说明

#### 1. 媒体上传 `/admin/upload/media`
- **用途：** 上传菜品图片和视频到S3存储
- **要求：** 管理员权限 + S3配置完成
- **返回：** S3文件URL
- **流程：** 上传文件 → 获取URL → 创建/更新菜品

#### 2. 随机搭配 `/combo/generate`
- **用途：** 根据荤素配置生成随机菜品组合
- **参数：** 荤菜数量、素菜数量、排除列表
- **返回：** 随机选择的菜品列表和总价

#### 3. 公共配置 `/config/public`
- **用途：** 获取系统公共配置（无需认证）
- **内容：** 默认参数、功能开关等
- **用于：** 前端初始化配置

### 环境变量完整列表

```bash
# 服务器
PORT=8080

# 数据库
DATABASE_URL=postgres://postgres:password@localhost/food_ordering?sslmode=disable
TEST_DATABASE_URL=postgres://postgres:password@localhost/food_ordering_test?sslmode=disable

# JWT认证
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24h

# S3存储（媒体上传必需）
S3_ENDPOINT=https://s3.amazonaws.com
S3_ACCESS_KEY=your-access-key
S3_SECRET_KEY=your-secret-key
S3_BUCKET=your-bucket
S3_REGION=us-west-1
S3_PATH_STYLE=false
```

### 常见开发场景

#### 场景1：管理员添加新菜品（含图片）

```bash
# 1. 登录获取Token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq -r '.token')

# 2. 上传菜品图片
IMAGE_URL=$(curl -s -X POST http://localhost:8080/api/v1/admin/upload/media \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@dish.jpg" | jq -r '.url')

# 3. 创建菜品
curl -X POST http://localhost:8080/api/v1/admin/dishes \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"红烧肉\",
    \"description\": \"经典家常菜\",
    \"category_id\": 1,
    \"price\": 35.00,
    \"image_url\": \"$IMAGE_URL\",
    \"is_seasonal\": false
  }"
```

#### 场景2：用户生成随机搭配并下单

```bash
# 1. 登录
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user","password":"user123"}' | jq -r '.token')

# 2. 生成随机搭配
curl -X POST http://localhost:8080/api/v1/combo/generate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"meat_count":1,"vegetable_count":2}'

# 3. 创建订单（使用返回的菜品ID）
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"items":[{"dish_id":1,"quantity":1},{"dish_id":5,"quantity":1}]}'
```

---

**📖 相关文档：**
- [README.md](../README.md) - 项目概述和快速开始
- [DEPLOYMENT.md](DEPLOYMENT.md) - 生产环境部署指南
- [API测试工具](http://localhost:8080/api-tester.html) - 可视化API测试

**🔗 有用链接：**
- [PostgreSQL文档](https://www.postgresql.org/docs/)
- [AWS S3文档](https://docs.aws.amazon.com/s3/)
- [JWT.io](https://jwt.io/) - JWT调试工具
- [Postman](https://www.postman.com/) - API测试工具
```