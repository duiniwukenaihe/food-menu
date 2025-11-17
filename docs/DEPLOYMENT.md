# 部署说明

## 环境要求

### 后端环境
- Go 1.21+
- PostgreSQL 12+
- Git

### 前端环境
- Node.js 16+
- npm 或 yarn

## 数据库设置

### 1. 安装PostgreSQL

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```

**CentOS/RHEL:**
```bash
sudo yum install postgresql-server postgresql-contrib
sudo postgresql-setup initdb
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

**macOS (使用Homebrew):**
```bash
brew install postgresql
brew services start postgresql
```

### 2. 创建数据库和用户

```bash
# 切换到postgres用户
sudo -u postgres psql

# 在PostgreSQL命令行中执行
CREATE DATABASE food_ordering;
CREATE USER food_ordering_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE food_ordering TO food_ordering_user;
\q
```

### 3. 初始化数据库表

```bash
cd database
psql -U food_ordering_user -d food_ordering -f schema.sql
```

## 后端部署

### 1. 进入后端目录

```bash
cd backend
```

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置环境变量

复制示例配置文件并根据实际情况修改：

```bash
cp .env.example .env
```

编辑 `.env` 文件：

```bash
# 服务器配置
PORT=8080

# 数据库配置
DATABASE_URL=postgres://food_ordering_user:your_password@localhost/food_ordering?sslmode=disable

# JWT密钥（请使用强密码）
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# S3配置（可选）
S3_ENDPOINT=https://oss-cn-beijing.aliyuncs.com
S3_ACCESS_KEY=your_access_key
S3_SECRET_KEY=your_secret_key
S3_BUCKET=your_bucket_name
S3_REGION=oss-cn-beijing
```

### 4. 运行后端服务

**开发模式:**
```bash
go run main.go
```

**生产模式:**
```bash
# 编译
go build -o food-ordering-server main.go

# 运行
./food-ordering-server
```

### 5. 使用systemd管理服务（生产环境推荐）

创建服务文件：

```bash
sudo nano /etc/systemd/system/food-ordering.service
```

内容如下：

```ini
[Unit]
Description=Food Ordering Service
After=network.target

[Service]
Type=simple
User=your_username
WorkingDirectory=/path/to/food-ordering/backend
ExecStart=/path/to/food-ordering/backend/food-ordering-server
Restart=always
RestartSec=3
Environment=DATABASE_URL=postgres://food_ordering_user:your_password@localhost/food_ordering?sslmode=disable
Environment=JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

[Install]
WantedBy=multi-user.target
```

启动和启用服务：

```bash
sudo systemctl daemon-reload
sudo systemctl start food-ordering
sudo systemctl enable food-ordering
```

## 前端部署

### 1. 进入前端目录

```bash
cd frontend
```

### 2. 安装依赖

```bash
npm install
```

### 3. 配置环境

如果需要修改API地址，编辑 `vite.config.ts` 中的proxy配置：

```typescript
export default defineConfig({
  // ...其他配置
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080', // 修改为实际的API地址
        changeOrigin: true,
      },
    },
  },
})
```

### 4. 开发模式运行

```bash
npm run dev
```

### 5. 生产环境构建

```bash
npm run build
```

构建完成后，`dist` 目录包含了所有静态文件。

### 6. 使用Nginx部署前端

**安装Nginx:**
```bash
# Ubuntu/Debian
sudo apt install nginx

# CentOS/RHEL
sudo yum install nginx
```

**配置Nginx:**

创建配置文件：

```bash
sudo nano /etc/nginx/sites-available/food-ordering
```

内容如下：

```nginx
server {
    listen 80;
    server_name your-domain.com;  # 替换为你的域名

    # 前端静态文件
    location / {
        root /path/to/food-ordering/frontend/dist;
        try_files $uri $uri/ /index.html;
    }

    # API代理到后端
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

启用配置：

```bash
sudo ln -s /etc/nginx/sites-available/food-ordering /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 使用Docker部署

### 1. 创建Docker Compose文件

