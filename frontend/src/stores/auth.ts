import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

import type { User } from '../types/user'

const SESSION_KEY = 'gourmet-explorer:session'

export type SessionPayload = {
  token: string
  user: User
}

const readSession = (): SessionPayload | null => {
  if (typeof window === 'undefined') return null
  const stored = window.localStorage.getItem(SESSION_KEY)
  if (!stored) return null
  try {
    return JSON.parse(stored) as SessionPayload
  } catch (error) {
    console.warn('[auth] Failed to parse stored session', error)
    return null
  }
}

const writeSession = (session: SessionPayload | null) => {
  if (typeof window === 'undefined') return
  if (!session) {
    window.localStorage.removeItem(SESSION_KEY)
    return
  }
  window.localStorage.setItem(SESSION_KEY, JSON.stringify(session))
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const isLoading = ref(true)

  const isAuthenticated = computed(() => Boolean(token.value))
  const role = computed(() => user.value?.role ?? 'guest')

  const setSession = (session: SessionPayload) => {
    user.value = session.user
    token.value = session.token
    writeSession(session)
  }

  const clearSession = () => {
    user.value = null
    token.value = null
    writeSession(null)
  }

  const initialize = async () => {
    isLoading.value = true
    try {
      const session = readSession()
      if (session) {
        user.value = session.user
        token.value = session.token
      }
    } finally {
      isLoading.value = false
    }
  }

  const logout = () => {
    clearSession()
  }

  const handleUnauthorized = () => {
    clearSession()
  }

  const startDemoSession = () => {
    setSession({
      token: 'demo-admin-token',
      user: {
        id: 'demo-admin',
        name: 'Demo Admin',
        email: 'demo@gourmet.app',
        role: 'admin',
      },
    })
  }

  const getAccessToken = () => token.value

  return {
    user,
    token,
    role,
    isLoading,
    isAuthenticated,
    initialize,
    setSession,
    logout,
    handleUnauthorized,
    startDemoSession,
    getAccessToken,
  }
})
