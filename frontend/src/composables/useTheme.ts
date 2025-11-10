import { computed, onMounted, ref, watch } from 'vue'

type Theme = 'light' | 'dark'

const STORAGE_KEY = 'gourmet-explorer:theme'
const theme = ref<Theme>('light')
const isInitialized = ref(false)

const applyTheme = (value: Theme) => {
  if (typeof document === 'undefined') return
  const root = document.documentElement
  root.classList.toggle('dark', value === 'dark')
  root.classList.toggle('light', value === 'light')
  root.dataset.theme = value
  root.classList.add(value)
  root.classList.remove(value === 'dark' ? 'light' : 'dark')
}

watch(theme, value => {
  applyTheme(value)
  if (typeof window !== 'undefined') {
    window.localStorage.setItem(STORAGE_KEY, value)
  }
})

const initialize = () => {
  if (isInitialized.value) return
  if (typeof window === 'undefined') return

  const storedTheme = window.localStorage.getItem(STORAGE_KEY) as Theme | null
  const prefersDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches

  theme.value = storedTheme ?? (prefersDark ? 'dark' : 'light')
  applyTheme(theme.value)
  isInitialized.value = true
}

const setTheme = (value: Theme) => {
  theme.value = value
}

const toggleTheme = () => {
  theme.value = theme.value === 'dark' ? 'light' : 'dark'
}

export const useTheme = () => {
  onMounted(() => {
    initialize()
  })

  return {
    theme: computed(() => theme.value),
    isDark: computed(() => theme.value === 'dark'),
    setTheme,
    toggleTheme,
    initialize,
  }
}
