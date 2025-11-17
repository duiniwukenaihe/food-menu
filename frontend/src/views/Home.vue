<template>
  <div class="home">
    <!-- 头部导航 -->
    <header class="header">
      <div class="container">
        <div class="flex-between">
          <h1 class="logo polar-bear">🐻 美食点餐系统</h1>
          <div class="header-actions">
            <template v-if="userStore.isAuthenticated">
              <el-dropdown @command="handleCommand">
                <span class="user-info">
                  <el-icon><User /></el-icon>
                  {{ userStore.user?.username }}
                  <el-icon class="el-icon--right"><arrow-down /></el-icon>
                </span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                    <el-dropdown-item command="orders">我的订单</el-dropdown-item>
                    <el-dropdown-item command="favorites">我的收藏</el-dropdown-item>
                    <el-dropdown-item v-if="userStore.isAdmin" command="admin">管理后台</el-dropdown-item>
                    <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </template>
            <template v-else>
              <el-button type="primary" @click="showLoginDialog = true" class="login-btn">
                <el-icon><User /></el-icon>
                用户登录
              </el-button>
            </template>
          </div>
        </div>
      </div>
    </header>

    <!-- 主要内容 -->
    <main class="main">
      <div class="container">
        <!-- 欢迎信息 -->
        <div v-if="userStore.isAuthenticated" class="welcome-section">
          <div class="card text-center">
            <div class="bear-paw">🐾</div>
            <h2>欢迎点餐！</h2>
            <p>亲爱的 {{ userStore.user?.username }}，今天想吃什么美味呢？</p>
          </div>
        </div>

        <!-- 应季菜品推荐 -->
        <section v-if="dishStore.seasonalDishes.length > 0" class="seasonal-section">
          <h2 class="section-title">🍂 应季菜品推荐</h2>
          <div class="dish-grid">
            <div
              v-for="dish in dishStore.seasonalDishes"
              :key="dish.id"
              class="dish-card"
              @click="goToDishDetail(dish.id)"
            >
              <img :src="dish.image_url || '/placeholder-food.svg'" :alt="dish.name" class="dish-image" />
              <div class="dish-info">
                <h3>{{ dish.name }}</h3>
                <p class="price">¥{{ dish.price }}</p>
                <el-tag v-if="dish.is_seasonal" type="success" size="small">应季</el-tag>
              </div>
            </div>
          </div>
        </section>

        <!-- 推荐搭配 -->
        <section v-if="dishStore.recommendations.length > 0" class="recommendation-section">
          <h2 class="section-title">🍽️ 推荐搭配</h2>
          <div class="recommendation-list">
            <div
              v-for="rec in dishStore.recommendations"
              :key="rec.id"
              class="recommendation-card"
            >
              <h3>{{ rec.name }}</h3>
              <p>{{ rec.description }}</p>
              <p class="config">{{ rec.meat_count }}荤{{ rec.vegetable_count }}素</p>
              <el-button type="primary" @click="generateRandomCombo(rec)">随机搭配</el-button>
            </div>
          </div>
        </section>

        <!-- 猜你喜欢 -->
        <section class="random-section">
          <h2 class="section-title">🎲 猜你喜欢</h2>
          <div class="flex-between mb-20">
            <el-select v-model="meatCount" placeholder="荤菜数量" style="width: 120px">
              <el-option label="0荤" :value="0" />
              <el-option label="1荤" :value="1" />
              <el-option label="2荤" :value="2" />
              <el-option label="3荤" :value="3" />
            </el-select>
            <el-select v-model="vegetableCount" placeholder="素菜数量" style="width: 120px">
              <el-option label="1素" :value="1" />
              <el-option label="2素" :value="2" />
              <el-option label="3素" :value="3" />
            </el-select>
            <el-button type="primary" @click="generateRandomDishes">
              <el-icon><Refresh /></el-icon>
              随机生成
            </el-button>
          </div>
          <div v-if="randomDishes.length > 0" class="dish-grid">
            <div
              v-for="dish in randomDishes"
              :key="dish.id"
              class="dish-card"
              @click="goToDishDetail(dish.id)"
            >
              <img :src="dish.image_url || '/placeholder-food.svg'" :alt="dish.name" class="dish-image" />
              <div class="dish-info">
                <h3>{{ dish.name }}</h3>
                <p class="price">¥{{ dish.price }}</p>
                <el-button v-if="userStore.isAuthenticated" type="primary" size="small" @click.stop="addToOrder(dish)">
                  加入订单
                </el-button>
              </div>
            </div>
          </div>
        </section>

        <!-- 所有菜品 -->
        <section class="all-dishes-section">
          <h2 class="section-title">📋 所有菜品</h2>
          
          <!-- 筛选器 -->
          <div class="filters mb-20">
            <el-select v-model="selectedCategory" placeholder="选择分类" clearable style="width: 200px">
              <el-option
                v-for="category in dishStore.categories"
                :key="category.id"
                :label="category.name"
                :value="category.id"
              />
            </el-select>
            <el-input
              v-model="searchQuery"
              placeholder="搜索菜品"
              style="width: 300px; margin-left: 10px"
              clearable
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            <el-button type="primary" @click="loadDishes" style="margin-left: 10px">
              搜索
            </el-button>
          </div>

          <!-- 菜品网格 -->
          <div v-loading="dishStore.isLoading" class="dish-grid">
            <div
              v-for="dish in dishStore.dishes"
              :key="dish.id"
              class="dish-card"
              @click="goToDishDetail(dish.id)"
            >
              <img :src="dish.image_url || '/placeholder-food.svg'" :alt="dish.name" class="dish-image" />
              <div class="dish-info">
                <h3>{{ dish.name }}</h3>
                <p class="description">{{ dish.description }}</p>
                <p class="price">¥{{ dish.price }}</p>
                <div class="dish-tags">
                  <el-tag v-if="dish.is_seasonal" type="success" size="small">应季</el-tag>
                  <el-tag v-if="dish.category" size="small">{{ dish.category.name }}</el-tag>
                </div>
                <el-button v-if="userStore.isAuthenticated" type="primary" size="small" @click.stop="addToOrder(dish)">
                  加入订单
                </el-button>
              </div>
            </div>
          </div>

          <!-- 分页 -->
          <div class="pagination-wrapper">
            <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :total="dishStore.total"
              :page-sizes="[12, 24, 48]"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
            />
          </div>
        </section>
      </div>
    </main>

    <!-- 登录对话框 -->
    <LoginDialog v-model="showLoginDialog" @success="handleLoginSuccess" />

    <!-- 订单对话框 -->
    <OrderDialog v-model="showOrderDialog" :items="orderItems" @success="handleOrderSuccess" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, ArrowDown, Refresh, Search } from '@element-plus/icons-vue'
