import axios, { AxiosInstance, AxiosResponse } from 'axios'
import type { 
  User, 
  Dish, 
  Category, 
  Order, 
  Recommendation, 
  UserFavorite,
  LoginRequest,
  LoginResponse,
  CreateOrderRequest,
  CreateDishRequest,
  UpdateDishRequest,
  CreateCategoryRequest,
  UpdateCategoryRequest,
  SystemConfig,
  ApiResponse,
  PaginatedResponse
} from '@/types'

class ApiClient {
  private client: AxiosInstance

  constructor() {
    this.client = axios.create({
      baseURL: '/api/v1',
      timeout: 10000,
    })

    // 请求拦截器 - 添加token
    this.client.interceptors.request.use(
      (config) => {
        const token = localStorage.getItem('token')
        if (token) {
          config.headers.Authorization = `Bearer ${token}`
        }
        return config
      },
      (error) => {
        return Promise.reject(error)
      }
    )

    // 响应拦截器 - 处理错误
    this.client.interceptors.response.use(
      (response: AxiosResponse) => {
        return response
      },
      (error) => {
        if (error.response?.status === 401) {
          localStorage.removeItem('token')
          localStorage.removeItem('user')
          window.location.href = '/login'
        }
        return Promise.reject(error)
      }
    )
  }

  // 认证相关
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await this.client.post<LoginResponse>('/login', credentials)
    return response.data
  }

  async getProfile(): Promise<User> {
    const response = await this.client.get<User>('/profile')
    return response.data
  }

  // 菜品相关
  async getDishes(params?: {
    page?: number
    limit?: number
    category_id?: number
    search?: string
  }): Promise<PaginatedResponse<Dish>> {
    const response = await this.client.get<PaginatedResponse<Dish>>('/dishes', { params })
    return response.data
  }

  async getDish(id: number): Promise<Dish> {
    const response = await this.client.get<Dish>(`/dishes/${id}`)
    return response.data
  }

  async createDish(dish: CreateDishRequest): Promise<Dish> {
    const response = await this.client.post<Dish>('/admin/dishes', dish)
    return response.data
  }

  async updateDish(id: number, dish: UpdateDishRequest): Promise<Dish> {
    const response = await this.client.put<Dish>(`/admin/dishes/${id}`, dish)
    return response.data
  }

  async deleteDish(id: number): Promise<void> {
    await this.client.delete(`/admin/dishes/${id}`)
  }

  // 分类相关
  async getCategories(): Promise<Category[]> {
    const response = await this.client.get<Category[]>('/categories')
    return response.data
  }

  async createCategory(category: CreateCategoryRequest): Promise<Category> {
    const response = await this.client.post<Category>('/admin/categories', category)
    return response.data
  }

  async updateCategory(id: number, category: UpdateCategoryRequest): Promise<Category> {
    const response = await this.client.put<Category>(`/admin/categories/${id}`, category)
    return response.data
  }

  async deleteCategory(id: number): Promise<void> {
    await this.client.delete(`/admin/categories/${id}`)
  }

  // 推荐相关
  async getRecommendations(): Promise<Recommendation[]> {
    const response = await this.client.get<Recommendation[]>('/recommendations')
    return response.data
  }

  async getSeasonalDishes(): Promise<Dish[]> {
    const response = await this.client.get<Dish[]>('/seasonal-dishes')
    return response.data
  }

  // 订单相关
  async createOrder(order: CreateOrderRequest): Promise<Order> {
    const response = await this.client.post<Order>('/orders', order)
    return response.data
  }

  async getOrders(params?: { page?: number; limit?: number }): Promise<PaginatedResponse<Order>> {
    const response = await this.client.get<PaginatedResponse<Order>>('/orders', { params })
    return response.data
  }

  // 收藏相关
  async addToFavorites(dishId: number): Promise<void> {
    await this.client.post(`/favorites/${dishId}`)
  }

  async removeFromFavorites(dishId: number): Promise<void> {
    await this.client.delete(`/favorites/${dishId}`)
  }

  async getFavorites(params?: { page?: number; limit?: number }): Promise<PaginatedResponse<UserFavorite>> {
    const response = await this.client.get<PaginatedResponse<UserFavorite>>('/favorites', { params })
    return response.data
  }

  // 管理员相关
  async getUsers(params?: { page?: number; limit?: number; search?: string }): Promise<PaginatedResponse<User>> {
    const response = await this.client.get<PaginatedResponse<User>>('/admin/users', { params })
    return response.data
  }

  async getConfig(): Promise<SystemConfig[]> {
    const response = await this.client.get<SystemConfig[]>('/admin/config')
    return response.data
  }

  async updateConfig(config: Record<string, string>): Promise<void> {
    await this.client.put('/admin/config', config)
  }
}

export const api = new ApiClient()