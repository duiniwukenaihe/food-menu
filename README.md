# 🍽️ 食物点餐系统

一个基于Go后端和Vue3前端的现代化食物点餐网站，支持用户点餐、菜品管理、随机搭配等功能。

## ✨ 功能特性

### 🎨 用户界面
- **北极熊主题UI** - 可爱的北极熊登录界面和熊掌欢迎图案
- **响应式设计** - 完美适配桌面和移动设备
- **现代化界面** - 基于Element Plus的美观UI组件

### 🍜 核心功能
- **菜品浏览** - 网格化展示所有菜品，支持分类筛选和搜索
- **菜品详情** - 查看菜品图片、制作步骤、视频教程
- **应季推荐** - 根据季节推荐特色菜品
- **随机搭配** - 智能推荐荤素搭配，支持自定义配置
- **猜你喜欢** - 随机生成个性化菜品组合

### 👤 用户系统
- **用户登录** - 安全的JWT认证系统
- **个人中心** - 查看个人信息和统计数据
- **订单管理** - 创建、查看、管理个人订单
- **收藏功能** - 收藏喜欢的菜品

### 🛠️ 管理后台
- **菜品管理** - 添加、编辑、删除菜品
- **媒体上传** - 上传图片和视频到S3存储（新功能）
- **分类管理** - 管理菜品分类
- **用户管理** - 查看和管理系统用户
- **系统配置** - 配置S3存储、JWT、默认参数等

> **工作流程变更说明：**
> - ✅ **新方式：** 使用媒体上传接口上传文件到S3，获取URL后创建/更新菜品
> - ❌ **旧方式：** 手动输入外部图片URL（已不推荐）

### 🗄️ 技术特性
- **PostgreSQL数据库** - 高性能关系型数据库
- **S3对象存储** - 支持阿里云OSS、AWS S3、腾讯云COS、MinIO等
- **JWT认证** - 安全的用户认证机制
- **RESTful API** - 标准化的API接口设计

## 🚀 快速开始

### 环境要求

- **Go 1.21+**
- **Node.js 16+**
- **PostgreSQL 12+**

### 一键启动

```bash
# 克隆项目
git clone <repository-url>
cd food-ordering

# 运行启动脚本
./start.sh
```

启动脚本会自动：
- 检查系统依赖
- 初始化数据库
- 安装前后端依赖
- 启动开发服务器

### 手动启动

#### 1. 数据库设置

```bash
# 创建数据库
createdb food_ordering

# 初始化表结构
psql -d food_ordering -f database/schema.sql
```

#### 2. 后端启动

```bash
cd backend

# 安装依赖
go mod tidy

# 配置环境变量
cp .env.example .env
# 编辑 .env 文件配置数据库连接等

# 启动服务
go run main.go
```

#### 3. 前端启动

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

## 🌐 访问地址

启动成功后，可以通过以下地址访问：

- **前端应用**: http://localhost:3000
- **后端API**: http://localhost:8080
- **API文档**: [docs/API.md](docs/API.md)

## 👤 测试账号

系统提供以下测试账号：

| 角色 | 用户名 | 密码 |
|------|--------|------|
| 管理员 | admin | admin123 |
| 普通用户 | user | user123 |

## 📁 项目结构

```
food-ordering/
├── backend/                 # Go后端代码
│   ├── config/             # 配置管理
│   ├── database/           # 数据库连接
│   ├── handlers/           # HTTP处理器
│   ├── middleware/         # 中间件
│   ├── models/             # 数据模型
│   └── main.go            # 主程序入口
├── frontend/               # Vue3前端代码
│   ├── src/
│   │   ├── components/     # Vue组件
│   │   ├── views/         # 页面视图
│   │   ├── stores/        # Pinia状态管理
│   │   ├── types/         # TypeScript类型定义
│   │   ├── utils/         # 工具函数
│   │   └── main.ts        # 前端入口
│   └── package.json       # 前端依赖配置
├── database/              # 数据库脚本
│   └── schema.sql        # 数据库表结构
├── docs/                 # 项目文档
│   ├── API.md           # API接口文档
│   └── DEPLOYMENT.md    # 部署说明
├── start.sh              # 一键启动脚本
└── README.md            # 项目说明
```

## 🎯 核心功能演示

### 1. 用户登录
- 访问首页，点击右上角"用户登录"
- 输入测试账号：admin/admin123
- 查看可爱的北极熊登录界面

<!-- 截图占位: docs/screenshots/login.png -->

### 2. 浏览菜品
- 登录后显示熊掌欢迎图案
- 浏览所有菜品，查看图片和价格
- 点击菜品查看详细信息、制作步骤和视频

<!-- 截图占位: docs/screenshots/dishes-grid.png -->
<!-- 截图占位: docs/screenshots/dish-detail.png -->

