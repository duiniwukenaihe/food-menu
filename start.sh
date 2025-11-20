#!/bin/bash

# 食物点餐系统启动脚本

set -e

echo "🍽️  食物点餐系统 - 启动脚本"
echo "=================================="

# 检查依赖
check_dependencies() {
    echo "📋 检查系统依赖..."
    
    # 检查Go
    if ! command -v go &> /dev/null; then
        echo "❌ Go未安装，请先安装Go 1.21+"
        exit 1
    fi
    echo "✅ Go版本: $(go version)"
    
    # 检查Node.js
    if ! command -v node &> /dev/null; then
        echo "❌ Node.js未安装，请先安装Node.js 16+"
        exit 1
    fi
    echo "✅ Node.js版本: $(node --version)"
    
    # 检查npm
    if ! command -v npm &> /dev/null; then
        echo "❌ npm未安装"
        exit 1
    fi
    echo "✅ npm版本: $(npm --version)"
    
    # 检查PostgreSQL
    if ! command -v psql &> /dev/null; then
        echo "⚠️  PostgreSQL客户端未找到，请确保PostgreSQL已安装"
    else
        echo "✅ PostgreSQL已安装"
    fi
}

# 设置数据库
setup_database() {
    echo ""
    echo "🗄️  设置数据库..."
    
    # 检查数据库连接
    if PGPASSWORD=password psql -h localhost -U postgres -d food_ordering -c "SELECT 1;" &> /dev/null; then
        echo "✅ 数据库连接成功"
    else
        echo "❌ 数据库连接失败"
        echo "请确保："
        echo "1. PostgreSQL服务已启动"
        echo "2. 数据库 'food_ordering' 已创建"
        echo "3. 用户 'postgres' 存在且密码为 'password'"
        echo ""
        echo "快速修复命令："
        echo "  # 检查PostgreSQL状态"
        echo "  sudo systemctl status postgresql"
        echo ""
        echo "  # 启动PostgreSQL"
        echo "  sudo systemctl start postgresql"
        echo ""
        echo "  # 创建数据库"
        echo "  createdb food_ordering"
        echo "  # 或"
        echo "  sudo -u postgres createdb food_ordering"
        echo ""
        read -p "是否继续启动应用？(y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
    
    # 初始化数据库表
    if [ -f "database/schema.sql" ]; then
        echo "📝 运行数据库迁移..."
        # 使用2>&1重定向错误输出，因为IF NOT EXISTS会产生通知
        PGPASSWORD=password psql -h localhost -U postgres -d food_ordering -f database/schema.sql 2>&1 | grep -v "NOTICE"
        if [ $? -eq 0 ] || [ $? -eq 1 ]; then
            echo "✅ 数据库迁移完成"
        else
            echo "⚠️  数据库迁移可能有警告（这通常是正常的）"
        fi
    else
        echo "⚠️  数据库脚本文件不存在: database/schema.sql"
    fi
}

# 启动后端
start_backend() {
    echo ""
    echo "🚀 启动后端服务..."
    
    cd backend
    
    # 检查go.mod
    if [ ! -f "go.mod" ]; then
        echo "📦 初始化Go模块..."
        go mod init food-ordering
    fi
    
    # 安装依赖
    echo "📦 安装Go依赖..."
    go mod tidy
    
    # 创建.env文件
    if [ ! -f ".env" ]; then
        echo "⚙️  创建环境配置文件..."
        cp .env.example .env
        echo "✅ 已创建 .env 文件，请根据需要修改配置"
    fi
    
    # 启动后端
    echo "🌟 启动Go服务器..."
    go run main.go &
    BACKEND_PID=$!
    echo "✅ 后端服务已启动 (PID: $BACKEND_PID)"
    
    cd ..
}

# 启动前端
start_frontend() {
    echo ""
    echo "🎨 启动前端服务..."
    
    cd frontend
    
    # 安装依赖
    if [ ! -d "node_modules" ]; then
        echo "📦 安装npm依赖..."
        npm install
    fi
    
    # 启动前端
    echo "🌟 启动Vue开发服务器..."
    npm run dev &
    FRONTEND_PID=$!
    echo "✅ 前端服务已启动 (PID: $FRONTEND_PID)"
    
    cd ..
}

