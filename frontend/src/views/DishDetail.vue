<template>
  <div class="dish-detail">
    <div class="container">
      <!-- 返回按钮 -->
      <div class="back-section">
        <el-button @click="goBack" type="default">
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
      </div>

      <!-- 菜品详情 -->
      <div v-if="dish" class="dish-content">
        <div class="dish-header">
          <div class="dish-image-section">
            <img 
              :src="dish.image_url || '/placeholder-food.jpg'" 
              :alt="dish.name" 
              class="main-dish-image"
            />
          </div>
          <div class="dish-info-section">
            <h1 class="dish-title">{{ dish.name }}</h1>
            <div class="dish-meta">
              <el-tag v-if="dish.is_seasonal" type="success" size="large">应季推荐</el-tag>
              <el-tag v-if="dish.category" type="info" size="large">{{ dish.category.name }}</el-tag>
            </div>
            <div class="price-section">
              <span class="price">¥{{ dish.price }}</span>
            </div>
            <div class="action-buttons">
              <el-button 
                v-if="userStore.isAuthenticated"
                type="primary" 
                size="large"
                @click="addToOrder"
              >
                <el-icon><ShoppingCart /></el-icon>
                加入订单
              </el-button>
              <el-button 
                v-else
                type="primary" 
                size="large"
                @click="goToLogin"
              >
                <el-icon><User /></el-icon>
                登录后点餐
              </el-button>
              
              <el-button 
                v-if="userStore.isAuthenticated"
                :type="isFavorite ? 'danger' : 'default'"
                size="large"
                @click="toggleFavorite"
              >
                <el-icon><Star /></el-icon>
                {{ isFavorite ? '取消收藏' : '收藏' }}
              </el-button>
            </div>
          </div>
        </div>

        <!-- 菜品描述 -->
        <div v-if="dish.description" class="description-section">
          <h2>菜品介绍</h2>
          <p class="description">{{ dish.description }}</p>
        </div>

        <!-- 制作步骤 -->
        <div v-if="dish.cooking_steps" class="cooking-steps-section">
          <h2>制作步骤</h2>
          <div class="steps-content">
            <div v-for="(step, index) in cookingSteps" :key="index" class="step-item">
              <div class="step-number">{{ index + 1 }}</div>
              <div class="step-text">{{ step }}</div>
            </div>
          </div>
        </div>

        <!-- 制作视频 -->
        <div v-if="dish.video_url" class="video-section">
          <h2>制作视频</h2>
          <div class="video-container">
            <video 
              :src="dish.video_url" 
              controls 
              preload="metadata"
              class="cooking-video"
            >
              您的浏览器不支持视频播放
            </video>
          </div>
        </div>

        <!-- 营养信息 -->
        <div class="nutrition-section">
          <h2>营养信息</h2>
          <div class="nutrition-grid">
            <div class="nutrition-item">
              <span class="nutrition-label">热量</span>
              <span class="nutrition-value">{{ nutrition?.calories || '--' }} kcal</span>
            </div>
            <div class="nutrition-item">
              <span class="nutrition-label">蛋白质</span>
              <span class="nutrition-value">{{ nutrition?.protein || '--' }}g</span>
            </div>
            <div class="nutrition-item">
              <span class="nutrition-label">脂肪</span>
              <span class="nutrition-value">{{ nutrition?.fat || '--' }}g</span>
            </div>
            <div class="nutrition-item">
              <span class="nutrition-label">碳水化合物</span>
              <span class="nutrition-value">{{ nutrition?.carbohydrates || '--' }}g</span>
            </div>
            <div class="nutrition-item">
              <span class="nutrition-label">纤维</span>
              <span class="nutrition-value">{{ nutrition?.fiber || '--' }}g</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 加载状态 -->
      <div v-else-if="loading" class="loading-section">
        <el-skeleton :rows="10" animated />
      </div>

      <!-- 错误状态 -->
      <div v-else class="error-section">
        <el-result
          icon="error"
          title="菜品不存在"
          sub-title="抱歉，您要查看的菜品不存在或已被删除"
        >
          <template #extra>
            <el-button type="primary" @click="goBack">返回首页</el-button>
          </template>
        </el-result>
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
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft, ShoppingCart, User, Star } from '@element-plus/icons-vue'
import { useUserStore, useDishStore, useOrderStore } from '@/stores'
import { api } from '@/utils/api'
import OrderDialog from '@/components/OrderDialog.vue'
import type { Dish } from '@/types'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const dishStore = useDishStore()
const orderStore = useOrderStore()

