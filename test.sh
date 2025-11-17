#!/bin/bash

# 食物点餐系统 - 功能测试脚本

set -e

API_BASE="http://localhost:8080/api/v1"

echo "🧪 食物点餐系统 - 功能测试"
echo "=================================="

# 检查是否只运行Go单元测试
if [ "$1" == "unit" ]; then
    echo ""
    echo "🧪 运行Go单元测试..."
    echo "======================================"
    cd backend
    export TEST_DATABASE_URL="postgres://postgres:password@localhost/food_ordering_test?sslmode=disable"
    go test -v ./handlers
    exit $?
fi

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试结果统计
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 测试函数
test_api() {
    local method=$1
    local endpoint=$2
    local data=$3
    local expected_status=$4
    local description=$5
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    echo -n "测试 $TOTAL_TESTS: $description ... "
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" "$API_BASE$endpoint")
    elif [ "$method" = "POST" ]; then
        response=$(curl -s -w "\n%{http_code}" -X POST -H "Content-Type: application/json" -d "$data" "$API_BASE$endpoint")
    elif [ "$method" = "PUT" ]; then
        response=$(curl -s -w "\n%{http_code}" -X PUT -H "Content-Type: application/json" -d "$data" "$API_BASE$endpoint")
    elif [ "$method" = "DELETE" ]; then
        response=$(curl -s -w "\n%{http_code}" -X DELETE "$API_BASE$endpoint")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" = "$expected_status" ]; then
        echo -e "${GREEN}✅ 通过${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}❌ 失败${NC}"
        echo "   期望状态码: $expected_status, 实际: $http_code"
        echo "   响应内容: $body"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
}

# 检查服务是否运行
check_server() {
    echo "🔍 检查服务状态..."
    
    if curl -s "$API_BASE/health" > /dev/null; then
        echo -e "${GREEN}✅ 后端服务运行正常${NC}"
    else
        echo -e "${RED}❌ 后端服务未运行，请先启动服务${NC}"
        echo "运行命令: ./start.sh"
        exit 1
    fi
}

# 开始测试
echo ""
echo "🚀 开始功能测试..."

# 检查服务状态
check_server

echo ""
echo "📋 基础API测试"
echo "=================="

# 1. 健康检查
test_api "GET" "/health" "" "200" "健康检查"

# 2. 获取分类列表
test_api "GET" "/categories" "" "200" "获取分类列表"

# 3. 获取菜品列表
test_api "GET" "/dishes" "" "200" "获取菜品列表"

# 4. 获取推荐配置
test_api "GET" "/recommendations" "" "200" "获取推荐配置"

# 5. 获取应季菜品
test_api "GET" "/seasonal-dishes" "" "200" "获取应季菜品"

echo ""
echo "🔐 认证相关测试"
echo "=================="

# 6. 用户登录 - 管理员
test_api "POST" "/login" '{"username":"admin","password":"admin123"}' "200" "管理员登录"

# 提取token
admin_token=$(curl -s -X POST -H "Content-Type: application/json" -d '{"username":"admin","password":"admin123"}' "$API_BASE/login" | jq -r '.token')

# 7. 用户登录 - 普通用户
test_api "POST" "/login" '{"username":"user","password":"user123"}' "200" "普通用户登录"

# 提取token
user_token=$(curl -s -X POST -H "Content-Type: application/json" -d '{"username":"user","password":"user123"}' "$API_BASE/login" | jq -r '.token')

# 8. 错误登录测试
test_api "POST" "/login" '{"username":"wrong","password":"wrong"}' "401" "错误登录"

echo ""
echo "👤 用户功能测试"
echo "=================="

# 9. 获取用户信息（管理员）
test_api "GET" "/profile" "" "200" "获取管理员信息"

# 10. 获取用户信息（普通用户）
test_api "GET" "/profile" "" "200" "获取普通用户信息"

echo ""
echo "🍽️ 菜品管理测试"
echo "=================="

# 11. 创建菜品（管理员）
new_dish='{"name":"测试菜品","description":"测试描述","category_id":1,"price":25.00,"image_url":"https://example.com/test.jpg","cooking_steps":"1. 准备食材\n2. 开始烹饪\n3. 完成装盘","is_seasonal":false}'
test_api "POST" "/admin/dishes" "$new_dish" "201" "创建菜品"

# 获取新创建的菜品ID
dish_id=$(curl -s -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $admin_token" -d "$new_dish" "$API_BASE/admin/dishes" | jq -r '.id')

# 12. 获取单个菜品
test_api "GET" "/dishes/$dish_id" "" "200" "获取单个菜品"

# 13. 更新菜品（管理员）
update_dish='{"name":"更新的测试菜品","price":30.00}'
test_api "PUT" "/admin/dishes/$dish_id" "$update_dish" "200" "更新菜品"

# 14. 获取用户列表（管理员）
test_api "GET" "/admin/users" "" "200" "获取用户列表"

# 15. 获取系统配置（管理员）
test_api "GET" "/admin/config" "" "200" "获取系统配置"

echo ""
echo "🛒 订单功能测试"
echo "=================="

# 16. 创建订单（普通用户）
create_order='{"items":[{"dish_id":1,"quantity":2}]}'
test_api "POST" "/orders" "$create_order" "201" "创建订单"

# 17. 获取订单列表
test_api "GET" "/orders" "" "200" "获取订单列表"

echo ""
echo "❌ 错误处理测试"
echo "=================="

# 18. 未授权访问管理接口
test_api "GET" "/admin/users" "" "401" "未授权访问管理接口"

# 19. 访问不存在的菜品
test_api "GET" "/dishes/99999" "" "404" "访问不存在的菜品"

# 20. 创建无效订单
invalid_order='{"items":[]}'
test_api "POST" "/orders" "$invalid_order" "400" "创建无效订单"

echo ""
echo "🧹 清理测试数据"
echo "=================="

# 删除测试菜品
if [ "$dish_id" != "null" ] && [ "$dish_id" != "" ]; then
    curl -s -X DELETE -H "Authorization: Bearer $admin_token" "$API_BASE/admin/dishes/$dish_id" > /dev/null
    echo "🗑️  已删除测试菜品 (ID: $dish_id)"
fi

echo ""
echo "📊 测试结果统计"
echo "=================="
echo -e "总测试数: $TOTAL_TESTS"
echo -e "${GREEN}通过: $PASSED_TESTS${NC}"
echo -e "${RED}失败: $FAILED_TESTS${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}🎉 所有测试通过！${NC}"
    exit 0
else
    echo -e "${RED}❌ 有 $FAILED_TESTS 个测试失败${NC}"
    exit 1
fi