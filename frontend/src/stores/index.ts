import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, Dish, Category, Order, Recommendation } from '@/types'
import { api } from '@/utils/api'

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const isLoading = ref(false)

  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  // 初始化用户状态
  function initAuth() {
    const savedToken = localStorage.getItem('token')
    const savedUser = localStorage.getItem('user')
    
    if (savedToken && savedUser) {
      token.value = savedToken
      user.value = JSON.parse(savedUser)
    }
  }

  // 登录
  async function login(username: string, password: string) {
    try {
      isLoading.value = true
      const response = await api.login({ username, password })
      
      token.value = response.token
      user.value = response.user
      
      localStorage.setItem('token', response.token)
      localStorage.setItem('user', JSON.stringify(response.user))
      
      return true
    } catch (error) {
      console.error('Login failed:', error)
      return false
    } finally {
      isLoading.value = false
    }
  }

  // 登出
  function logout() {
    user.value = null
    token.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  // 获取用户信息
  async function fetchProfile() {
    try {
      const profile = await api.getProfile()
      user.value = profile
      localStorage.setItem('user', JSON.stringify(profile))
    } catch (error) {
      console.error('Failed to fetch profile:', error)
      logout()
    }
  }

  return {
    user,
    token,
    isLoading,
    isAuthenticated,
    isAdmin,
    initAuth,
    login,
    logout,
    fetchProfile,
  }
})

export const useDishStore = defineStore('dish', () => {
  const dishes = ref<Dish[]>([])
  const categories = ref<Category[]>([])
  const seasonalDishes = ref<Dish[]>([])
  const recommendations = ref<Recommendation[]>([])
  const isLoading = ref(false)
  const total = ref(0)

  // 获取菜品列表
  async function fetchDishes(params?: {
    page?: number
    limit?: number
    category_id?: number
    search?: string
  }) {
    try {
      isLoading.value = true
      const response = await api.getDishes(params)
      dishes.value = response.dishes || []
      total.value = response.total || 0
    } catch (error) {
      console.error('Failed to fetch dishes:', error)
    } finally {
      isLoading.value = false
    }
  }

  // 获取分类列表
  async function fetchCategories() {
    try {
      categories.value = await api.getCategories()
    } catch (error) {
      console.error('Failed to fetch categories:', error)
    }
  }

  // 获取应季菜品
  async function fetchSeasonalDishes() {
    try {
      seasonalDishes.value = await api.getSeasonalDishes()
    } catch (error) {
      console.error('Failed to fetch seasonal dishes:', error)
    }
  }

  // 获取推荐配置
  async function fetchRecommendations() {
    try {
      recommendations.value = await api.getRecommendations()
    } catch (error) {
      console.error('Failed to fetch recommendations:', error)
    }
  }

  // 获取单个菜品
  async function getDish(id: number): Promise<Dish | null> {
    try {
      return await api.getDish(id)
    } catch (error) {
      console.error('Failed to get dish:', error)
      return null
    }
  }

  return {
    dishes,
    categories,
    seasonalDishes,
    recommendations,
    isLoading,
    total,
    fetchDishes,
    fetchCategories,
    fetchSeasonalDishes,
    fetchRecommendations,
    getDish,
  }
})

export const useOrderStore = defineStore('order', () => {
  const orders = ref<Order[]>([])
  const favorites = ref<Dish[]>([])
  const isLoading = ref(false)
  const total = ref(0)

  // 创建订单
  async function createOrder(items: { dish_id: number; quantity: number }[]) {
    try {
      isLoading.value = true
      const order = await api.createOrder({ items })
      orders.value.unshift(order)
      return order
    } catch (error) {
      console.error('Failed to create order:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  // 获取订单列表
  async function fetchOrders(params?: { page?: number; limit?: number }) {
    try {
      isLoading.value = true
      const response = await api.getOrders(params)
      orders.value = response.orders || []
      total.value = response.total || 0
    } catch (error) {
      console.error('Failed to fetch orders:', error)
    } finally {
      isLoading.value = false
    }
  }

  // 获取收藏列表
  async function fetchFavorites(params?: { page?: number; limit?: number }) {
    try {
      isLoading.value = true
      const response = await api.getFavorites(params)
      favorites.value = response.favorites?.map(fav => fav.dish!).filter(Boolean) || []
      total.value = response.total || 0
    } catch (error) {
      console.error('Failed to fetch favorites:', error)
    } finally {
      isLoading.value = false
    }
  }

  // 添加到收藏
  async function addToFavorites(dishId: number) {
    try {
      await api.addToFavorites(dishId)
      // 重新获取收藏列表
      await fetchFavorites()
    } catch (error) {
      console.error('Failed to add to favorites:', error)
      throw error
    }
  }

  // 从收藏中移除
  async function removeFromFavorites(dishId: number) {
    try {
      await api.removeFromFavorites(dishId)
      // 重新获取收藏列表
      await fetchFavorites()
    } catch (error) {
      console.error('Failed to remove from favorites:', error)
      throw error
    }
  }

  return {
    orders,
    favorites,
    isLoading,
    total,
    createOrder,
    fetchOrders,
    fetchFavorites,
    addToFavorites,
    removeFromFavorites,
  }
})