### 3. **媒体上传（新功能）**
- 管理员登录后台
- 点击"上传图片/视频"
- 选择文件上传到S3存储
- 系统自动返回可用的URL地址
- 在创建菜品时使用上传的媒体URL

<!-- 截图占位: docs/screenshots/media-upload.png -->

### 4. 随机搭配
- 使用"猜你喜欢"功能
- 自定义荤素数量配置
- 一键生成个性化菜品组合
- 查看推荐的荤素搭配

<!-- 截图占位: docs/screenshots/combo-generator.png -->

### 5. 下单点餐
- 选择喜欢的菜品加入订单
- 确认订单信息并提交
- 在个人中心查看订单状态

<!-- 截图占位: docs/screenshots/order.png -->

### 6. 管理后台
- 使用管理员账号登录
- 访问管理后台进行菜品管理
- 使用媒体上传功能
- 配置S3存储和系统参数

<!-- 截图占位: docs/screenshots/admin-panel.png -->
<!-- 截图占位: docs/screenshots/admin-config.png -->

### 7. **API测试工具（新功能）**
- 访问 http://localhost:8080/api-tester.html
- 可视化测试所有API接口
- 支持文件上传测试
- 自动管理认证Token

<!-- 截图占位: docs/screenshots/api-tester.png -->

## 🔧 配置说明

### 环境变量配置

创建或编辑 `backend/.env` 文件，配置以下环境变量：

```bash
# 服务器配置
PORT=8080

# 数据库配置
DATABASE_URL=postgres://postgres:password@localhost/food_ordering?sslmode=disable
TEST_DATABASE_URL=postgres://postgres:password@localhost/food_ordering_test?sslmode=disable

# JWT认证配置
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRATION=24h  # Token过期时间：24h, 30m, 3600s等

# S3对象存储配置（必需，用于媒体文件上传）
S3_ENDPOINT=https://s3.us-west-2.amazonaws.com
S3_ACCESS_KEY=your_access_key_id
S3_SECRET_KEY=your_secret_access_key
S3_BUCKET=your_bucket_name
S3_REGION=us-west-2
S3_PATH_STYLE=false  # MinIO需要设为true
```

### 数据库配置

#### 主数据库

```bash
DATABASE_URL=postgres://username:password@localhost/food_ordering?sslmode=disable
```

#### 测试数据库（可选）

用于运行测试，避免影响主数据库：

```bash
TEST_DATABASE_URL=postgres://username:password@localhost/food_ordering_test?sslmode=disable
```

### S3对象存储配置

系统使用S3存储菜品图片和视频，支持多种S3兼容存储服务。

#### AWS S3

```bash
S3_ENDPOINT=https://s3.us-west-2.amazonaws.com
S3_ACCESS_KEY=AKIAIOSFODNN7EXAMPLE
S3_SECRET_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
S3_BUCKET=my-food-ordering-bucket
S3_REGION=us-west-2
S3_PATH_STYLE=false
```

#### 阿里云 OSS

```bash
S3_ENDPOINT=https://oss-cn-beijing.aliyuncs.com
S3_ACCESS_KEY=your_access_key
S3_SECRET_KEY=your_secret_key
S3_BUCKET=your_bucket_name
S3_REGION=oss-cn-beijing
S3_PATH_STYLE=false
```

#### 腾讯云 COS

```bash
S3_ENDPOINT=https://cos.ap-beijing.myqcloud.com
S3_ACCESS_KEY=your_secret_id
S3_SECRET_KEY=your_secret_key
S3_BUCKET=your_bucket_name
S3_REGION=ap-beijing
S3_PATH_STYLE=false
```

#### MinIO（本地开发）

```bash
S3_ENDPOINT=http://localhost:9000
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_BUCKET=food-ordering
S3_REGION=us-east-1
S3_PATH_STYLE=true  # ⚠️ MinIO必须设为true
```

**MinIO快速启动：**

```bash
# 使用Docker运行MinIO
docker run -d \
  -p 9000:9000 \
  -p 9001:9001 \
  --name minio \
  -e "MINIO_ROOT_USER=minioadmin" \
  -e "MINIO_ROOT_PASSWORD=minioadmin" \
  minio/minio server /data --console-address ":9001"

# 访问MinIO控制台创建bucket
# http://localhost:9001 (用户名：minioadmin，密码：minioadmin)
```

### JWT认证配置

```bash
# JWT密钥（生产环境请使用强密钥）
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Token过期时间
JWT_EXPIRATION=24h  # 24小时
# 其他示例：30m（30分钟）、168h（7天）、3600s（1小时）
```

**生成安全的JWT密钥：**

```bash
# Linux/Mac
openssl rand -base64 32

# 或使用在线工具
# https://www.random.org/strings/
```

## 🗄️ 数据库迁移

### 初始化数据库