import { useUserStore, useDishStore, useOrderStore } from '@/stores'
import { api } from '@/utils/api'
import LoginDialog from '@/components/LoginDialog.vue'
import OrderDialog from '@/components/OrderDialog.vue'
import type { Dish, Recommendation } from '@/types'

const router = useRouter()
const userStore = useUserStore()
const dishStore = useDishStore()
const orderStore = useOrderStore()

// 响应式数据
const showLoginDialog = ref(false)
const showOrderDialog = ref(false)
const orderItems = ref<{ dish: Dish; quantity: number }[]>([])
const randomDishes = ref<Dish[]>([])

// 筛选和分页
const selectedCategory = ref<number>()
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(12)

// 随机搭配配置
const meatCount = ref(1)
const vegetableCount = ref(2)

// 计算属性
const meatDishes = computed(() => 
  dishStore.dishes.filter(dish => dish.category?.name === '肉类')
)

const vegetableDishes = computed(() => 
  dishStore.dishes.filter(dish => dish.category?.name === '蔬菜类')
)

// 生命周期
onMounted(async () => {
  await Promise.all([
    dishStore.fetchCategories(),
    dishStore.fetchSeasonalDishes(),
    dishStore.fetchRecommendations(),
    loadDishes()
  ])
})

// 监听筛选条件变化
watch([selectedCategory, searchQuery], () => {
  currentPage.value = 1
  loadDishes()
})

// 方法
async function loadDishes() {
  await dishStore.fetchDishes({
    page: currentPage.value,
    limit: pageSize.value,
    category_id: selectedCategory.value,
    search: searchQuery.value
  })
}

function goToDishDetail(id: number) {
  router.push(`/dish/${id}`)
}

function handleCommand(command: string) {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'orders':
      router.push('/orders')
      break
    case 'favorites':
      router.push('/favorites')
      break
    case 'admin':
      router.push('/admin')
      break
    case 'logout':
      userStore.logout()
      ElMessage.success('退出登录成功')
      break
  }
}

