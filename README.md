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
- **分类管理** - 管理菜品分类
- **用户管理** - 查看和管理系统用户
- **系统配置** - 配置S3存储、默认参数等

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

### 2. 浏览菜品
- 登录后显示熊掌欢迎图案
- 浏览所有菜品，查看图片和价格
- 点击菜品查看详细信息、制作步骤和视频

### 3. 随机搭配
- 使用"猜你喜欢"功能
- 自定义荤素数量配置
- 一键生成个性化菜品组合

### 4. 下单点餐
- 选择喜欢的菜品加入订单
- 确认订单信息并提交
- 在个人中心查看订单状态

### 5. 管理后台
- 使用管理员账号登录
- 访问管理后台进行菜品管理
- 配置S3存储和系统参数

## 🔧 配置说明

### 数据库配置

编辑 `backend/.env` 文件：

```bash
DATABASE_URL=postgres://username:password@localhost/food_ordering?sslmode=disable
```

### S3对象存储配置

支持多种S3兼容存储服务：

```bash
# 阿里云OSS
S3_ENDPOINT=https://oss-cn-beijing.aliyuncs.com
S3_ACCESS_KEY=your_access_key
S3_SECRET_KEY=your_secret_key
S3_BUCKET=your_bucket_name
S3_REGION=oss-cn-beijing

# AWS S3
S3_ENDPOINT=https://s3.us-west-2.amazonaws.com
S3_REGION=us-west-2

# 腾讯云COS
S3_ENDPOINT=https://cos.ap-beijing.myqcloud.com
S3_REGION=ap-beijing

# MinIO
S3_ENDPOINT=http://localhost:9000
S3_REGION=us-east-1
```

## 🐳 Docker部署

使用Docker Compose一键部署：

```bash
docker-compose up -d
```

## 📚 API文档

详细的API接口文档请参考：[docs/API.md](docs/API.md)

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

## 📞 联系方式

如有问题或建议，请通过以下方式联系：

- 提交 Issue
- 发送邮件
- 创建 Pull Request

---

**🍽️ 享受美食，享受生活！**