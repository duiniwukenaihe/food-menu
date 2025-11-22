<template>
  <div class="favorites">
    <div class="container">
      <div class="favorites-header">
        <h1>我的收藏</h1>
        <el-button @click="goBack" type="default">
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
      </div>

      <div class="favorites-content">
        <div v-if="orderStore.isLoading" class="loading-section">
          <el-skeleton :rows="5" animated />
        </div>

        <div v-else-if="orderStore.favorites.length === 0" class="empty-section">
          <el-empty description="暂无收藏的菜品">
            <el-button type="primary" @click="goToHome">去浏览菜品</el-button>
          </el-empty>
        </div>

        <div v-else class="favorites-grid">
          <div
            v-for="dish in orderStore.favorites"
            :key="dish.id"
            class="favorite-card"
          >
            <img 
              :src="dish.image_url || '/placeholder-food.svg'" 
              :alt="dish.name" 
              class="dish-image"
              @click="goToDishDetail(dish.id)"
            />
            <div class="dish-info">
              <h3 @click="goToDishDetail(dish.id)">{{ dish.name }}</h3>
              <p class="description">{{ dish.description }}</p>
              <p class="price">¥{{ dish.price }}</p>
              <div class="dish-tags">
                <el-tag v-if="dish.is_seasonal" type="success" size="small">应季</el-tag>
                <el-tag v-if="dish.category" size="small">{{ dish.category.name }}</el-tag>
              </div>
              <div class="action-buttons">
                <el-button type="primary" size="small" @click="addToOrder(dish)">
                  加入订单
                </el-button>
                <el-button type="danger" size="small" @click="removeFromFavorites(dish.id)">
                  取消收藏
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
            :page-sizes="[12, 24, 48]"
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
import { ElMessage } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import { useOrderStore } from '@/stores'
import OrderDialog from '@/components/OrderDialog.vue'
import type { Dish } from '@/types'

const router = useRouter()
const orderStore = useOrderStore()

const currentPage = ref(1)
const pageSize = ref(12)
const showOrderDialog = ref(false)
const orderItems = ref<{ dish: Dish; quantity: number }[]>([])

onMounted(() => {
  loadFavorites()
})

async function loadFavorites() {
  await orderStore.fetchFavorites({
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

function goToDishDetail(id: number) {
  router.push(`/dish/${id}`)
}

async function removeFromFavorites(dishId: number) {
  try {
    await orderStore.removeFromFavorites(dishId)
    ElMessage.success('已取消收藏')
  } catch (error) {
    console.error('Remove from favorites failed:', error)
    ElMessage.error('取消收藏失败，请重试')
  }
}

function addToOrder(dish: Dish) {
  orderItems.value = [{ dish, quantity: 1 }]
  showOrderDialog.value = true
}

function handleOrderSuccess() {
  orderItems.value = []
  showOrderDialog.value = false
  ElMessage.success('订单创建成功！')
}

function handleSizeChange(size: number) {
  pageSize.value = size
  currentPage.value = 1
  loadFavorites()
}

function handleCurrentChange(page: number) {
  currentPage.value = page
  loadFavorites()
}
</script>

<style scoped>
.favorites {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 20px 0;
}

.favorites-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.favorites-header h1 {
  margin: 0;
  color: #333;
}

.favorites-content {
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

.favorites-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 24px;
  margin-bottom: 30px;
}

.favorite-card {
  background: white;
  border: 1px solid #eee;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s, box-shadow 0.3s;
}

.favorite-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.dish-image {
  width: 100%;
  height: 200px;
  object-fit: cover;
  cursor: pointer;
}

.dish-info {
  padding: 20px;
}

.dish-info h3 {
  margin: 0 0 8px 0;
  font-size: 1.2rem;
  color: #333;
  cursor: pointer;
  transition: color 0.3s;
}

.dish-info h3:hover {
  color: #409eff;
}

.description {
  color: #666;
  font-size: 0.9rem;
  margin-bottom: 12px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.4;
}

.price {
  color: #f56c6c;
  font-weight: bold;
  font-size: 1.3rem;
  margin: 12px 0;
}

.dish-tags {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.action-buttons {
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
  .favorites-header {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;
  }
  
  .favorites-grid {
    grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
    gap: 16px;
  }
  
  .favorites-content {
    padding: 20px;
  }
}
</style>