<template>
  <header
    class="sticky top-0 z-50 border-b border-[var(--border-color)] bg-[var(--card-bg)]/95 backdrop-blur supports-[backdrop-filter]:bg-[var(--card-bg)]/80"
  >
    <nav class="mx-auto flex w-full max-w-6xl items-center justify-between px-4 py-4 sm:px-6 lg:px-8">
      <RouterLink to="/" class="text-lg font-semibold tracking-tight">
        Gourmet Explorer
      </RouterLink>

      <button class="inline-flex items-center gap-2 rounded-md border border-transparent px-3 py-2 text-sm font-medium shadow-sm transition hover:border-[var(--border-color)] sm:hidden" @click="isMenuOpen = !isMenuOpen">
        <span>{{ isMenuOpen ? $t('navigation.close') : $t('navigation.menu') }}</span>
        <span aria-hidden="true">☰</span>
      </button>

      <div
        class="flex-1 items-center justify-end gap-8 sm:flex"
        :class="isMenuOpen ? 'mt-4 flex flex-col gap-6 sm:mt-0 sm:flex-row sm:items-center sm:justify-end' : 'hidden sm:flex'"
      >
        <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:gap-6">
          <RouterLink v-for="link in navigation" :key="link.to" :to="link.to" class="text-sm font-medium text-[var(--color-muted)] transition hover:text-[var(--color-text)]" active-class="text-[var(--color-text)]">
            {{ $t(link.label) }}
          </RouterLink>
        </div>

        <div class="flex items-center gap-3">
          <button
            class="inline-flex items-center justify-center rounded-full border border-[var(--border-color)] bg-transparent p-2 text-base transition hover:ring-2 hover:ring-[var(--ring-color)]"
            type="button"
            @click="toggleTheme"
            :aria-label="$t(isDark ? 'navigation.switchToLight' : 'navigation.switchToDark')"
          >
            <span aria-hidden="true">{{ isDark ? '🌙' : '☀️' }}</span>
          </button>

          <div v-if="auth.user" class="flex items-center gap-2 text-sm">
            <span class="hidden text-[var(--color-muted)] sm:inline">
              {{ auth.user.name }}
            </span>
            <button
              class="rounded-md border border-transparent bg-rose-500 px-3 py-2 text-xs font-semibold text-white transition hover:bg-rose-600"
              type="button"
              @click="auth.logout()"
            >
              {{ $t('auth.signOut') }}
            </button>
          </div>
          <button
            v-else
            class="rounded-md border border-[var(--border-color)] bg-transparent px-3 py-2 text-xs font-semibold transition hover:border-transparent hover:bg-emerald-500 hover:text-white"
            type="button"
            @click="auth.startDemoSession()"
          >
            {{ $t('auth.startDemo') }}
          </button>
        </div>
      </div>
    </nav>
  </header>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'

import { useTheme } from '../../composables/useTheme'
import { useAuthStore } from '../../stores/auth'

const auth = useAuthStore()
const navigation = computed(() => [
  { to: '/', label: 'navigation.home' },
  { to: '/dish/artisan-pasta', label: 'navigation.featuredDish' },
  { to: '/admin', label: 'navigation.admin' },
])

const { toggleTheme, isDark, initialize } = useTheme()
const isMenuOpen = ref(false)
const route = useRoute()

onMounted(() => {
  initialize()
})

watch(
  () => route.fullPath,
  () => {
    isMenuOpen.value = false
  }
)
</script>
