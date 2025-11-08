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
  isPublished: boolean
  viewCount: number
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

export interface Dish {
  id: number
  name: string
  description: string
  tags: string
  isActive: boolean
  isSeasonal: boolean
  availableMonths: string
  seasonalNote: string
  imageUrl: string
  thumbnailUrl: string
  galleryUrls: string
  createdAt: string
  updatedAt: string
}

export interface CreateDishRequest {
  name: string
  description?: string
  tags?: string
  isActive?: boolean
  isSeasonal?: boolean
  availableMonths?: string
  seasonalNote?: string
  imageUrl?: string
  thumbnailUrl?: string
  galleryUrls?: string
}

export interface UpdateDishRequest {
  name?: string
  description?: string
  tags?: string
  isActive?: boolean
  isSeasonal?: boolean
  availableMonths?: string
  seasonalNote?: string
  imageUrl?: string
  thumbnailUrl?: string
  galleryUrls?: string
}

export interface MediaFile {
  url: string
  thumbnailUrl?: string
  fileName: string
  contentType: string
  size: number
}

export interface UploadURLResponse {
  uploadUrl: string
  fileUrl: string
  key: string
}