const dish = ref<Dish | null>(null)
const loading = ref(true)
const nutrition = ref<any>(null)
const showOrderDialog = ref(false)
const orderItems = ref<{ dish: Dish; quantity: number }[]>([])

const dishId = computed(() => parseInt(route.params.id as string))

const isFavorite = computed(() => {
  return orderStore.favorites.some(fav => fav.id === dish.value?.id)
})

const cookingSteps = computed(() => {
  if (!dish.value?.cooking_steps) return []
  return dish.value.cooking_steps.split(/[\n\r]+/).filter(step => step.trim())
})

onMounted(async () => {
  await fetchDishDetail()
  if (userStore.isAuthenticated) {
    await orderStore.fetchFavorites()
  }
})

async function fetchDishDetail() {
  try {
    loading.value = true
    dish.value = await dishStore.getDish(dishId.value)
    
    if (dish.value) {
      // 这里可以获取营养信息，但需要后端支持
      // nutrition.value = await api.getDishNutrition(dishId.value)
    }
  } catch (error) {
    console.error('Failed to fetch dish detail:', error)
    dish.value = null
  } finally {
    loading.value = false
  }
}

function goBack() {
  router.go(-1)
}

function goToLogin() {
  router.push('/login')
}

function addToOrder() {
  if (!dish.value) return
  
  orderItems.value = [{ dish: dish.value, quantity: 1 }]
  showOrderDialog.value = true
}

function handleOrderSuccess() {
  orderItems.value = []
  showOrderDialog.value = false
  ElMessage.success('订单创建成功！')
}

async function toggleFavorite() {
  if (!dish.value) return
  
  try {
    if (isFavorite.value) {
      await orderStore.removeFromFavorites(dish.value.id)
      ElMessage.success('已取消收藏')
    } else {
      await orderStore.addToFavorites(dish.value.id)
      ElMessage.success('已添加到收藏')
    }
  } catch (error) {
    console.error('Toggle favorite failed:', error)
    ElMessage.error('操作失败，请重试')
  }
}
</script>

<style scoped>
.dish-detail {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 20px 0;
}

.back-section {
  margin-bottom: 20px;
}

.dish-content {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.dish-header {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 40px;
  padding: 40px;
}

.dish-image-section {
  display: flex;
  justify-content: center;
  align-items: center;
}

.main-dish-image {
  width: 100%;
  max-width: 400px;
  height: 300px;
  object-fit: cover;
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.dish-info-section {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.dish-title {
  font-size: 2.5rem;
  color: #333;
  margin-bottom: 16px;
  line-height: 1.2;
}

.dish-meta {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
}

.price-section {
  margin-bottom: 32px;
}

.price {
  font-size: 2rem;
  color: #f56c6c;
  font-weight: bold;
}

.action-buttons {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.description-section,
.cooking-steps-section,
.video-section,
.nutrition-section {
  padding: 32px 40px;
  border-top: 1px solid #eee;
}

.description-section h2,
.cooking-steps-section h2,
.video-section h2,
.nutrition-section h2 {
  font-size: 1.5rem;
  color: #333;
  margin-bottom: 16px;
}

.description {
  font-size: 1.1rem;
  line-height: 1.6;
  color: #666;
}

.steps-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.step-item {
  display: flex;
  align-items: flex-start;
  gap: 16px;
}

.step-number {
  width: 32px;
  height: 32px;
  background: #409eff;
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  flex-shrink: 0;
}

.step-text {
  flex: 1;
  line-height: 1.6;
  color: #333;
  padding-top: 4px;
}

.video-container {
  position: relative;
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
}

.cooking-video {
  width: 100%;
  height: auto;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.nutrition-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 20px;
}

.nutrition-item {
  text-align: center;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
}

.nutrition-label {
  display: block;
  font-size: 0.9rem;
  color: #666;
  margin-bottom: 8px;
}

.nutrition-value {
  display: block;
  font-size: 1.2rem;
  font-weight: bold;
  color: #333;
}

.loading-section,
.error-section {
  padding: 40px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

@media (max-width: 768px) {
  .dish-header {
    grid-template-columns: 1fr;
    gap: 20px;
    padding: 20px;
  }
  
  .dish-title {
    font-size: 2rem;
  }
  
  .price {
    font-size: 1.5rem;
  }
  
  .description-section,
  .cooking-steps-section,
  .video-section,
  .nutrition-section {
    padding: 20px;
  }
}
</style>