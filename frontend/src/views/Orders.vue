<template>
  <div class="orders">
    <div class="container">
      <div class="orders-header">
        <h1>我的订单</h1>
        <el-button @click="goBack" type="default">
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
      </div>

      <div class="orders-content">
        <div v-if="orderStore.isLoading" class="loading-section">
          <el-skeleton :rows="5" animated />
        </div>

        <div v-else-if="orderStore.orders.length === 0" class="empty-section">
          <el-empty description="暂无订单">
            <el-button type="primary" @click="goToHome">去点餐</el-button>
          </el-empty>
        </div>

        <div v-else class="orders-list">
          <div
            v-for="order in orderStore.orders"
            :key="order.id"
            class="order-card"
          >
            <div class="order-header">
              <div class="order-info">
                <span class="order-id">订单号：{{ order.id }}</span>
                <span class="order-time">{{ formatDate(order.created_at) }}</span>
              </div>
              <div class="order-status">
                <el-tag :type="getStatusType(order.status)">
                  {{ getStatusText(order.status) }}
                </el-tag>
              </div>
            </div>

            <div class="order-items">
              <div
                v-for="item in order.items"
                :key="item.id"
                class="order-item"
              >
                <img 
                  :src="item.dish?.image_url || '/placeholder-food.jpg'" 
                  :alt="item.dish?.name" 
                  class="item-image"
                />
                <div class="item-details">
                  <h4>{{ item.dish?.name }}</h4>
                  <p>数量：{{ item.quantity }}</p>
                </div>
                <div class="item-price">
                  ¥{{ (item.price * item.quantity).toFixed(2) }}
                </div>
              </div>
            </div>

            <div class="order-footer">
              <div class="total-amount">
                总计：<span>¥{{ order.total_amount.toFixed(2) }}</span>
              </div>
              <div class="order-actions">
                <el-button 
                  v-if="order.status === 'pending'"
                  type="danger" 
                  size="small"
                  @click="cancelOrder(order.id)"
                >
                  取消订单
                </el-button>
                <el-button 
                  type="primary" 
                  size="small"
                  @click="reorder(order)"
                >
                  再来一单
                </el-button>
              </div>
            </div>
          </div>
        </div>

        <!-- 分页 -->
        <div v-if="orderStore.total > 0" class="pagination-wrapper">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :total="orderStore.total"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </div>

    <!-- 订单对话框 -->
    <OrderDialog 
      v-model="showOrderDialog" 
      :items="orderItems" 
      @success="handleOrderSuccess" 
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import { useOrderStore } from '@/stores'
import { api } from '@/utils/api'
import OrderDialog from '@/components/OrderDialog.vue'
import type { Order, Dish } from '@/types'

const router = useRouter()
const orderStore = useOrderStore()

const currentPage = ref(1)
const pageSize = ref(10)
const showOrderDialog = ref(false)
const orderItems = ref<{ dish: Dish; quantity: number }[]>([])

onMounted(() => {
  loadOrders()
})

async function loadOrders() {
  await orderStore.fetchOrders({
    page: currentPage.value,
    limit: pageSize.value
  })
}

function goBack() {
  router.go(-1)
}

function goToHome() {
  router.push('/')
}

function formatDate(dateString: string) {
  return new Date(dateString).toLocaleString('zh-CN')
}

function getStatusType(status: string) {
  const statusMap: Record<string, string> = {
    pending: 'warning',
    confirmed: 'primary',
    preparing: 'info',
    ready: 'success',
    completed: 'success',
    cancelled: 'danger'
  }
  return statusMap[status] || 'info'
}

function getStatusText(status: string) {
  const statusMap: Record<string, string> = {
    pending: '待确认',
    confirmed: '已确认',
    preparing: '准备中',
    ready: '已完成',
    completed: '已完成',
    cancelled: '已取消'
  }
  return statusMap[status] || status
}

async function cancelOrder(orderId: number) {
  try {
    await ElMessageBox.confirm(
      '确定要取消这个订单吗？',
      '取消订单',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    // 这里应该调用取消订单的API
    // await api.cancelOrder(orderId)
    
    ElMessage.success('订单已取消')
    await loadOrders()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Cancel order failed:', error)
      ElMessage.error('取消订单失败')
    }
  }
}

function reorder(order: Order) {
  if (!order.items) return
  
  const items = order.items
    .filter(item => item.dish)
    .map(item => ({
      dish: item.dish!,
      quantity: item.quantity
    }))
  
  orderItems.value = items
  showOrderDialog.value = true
}

function handleOrderSuccess() {
  orderItems.value = []
  showOrderDialog.value = false
  ElMessage.success('订单创建成功！')
  loadOrders()
}

function handleSizeChange(size: number) {
  pageSize.value = size
  currentPage.value = 1
  loadOrders()
}

function handleCurrentChange(page: number) {
  currentPage.value = page
  loadOrders()
}
</script>

<style scoped>
.orders {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 20px 0;
}

.orders-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.orders-header h1 {
  margin: 0;
  color: #333;
}

.orders-content {
  background: white;
  border-radius: 12px;
  padding: 30px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.loading-section,
.empty-section {
  padding: 60px 20px;
  text-align: center;
}

.orders-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.order-card {
  border: 1px solid #eee;
  border-radius: 8px;
  overflow: hidden;
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: #f8f9fa;
  border-bottom: 1px solid #eee;
}

.order-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.order-id {
  font-weight: 500;
  color: #333;
}

.order-time {
  font-size: 0.9rem;
  color: #666;
}

.order-items {
  padding: 20px;
}

.order-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.order-item:last-child {
  border-bottom: none;
}

.item-image {
  width: 60px;
  height: 60px;
  border-radius: 4px;
  object-fit: cover;
}

.item-details {
  flex: 1;
}

.item-details h4 {
  margin: 0 0 4px 0;
  color: #333;
}

.item-details p {
  margin: 0;
  color: #666;
  font-size: 0.9rem;
}

.item-price {
  font-weight: bold;
  color: #f56c6c;
}

.order-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: #f8f9fa;
  border-top: 1px solid #eee;
}

.total-amount {
  font-size: 1.1rem;
  color: #333;
}

.total-amount span {
  color: #f56c6c;
  font-weight: bold;
  font-size: 1.2rem;
}

.order-actions {
  display: flex;
  gap: 8px;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

@media (max-width: 768px) {
  .orders-header {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;
  }
  
  .order-header {
    flex-direction: column;
    gap: 12px;
    align-items: flex-start;
  }
  
  .order-footer {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;
  }
  
  .order-item {
    flex-direction: column;
    align-items: flex-start;
    text-align: center;
  }
  
  .item-details {
    text-align: left;
    width: 100%;
  }
}
</style>