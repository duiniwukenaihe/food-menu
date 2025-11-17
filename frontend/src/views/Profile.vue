<template>
  <div class="profile">
    <div class="container">
      <div class="profile-header">
        <h1>个人中心</h1>
        <el-button @click="goBack" type="default">
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
      </div>

      <div class="profile-content">
        <div class="user-info-card">
          <h2>用户信息</h2>
          <div class="info-grid">
            <div class="info-item">
              <label>用户名：</label>
              <span>{{ userStore.user?.username }}</span>
            </div>
            <div class="info-item">
              <label>邮箱：</label>
              <span>{{ userStore.user?.email || '未设置' }}</span>
            </div>
            <div class="info-item">
              <label>角色：</label>
              <el-tag :type="userStore.isAdmin ? 'danger' : 'primary'">
                {{ userStore.isAdmin ? '管理员' : '普通用户' }}
              </el-tag>
            </div>
            <div class="info-item">
              <label>注册时间：</label>
              <span>{{ formatDate(userStore.user?.created_at) }}</span>
            </div>
          </div>
        </div>

        <div class="stats-card">
          <h2>统计信息</h2>
          <div class="stats-grid">
            <div class="stat-item">
              <div class="stat-number">{{ orderStore.orders.length }}</div>
              <div class="stat-label">总订单数</div>
            </div>
            <div class="stat-item">
              <div class="stat-number">{{ orderStore.favorites.length }}</div>
              <div class="stat-label">收藏菜品</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import { useUserStore, useOrderStore } from '@/stores'

const router = useRouter()
const userStore = useUserStore()
const orderStore = useOrderStore()

onMounted(async () => {
  await Promise.all([
    orderStore.fetchOrders({ limit: 10 }),
    orderStore.fetchFavorites({ limit: 10 })
  ])
})

function goBack() {
  router.go(-1)
}

function formatDate(dateString?: string) {
  if (!dateString) return '--'
  return new Date(dateString).toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.profile {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 20px 0;
}

.profile-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.profile-header h1 {
  margin: 0;
  color: #333;
}

.profile-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 30px;
}

.user-info-card,
.stats-card {
  background: white;
  padding: 30px;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.user-info-card h2,
.stats-card h2 {
  margin: 0 0 20px 0;
  color: #333;
  font-size: 1.3rem;
}

.info-grid {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.info-item {
  display: flex;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.info-item:last-child {
  border-bottom: none;
}

.info-item label {
  font-weight: 500;
  color: #666;
  width: 100px;
  flex-shrink: 0;
}

.info-item span {
  color: #333;
}

.stats-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.stat-item {
  text-align: center;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
}

.stat-number {
  font-size: 2rem;
  font-weight: bold;
  color: #409eff;
  margin-bottom: 8px;
}

.stat-label {
  color: #666;
  font-size: 0.9rem;
}

@media (max-width: 768px) {
  .profile-content {
    grid-template-columns: 1fr;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>