<template>
  <div class="admin-layout">
    <div class="admin-sidebar">
      <div class="sidebar-header">
        <h2>管理后台</h2>
      </div>
      <el-menu
        :default-active="$route.path"
        router
        class="sidebar-menu"
      >
        <el-menu-item index="/admin/dishes">
          <el-icon><Food /></el-icon>
          <span>菜品管理</span>
        </el-menu-item>
        <el-menu-item index="/admin/categories">
          <el-icon><Menu /></el-icon>
          <span>分类管理</span>
        </el-menu-item>
        <el-menu-item index="/admin/users">
          <el-icon><User /></el-icon>
          <span>用户管理</span>
        </el-menu-item>
        <el-menu-item index="/admin/config">
          <el-icon><Setting /></el-icon>
          <span>系统配置</span>
        </el-menu-item>
      </el-menu>
    </div>

    <div class="admin-main">
      <div class="admin-header">
        <div class="header-title">
          <h1>{{ getPageTitle() }}</h1>
        </div>
        <div class="header-actions">
          <el-button @click="goToHome" type="default">
            <el-icon><House /></el-icon>
            返回首页
          </el-button>
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-icon><User /></el-icon>
              {{ userStore.user?.username }}
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>

      <div class="admin-content">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { 
  Food, 
  Menu, 
  User, 
  Setting, 
  House, 
  ArrowDown 
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const getPageTitle = () => {
  const titleMap: Record<string, string> = {
    '/admin/dishes': '菜品管理',
    '/admin/categories': '分类管理',
    '/admin/users': '用户管理',
    '/admin/config': '系统配置'
  }
  return titleMap[route.path] || '管理后台'
}

function goToHome() {
  router.push('/')
}

function handleCommand(command: string) {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'logout':
      userStore.logout()
      ElMessage.success('退出登录成功')
      router.push('/login')
      break
  }
}
</script>

<style scoped>
.admin-layout {
  min-height: 100vh;
  display: flex;
  background: #f5f5f5;
}

.admin-sidebar {
  width: 250px;
  background: #001529;
  box-shadow: 2px 0 6px rgba(0, 21, 41, 0.35);
}

.sidebar-header {
  padding: 20px;
  text-align: center;
  border-bottom: 1px solid #1f1f1f;
}

.sidebar-header h2 {
  color: white;
  margin: 0;
  font-size: 1.2rem;
}

.sidebar-menu {
  border: none;
  background: #001529;
}

:deep(.el-menu-item) {
  color: rgba(255, 255, 255, 0.65);
  border-bottom: 1px solid #1f1f1f;
}

:deep(.el-menu-item:hover) {
  background-color: #1890ff !important;
  color: white;
}

:deep(.el-menu-item.is-active) {
  background-color: #1890ff !important;
  color: white;
}

.admin-main {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.admin-header {
  background: white;
  padding: 0 24px;
  height: 64px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
}

.header-title h1 {
  margin: 0;
  color: #262626;
  font-size: 1.5rem;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-info {
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 4px;
  transition: background-color 0.3s;
  color: #595959;
}

.user-info:hover {
  background-color: #f5f5f5;
}

.admin-content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
}
</style>