首次启动时，运行以下命令初始化数据库表结构：

```bash
# 创建数据库
createdb food_ordering

# 或使用psql
psql -U postgres -c "CREATE DATABASE food_ordering;"

# 运行初始化脚本
psql -d food_ordering -f database/schema.sql
```

### 运行迁移

系统会在启动时自动运行数据库迁移。如需手动运行：

```bash
cd backend
go run main.go --migrate
```

### 创建测试数据库

```bash
# 创建测试数据库
createdb food_ordering_test

# 初始化测试数据库
psql -d food_ordering_test -f database/schema.sql
```

### 迁移文件位置

- 主schema：`database/schema.sql`
- 数据库备份：`database/backup/`（如有）

## 🧪 API测试工具

系统内置Web API测试工具，方便开发和调试。

### 访问测试页面

启动后端服务后，访问：

```
http://localhost:8080/api-tester.html
```

### 功能特性

- ✅ 可视化测试所有API接口
- ✅ 自动管理JWT Token
- ✅ 支持文件上传测试（图片/视频）
- ✅ 显示请求/响应详情
- ✅ 保存常用请求
- ✅ 错误信息高亮显示

### 使用示例

1. 打开API测试页面
2. 点击"登录"选项，输入测试账号（admin/admin123）
3. 获取Token后，选择要测试的接口
4. 填写参数，点击"发送请求"
5. 查看响应结果

## 🐳 Docker部署

使用Docker Compose一键部署：

```bash
docker-compose up -d
```

## 📚 API文档

详细的API接口文档请参考：[docs/API.md](docs/API.md)

主要接口包括：
- 🔐 认证接口（登录、用户信息）
- 🍽️ 菜品管理（CRUD操作）
- 📤 **媒体上传**（图片/视频上传到S3）
- 🎲 **随机搭配**（智能生成菜品组合）
- ⚙️ **公共配置**（获取系统配置）
- 👤 用户管理（管理员）
- 📦 订单管理
- ⭐ 收藏功能

## 🚀 部署指南

生产环境部署请参考：[docs/DEPLOYMENT.md](docs/DEPLOYMENT.md)

## 🛠️ 开发指南

### 后端开发

- 使用Gin框架构建RESTful API
- GORM作为ORM工具
- JWT进行用户认证
- 支持CORS跨域请求

### 前端开发

- Vue 3 + TypeScript
- Pinia状态管理
- Element Plus UI组件库
- Vite构建工具

### 数据库设计

- 用户表(users) - 存储用户信息
- 菜品表(dishes) - 存储菜品数据
- 分类表(categories) - 菜品分类
- 订单表(orders) - 用户订单
- 收藏表(user_favorites) - 用户收藏

## 🔍 故障排除

### S3存储相关问题

#### 问题：上传文件失败 - "S3 storage is not configured"

**原因：** S3环境变量未配置或配置不正确

**解决方案：**
```bash
# 1. 检查 backend/.env 文件是否存在
ls -la backend/.env

# 2. 确认S3配置项都已填写
cat backend/.env | grep S3

# 3. 确保以下配置都有值（不为空）
S3_ENDPOINT=https://your-endpoint
S3_ACCESS_KEY=your-key
S3_SECRET_KEY=your-secret
S3_BUCKET=your-bucket
S3_REGION=your-region
```

#### 问题：MinIO连接失败 - "Connection refused"

**原因：** MinIO服务未启动或端点配置错误

**解决方案：**
```bash
# 1. 检查MinIO容器是否运行
docker ps | grep minio

# 2. 如未运行，启动MinIO
docker run -d -p 9000:9000 -p 9001:9001 --name minio \
  -e "MINIO_ROOT_USER=minioadmin" \
  -e "MINIO_ROOT_PASSWORD=minioadmin" \
  minio/minio server /data --console-address ":9001"

# 3. 确认端点配置
S3_ENDPOINT=http://localhost:9000  # 注意是http不是https
S3_PATH_STYLE=true                 # MinIO必须设为true
```

#### 问题：阿里云OSS连接失败 - "SignatureDoesNotMatch"

**原因：** 访问密钥配置错误或Bucket名称错误

**解决方案：**
```bash
# 1. 验证访问密钥
# 登录阿里云控制台，检查AccessKey ID和AccessKey Secret

# 2. 确认Bucket名称正确
# 注意：不要包含协议和域名，只填写bucket名称

# 3. 检查区域设置
S3_REGION=oss-cn-beijing  # 确保与Bucket所在区域一致
```

#### 问题：文件上传后无法访问 - 403 Forbidden

**原因：** S3 Bucket权限设置问题

**解决方案：**
```bash
# AWS S3: 设置Bucket Policy允许公开读取
# 阿里云OSS: 设置Bucket ACL为"公共读"
# MinIO: 设置Bucket Policy

# MinIO示例（通过mc命令行工具）
mc anonymous set download myminio/food-ordering
```

