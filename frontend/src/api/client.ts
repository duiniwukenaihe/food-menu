import axios, { AxiosResponse } from 'axios'
import toast from 'react-hot-toast'
import { AuthResponse, ErrorResponse } from '../types'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor to add auth token
apiClient.interceptors.request.use(
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

// Response interceptor to handle errors
apiClient.interceptors.response.use(
  (response: AxiosResponse) => {
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
      toast.error('Session expired. Please login again.')
    } else if (error.response?.data?.message) {
      toast.error(error.response.data.message)
    } else {
      toast.error('An error occurred. Please try again.')
    }
    return Promise.reject(error)
  }
)

export const authApi = {
  login: async (credentials: { username: string; password: string }): Promise<AuthResponse> => {
    const response = await apiClient.post('/auth/login', credentials)
    return response.data
  },

  register: async (userData: {
    username: string
    email: string
    password: string
    firstName: string
    lastName: string
  }): Promise<AuthResponse> => {
    const response = await apiClient.post('/auth/register', userData)
    return response.data
  },

  getProfile: async (): Promise<{ success: boolean; data: any }> => {
    const response = await apiClient.get('/auth/profile')
    return response.data
  },

  updateProfile: async (userData: {
    firstName?: string
    lastName?: string
    avatar?: string
  }): Promise<{ success: boolean; data: any }> => {
    const response = await apiClient.put('/auth/profile', userData)
    return response.data
  },
}

export const contentApi = {
  getContent: async (params?: {
    page?: number
    limit?: number
    category?: number
    search?: string
  }): Promise<any> => {
    const response = await apiClient.get('/content', { params })
    return response.data
  },

  getContentById: async (id: number): Promise<any> => {
    const response = await apiClient.get(`/content/${id}`)
    return response.data
  },

  getCategories: async (): Promise<any> => {
    const response = await apiClient.get('/categories')
    return response.data
  },

  getRecommendations: async (limit?: number): Promise<any> => {
    const response = await apiClient.get('/recommendations', { params: { limit } })
    return response.data
  },
}

export const adminApi = {
  // Users
  getUsers: async (params?: { page?: number; limit?: number }): Promise<any> => {
    const response = await apiClient.get('/admin/users', { params })
    return response.data
  },

  createUser: async (userData: any): Promise<any> => {
    const response = await apiClient.post('/admin/users', userData)
    return response.data
  },

  updateUser: async (id: number, userData: any): Promise<any> => {
    const response = await apiClient.put(`/admin/users/${id}`, userData)
    return response.data
  },

  deleteUser: async (id: number): Promise<any> => {
    const response = await apiClient.delete(`/admin/users/${id}`)
    return response.data
  },

  // Content
  getAdminContent: async (params?: { page?: number; limit?: number }): Promise<any> => {
    const response = await apiClient.get('/admin/content', { params })
    return response.data
  },

  createContent: async (contentData: any): Promise<any> => {
    const response = await apiClient.post('/admin/content', contentData)
    return response.data
  },

  updateContent: async (id: number, contentData: any): Promise<any> => {
    const response = await apiClient.put(`/admin/content/${id}`, contentData)
    return response.data
  },

  deleteContent: async (id: number): Promise<any> => {
    const response = await apiClient.delete(`/admin/content/${id}`)
    return response.data
  },

  // Categories
  getAdminCategories: async (): Promise<any> => {
    const response = await apiClient.get('/admin/categories')
    return response.data
  },

  createCategory: async (categoryData: any): Promise<any> => {
    const response = await apiClient.post('/admin/categories', categoryData)
    return response.data
  },

  updateCategory: async (id: number, categoryData: any): Promise<any> => {
    const response = await apiClient.put(`/admin/categories/${id}`, categoryData)
    return response.data
  },

  deleteCategory: async (id: number): Promise<any> => {
    const response = await apiClient.delete(`/admin/categories/${id}`)
    return response.data
  },
}
