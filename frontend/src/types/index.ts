export interface User {
  id: number
  username: string
  email: string
  firstName: string
  lastName: string
  role: string
  avatar?: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface AuthResponse {
  token: string
  user: User
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
  firstName: string
  lastName: string
}

export interface Category {
  id: number
  name: string
  description?: string
  color?: string
  icon?: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface Content {
  id: number
  title: string
  description?: string
  body: string
  categoryId: number
  category: Category
  authorId: number
  author: User
  tags?: string
  imageUrl?: string
  videoUrl?: string
  isPublished: boolean
  viewCount: number
  difficulty?: string
  prepTime?: string
  servings?: string
  rating?: number
  createdAt: string
  updatedAt: string
}

export interface CreateContentRequest {
  title: string
  description?: string
  body: string
  categoryId: number
  tags?: string
  imageUrl?: string
  isPublished?: boolean
}

export interface UpdateContentRequest {
  title?: string
  description?: string
  body?: string
  categoryId?: number
  tags?: string
  imageUrl?: string
  isPublished?: boolean
}

export interface PaginatedResponse<T> {
  success: boolean
  data: T[]
  total: number
  page: number
  limit: number
}

export interface SuccessResponse<T> {
  success: boolean
  message?: string
  data?: T
}

export interface ErrorResponse {
  success: boolean
  message: string
  error?: string
}