function handleLoginSuccess() {
  showLoginDialog.value = false
  ElMessage.success('登录成功！')
}

function addToOrder(dish: Dish) {
  const existingItem = orderItems.value.find(item => item.dish.id === dish.id)
  if (existingItem) {
    existingItem.quantity++
  } else {
    orderItems.value.push({ dish, quantity: 1 })
  }
  ElMessage.success(`已添加 ${dish.name} 到订单`)
  showOrderDialog.value = true
}

function handleOrderSuccess() {
  orderItems.value = []
  showOrderDialog.value = false
  ElMessage.success('订单创建成功！')
}

async function generateRandomCombo(recommendation: Recommendation) {
  const meatDishesList = meatDishes.value
  const vegetableDishesList = vegetableDishes.value
  
  const selectedMeats: Dish[] = []
  const selectedVegetables: Dish[] = []
  
  // 随机选择荤菜
  for (let i = 0; i < recommendation.meat_count && i < meatDishesList.length; i++) {
    const randomIndex = Math.floor(Math.random() * meatDishesList.length)
    const dish = meatDishesList[randomIndex]
    if (!selectedMeats.find(d => d.id === dish.id)) {
      selectedMeats.push(dish)
    }
  }
  
  // 随机选择素菜
  for (let i = 0; i < recommendation.vegetable_count && i < vegetableDishesList.length; i++) {
    const randomIndex = Math.floor(Math.random() * vegetableDishesList.length)
    const dish = vegetableDishesList[randomIndex]
    if (!selectedVegetables.find(d => d.id === dish.id)) {
      selectedVegetables.push(dish)
    }
  }
  
  randomDishes.value = [...selectedMeats, ...selectedVegetables]
}

async function generateRandomDishes() {
  const meatDishesList = meatDishes.value
  const vegetableDishesList = vegetableDishes.value
  
  const selectedMeats: Dish[] = []
  const selectedVegetables: Dish[] = []
  
  // 随机选择荤菜
  for (let i = 0; i < meatCount.value && i < meatDishesList.length; i++) {
    const randomIndex = Math.floor(Math.random() * meatDishesList.length)
    const dish = meatDishesList[randomIndex]
    if (!selectedMeats.find(d => d.id === dish.id)) {
      selectedMeats.push(dish)
    }
  }
  
  // 随机选择素菜
  for (let i = 0; i < vegetableCount.value && i < vegetableDishesList.length; i++) {
    const randomIndex = Math.floor(Math.random() * vegetableDishesList.length)
    const dish = vegetableDishesList[randomIndex]
    if (!selectedVegetables.find(d => d.id === dish.id)) {
      selectedVegetables.push(dish)
    }
  }
  
  randomDishes.value = [...selectedMeats, ...selectedVegetables]
}

function handleSizeChange(size: number) {
  pageSize.value = size
  currentPage.value = 1
  loadDishes()
}

function handleCurrentChange(page: number) {
  currentPage.value = page
  loadDishes()
}
</script>

<style scoped>
.home {
  min-height: 100vh;
}

.header {
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 1rem 0;
  position: sticky;
  top: 0;
  z-index: 100;
}

.logo {
  font-size: 1.5rem;
  margin: 0;
}

.user-info {
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: #f5f5f5;
}

.login-btn {
  border-radius: 20px;
}

.main {
  padding: 2rem 0;
}

.welcome-section {
  margin-bottom: 2rem;
}

.bear-paw {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.section-title {
  margin-bottom: 1.5rem;
  color: #333;
  font-size: 1.5rem;
}

.dish-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.dish-card {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s, box-shadow 0.3s;
  cursor: pointer;
}

.dish-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.dish-image {
  width: 100%;
  height: 200px;
  object-fit: cover;
}

.dish-info {
  padding: 1rem;
}

.dish-info h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.1rem;
  color: #333;
}

.description {
  color: #666;
  font-size: 0.9rem;
  margin-bottom: 0.5rem;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.price {
  color: #f56c6c;
  font-weight: bold;
  font-size: 1.2rem;
  margin: 0.5rem 0;
}

.dish-tags {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.recommendation-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}

.recommendation-card {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.recommendation-card h3 {
  margin-bottom: 0.5rem;
  color: #333;
}

.config {
  color: #666;
  margin: 1rem 0;
}

.filters {
  display: flex;
  align-items: center;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 2rem;
}
</style>