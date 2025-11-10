import axios, { type AxiosInstance, type AxiosRequestConfig } from 'axios'

const baseURL = import.meta.env.VITE_API_BASE_URL ?? '/api'

export const apiClient: AxiosInstance = axios.create({
  baseURL,
  withCredentials: true,
  timeout: 15000,
})

export type AuthBridge = {
  getAccessToken: () => string | null
  handleUnauthorized: () => void
}

let authBridge: AuthBridge | null = null

export const configureApiClient = (bridge: AuthBridge) => {
  authBridge = bridge
}

apiClient.interceptors.request.use(config => {
  const token = authBridge?.getAccessToken()
  if (token) {
    config.headers = {
      ...(config.headers ?? {}),
      Authorization: `Bearer ${token}`,
    }
  }
  return config
})

apiClient.interceptors.response.use(
  response => response,
  error => {
    const status = error.response?.status
    if (status === 401) {
      authBridge?.handleUnauthorized()
    }
    return Promise.reject(error)
  }
)

export const createApiClient = (config?: AxiosRequestConfig) => {
  const instance = axios.create({ baseURL, timeout: 15000, ...config })
  if (authBridge) {
    instance.interceptors.request.use(requestConfig => {
      const token = authBridge?.getAccessToken()
      if (token) {
        requestConfig.headers = {
          ...(requestConfig.headers ?? {}),
          Authorization: `Bearer ${token}`,
        }
      }
      return requestConfig
    })
    instance.interceptors.response.use(
      response => response,
      error => {
        if (error.response?.status === 401) {
          authBridge?.handleUnauthorized()
        }
        return Promise.reject(error)
      }
    )
  }
  return instance
}