# 检查S3配置
check_s3_config() {
    echo ""
    echo "☁️  检查S3存储配置..."
    
    if [ -f "backend/.env" ]; then
        S3_CONFIGURED=false
        if grep -q "^S3_ENDPOINT=.\+" backend/.env && \
           grep -q "^S3_BUCKET=.\+" backend/.env; then
            S3_CONFIGURED=true
            echo "✅ S3配置已设置"
        else
            echo "⚠️  S3配置未完成"
            echo ""
            echo "媒体文件上传功能需要配置S3存储。"
            echo "请编辑 backend/.env 文件，配置以下变量："
            echo "  S3_ENDPOINT=https://your-s3-endpoint"
            echo "  S3_ACCESS_KEY=your-access-key"
            echo "  S3_SECRET_KEY=your-secret-key"
            echo "  S3_BUCKET=your-bucket-name"
            echo "  S3_REGION=your-region"
            echo "  S3_PATH_STYLE=false  # MinIO设为true"
            echo ""
            echo "支持的存储服务："
            echo "  - AWS S3"
            echo "  - 阿里云 OSS"
            echo "  - 腾讯云 COS"
            echo "  - MinIO (本地开发)"
            echo ""
            echo "详细配置说明请查看 README.md"
        fi
    fi
}

# 显示服务信息
show_info() {
    echo ""
    echo "🎉 服务启动成功！"
    echo "=================================="
    echo "📍 前端地址: http://localhost:3000"
    echo "📍 后端地址: http://localhost:8080"
    echo "📍 API文档: docs/API.md"
    echo "📍 API测试工具: http://localhost:8080/api-tester.html"
    echo ""
    echo "👤 测试账号:"
    echo "   管理员: admin / admin123"
    echo "   普通用户: user / user123"
    echo ""
    echo "📖 功能说明:"
    echo "   - 菜品管理（需登录）"
    echo "   - 媒体上传（需配置S3）"
    echo "   - 随机搭配推荐"
    echo "   - 订单管理"
    echo ""
    echo "🔧 配置文件:"
    echo "   - 后端配置: backend/.env"
    echo "   - 数据库迁移: database/schema.sql"
    echo ""
    echo "📚 文档:"
    echo "   - README.md - 快速开始和配置说明"
    echo "   - docs/API.md - API接口文档"
    echo "   - docs/DEPLOYMENT.md - 部署指南"
    echo ""
    echo "🛑 按 Ctrl+C 停止所有服务"
    echo "=================================="
}

# 清理函数
cleanup() {
    echo ""
    echo "🛑 正在停止服务..."
    
    if [ ! -z "$BACKEND_PID" ]; then
        kill $BACKEND_PID 2>/dev/null || true
        echo "✅ 后端服务已停止"
    fi
    
    if [ ! -z "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null || true
        echo "✅ 前端服务已停止"
    fi
    
    echo "👋 再见！"
    exit 0
}

# 设置信号处理
trap cleanup SIGINT SIGTERM

# 主函数
main() {
    # 解析命令行参数
    SKIP_DB=false
    SKIP_BACKEND=false
    SKIP_FRONTEND=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --skip-db)
                SKIP_DB=true
                shift
                ;;
            --skip-backend)
                SKIP_BACKEND=true
                shift
                ;;
            --skip-frontend)
                SKIP_FRONTEND=true
                shift
                ;;
            --help|-h)
                echo "用法: $0 [选项]"
                echo ""
                echo "选项:"
                echo "  --skip-db      跳过数据库设置"
                echo "  --skip-backend 跳过后端启动"
                echo "  --skip-frontend 跳过前端启动"
                echo "  --help, -h     显示帮助信息"
                exit 0
                ;;
            *)
                echo "未知选项: $1"
                echo "使用 --help 查看帮助信息"
                exit 1
                ;;
        esac
    done
    
    # 执行启动流程
    check_dependencies
    
    if [ "$SKIP_DB" = false ]; then
        setup_database
    fi
    
    # 检查S3配置（不阻塞启动）
    check_s3_config
    
    if [ "$SKIP_BACKEND" = false ]; then
        start_backend
    fi
    
    if [ "$SKIP_FRONTEND" = false ]; then
        start_frontend
    fi
    
    show_info
    
    # 等待信号
    wait
}

# 运行主函数
main "$@"