在项目根目录创建 `docker-compose.yml`：

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:13
    environment:
      POSTGRES_DB: food_ordering
      POSTGRES_USER: food_ordering_user
      POSTGRES_PASSWORD: your_password
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./database/schema.sql:/docker-entrypoint-initdb.d/schema.sql
    ports:
      - "5432:5432"

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    environment:
      DATABASE_URL: postgres://food_ordering_user:your_password@postgres:5432/food_ordering?sslmode=disable
      JWT_SECRET: your-super-secret-jwt-key-change-this-in-production
    ports:
      - "8080:8080"
    depends_on:
      - postgres

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    ports:
      - "80:80"
    depends_on:
      - backend

volumes:
  postgres_data:
```

### 2. 创建后端Dockerfile

在 `backend` 目录创建 `Dockerfile`：

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env.example .env

EXPOSE 8080
CMD ["./main"]
```

### 3. 创建前端Dockerfile

在 `frontend` 目录创建 `Dockerfile`：

```dockerfile
FROM node:16-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### 4. 创建Nginx配置

在 `frontend` 目录创建 `nginx.conf`：

```nginx
events {
    worker_connections 1024;
}

http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    server {
        listen 80;
        server_name localhost;

        location / {
            root /usr/share/nginx/html;
            index index.html index.htm;
            try_files $uri $uri/ /index.html;
        }

        location /api {
            proxy_pass http://backend:8080;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
    }
}
```

### 5. 启动服务

```bash
docker-compose up -d
```

## SSL证书配置（HTTPS）

### 使用Let's Encrypt

**安装Certbot:**
```bash
sudo apt install certbot python3-certbot-nginx
```

**获取证书:**
```bash
sudo certbot --nginx -d your-domain.com
```

**自动续期:**
```bash
sudo crontab -e
# 添加以下行
0 12 * * * /usr/bin/certbot renew --quiet
```

## 监控和日志

### 1. 查看服务状态

```bash
# systemd服务
sudo systemctl status food-ordering

# Docker容器
docker-compose ps
```

### 2. 查看日志

```bash
# systemd服务
sudo journalctl -u food-ordering -f

# Docker容器
docker-compose logs -f
```

### 3. 数据库备份

创建备份脚本：

```bash
#!/bin/bash
BACKUP_DIR="/path/to/backups"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME="food_ordering"
DB_USER="food_ordering_user"

mkdir -p $BACKUP_DIR

pg_dump -U $DB_USER -h localhost $DB_NAME > $BACKUP_DIR/backup_$DATE.sql

# 删除7天前的备份
find $BACKUP_DIR -name "backup_*.sql" -mtime +7 -delete
```

设置定时备份：

```bash
sudo crontab -e
# 添加每天凌晨2点备份
0 2 * * * /path/to/backup_script.sh
```

## 性能优化

### 1. 数据库优化

- 为常用查询字段添加索引
- 定期执行VACUUM和ANALYZE
- 监控慢查询日志

### 2. 前端优化

- 启用Gzip压缩
- 设置适当的缓存策略
- 使用CDN加速静态资源

### 3. 后端优化

- 使用连接池
- 实现缓存机制
- 监控内存和CPU使用情况

## 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查数据库服务状态
   - 验证连接字符串
   - 确认防火墙设置

2. **前端无法访问API**
   - 检查CORS配置
   - 验证代理设置
   - 确认后端服务运行状态

3. **JWT认证失败**
   - 检查JWT密钥配置
   - 验证token格式
   - 确认token未过期

### 日志分析

关键日志位置：
- 后端应用日志：`/var/log/food-ordering/`
- Nginx访问日志：`/var/log/nginx/access.log`
- Nginx错误日志：`/var/log/nginx/error.log`
- PostgreSQL日志：`/var/log/postgresql/`

## 安全建议

1. 定期更新依赖包
2. 使用强密码和密钥
3. 启用HTTPS
4. 配置防火墙规则
5. 定期备份数据
6. 监控异常访问
7. 限制API访问频率
8. 验证所有用户输入