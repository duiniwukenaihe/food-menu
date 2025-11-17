export interface User {
  id: number
  username: string
  email: string
  role: string
  created_at: string
  updated_at: string
}

export interface Category {
  id: number
  name: string
  description: string
  created_at: string
}

export interface Dish {
  id: number
  name: string
  description: string
  category_id: number
  category?: Category
  price: number
  image_url: string
  video_url: string
  cooking_steps: string
  is_seasonal: boolean
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface Order {
  id: number
  user_id: number
  total_amount: number
  status: 'pending' | 'confirmed' | 'preparing' | 'ready' | 'completed' | 'cancelled'
  created_at: string
  updated_at: string
  items?: OrderItem[]
}

export interface OrderItem {
  id: number
  order_id: number
  dish_id: number
  dish?: Dish
  quantity: number
  price: number
  created_at: string
}

export interface Recommendation {
  id: number
  name: string
  description: string
  meat_count: number
  vegetable_count: number
  is_active: boolean
  created_at: string
  dishes?: Dish[]
}

export interface UserFavorite {
  id: number
  user_id: number
  dish_id: number
  dish?: Dish
  created_at: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface CreateOrderRequest {
  items: CreateOrderItemRequest[]
}

export interface CreateOrderItemRequest {
  dish_id: number
  quantity: number
}

export interface CreateDishRequest {
  name: string
  description: string
  category_id: number
  price: number
  image_url: string
  video_url: string
  cooking_steps: string
  is_seasonal: boolean
}

export interface UpdateDishRequest {
  name?: string
  description?: string
  category_id?: number
  price?: number
  image_url?: string
  video_url?: string
  cooking_steps?: string
  is_seasonal?: boolean
  is_active?: boolean
}

export interface CreateCategoryRequest {
  name: string
  description: string
}

export interface UpdateCategoryRequest {
  name?: string
  description?: string
}

export interface SystemConfig {
  id: number
  config_key: string
  config_value: string
  description: string
  updated_at: string
}

export interface ApiResponse<T = any> {
  data?: T
  error?: string
  message?: string
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  limit: number
}