### 数据库迁移问题

#### 问题：数据库连接失败 - "connection refused"

**原因：** PostgreSQL服务未启动或连接配置错误

**解决方案：**
```bash
# 1. 检查PostgreSQL服务状态
sudo systemctl status postgresql
# 或
pg_isready

# 2. 启动PostgreSQL
sudo systemctl start postgresql

# 3. 验证连接配置
psql -h localhost -U postgres -d food_ordering

# 4. 检查 .env 中的 DATABASE_URL 格式
DATABASE_URL=postgres://用户名:密码@主机:端口/数据库名?sslmode=disable
```

#### 问题：迁移脚本执行失败 - "relation already exists"

**原因：** 数据库表已存在，重复执行了初始化脚本

**解决方案：**
```bash
# 方法1：删除并重新创建数据库
dropdb food_ordering
createdb food_ordering
psql -d food_ordering -f database/schema.sql

# 方法2：跳过已存在的表（schema.sql使用了IF NOT EXISTS）
# 直接忽略此错误，不影响使用
```

#### 问题：初始数据插入失败 - "duplicate key value"

**原因：** 默认数据已存在

**解决方案：**
```bash
# 这通常不是错误，可以安全忽略
# 或者清空特定表后重新插入

psql -d food_ordering -c "TRUNCATE categories CASCADE;"
psql -d food_ordering -c "TRUNCATE system_config CASCADE;"
# 然后重新运行schema.sql中的INSERT语句
```

### JWT认证问题

#### 问题：Token失效过快

**原因：** JWT_EXPIRATION设置太短

**解决方案：**
```bash
# 在 backend/.env 中调整过期时间
JWT_EXPIRATION=24h  # 设置为24小时
# 或
JWT_EXPIRATION=168h # 设置为7天

# 重启后端服务使配置生效
```

#### 问题：Token验证失败 - "Invalid token"

**原因：** JWT_SECRET不一致或Token格式错误

**解决方案：**
```bash
# 1. 确保JWT_SECRET在所有环境中一致
# 2. 清除旧Token，重新登录获取新Token
# 3. 检查Token格式：Authorization: Bearer <token>
```

### 前端连接问题

#### 问题：API请求失败 - CORS错误

**原因：** 跨域配置问题

**解决方案：**
```bash
# 检查 backend/main.go 中的CORS配置
# 确保前端地址在 AllowOrigins 列表中
AllowOrigins: []string{"http://localhost:3000"}
```

#### 问题：图片/视频无法显示

**原因：** S3 URL无法访问或Bucket权限问题

**解决方案：**
```bash
# 1. 检查S3 URL是否可以在浏览器中直接访问
# 2. 确认Bucket设置为公共读取
# 3. 检查CORS配置（S3 Bucket CORS规则）
```

### 常见启动错误

#### 问题：端口已被占用

```bash
# 查找占用端口的进程
lsof -i :8080  # 后端端口
lsof -i :3000  # 前端端口

# 杀死占用端口的进程
kill -9 <PID>

# 或修改端口配置
# backend/.env: PORT=8081
# frontend/vite.config.ts: server.port = 3001
```

#### 问题：依赖安装失败

```bash
# Go依赖问题
cd backend
go clean -modcache
go mod tidy
go mod download

# Node依赖问题
cd frontend
rm -rf node_modules package-lock.json
npm install
```

### 获取更多帮助

如果以上方法无法解决问题：

1. 查看详细错误日志
2. 访问 [API文档](docs/API.md) 了解接口详情
3. 查看 [部署指南](docs/DEPLOYMENT.md)
4. 提交 Issue 到项目仓库

## 🤝 贡献指南

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 🎉 致谢

感谢以下开源项目：

- [Gin](https://github.com/gin-gonic/gin) - Go Web框架
- [Vue.js](https://vuejs.org/) - 前端框架
- [Element Plus](https://element-plus.org/) - Vue组件库
- [PostgreSQL](https://www.postgresql.org/) - 数据库
- [GORM](https://gorm.io/) - Go ORM库

## 🏗️ 基础设施

项目包含 Terraform 基础设施代码，用于在 Proxmox VE 上部署 Kubernetes 集群：

- **位置**: `infrastructure/terraform/`
- **用途**: 自动化部署和管理 Kubernetes 基础设施
- **文档**: 查看 [infrastructure/README.md](infrastructure/README.md)

快速开始：
```bash
cd infrastructure/terraform
cp terraform.tfvars.example terraform.tfvars
# 编辑配置文件
terraform init
terraform plan
terraform apply
```

## 📞 联系方式

如有问题或建议，请通过以下方式联系：

- 提交 Issue
- 发送邮件
- 创建 Pull Request

---

**🍽️ 享受美食，享受生活！**