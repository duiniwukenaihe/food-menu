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
  "image_url": "https://example.com/image.jpg",
  "video_url": "https://example.com/video.mp4",
  "cooking_steps": "制作步骤",
  "is_seasonal": false
}
```

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
- `401` - 未授权
- `403` - 权限不足
- `404` - 资源不存在
- `500` - 服务器内部错误

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
```