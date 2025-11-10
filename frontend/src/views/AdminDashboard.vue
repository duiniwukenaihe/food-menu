<template>
  <section class="space-y-8">
    <header class="space-y-3">
      <p class="text-sm font-medium uppercase tracking-wider text-emerald-500">{{ $t('admin.eyebrow') }}</p>
      <h1 class="text-3xl font-bold tracking-tight">{{ $t('admin.title') }}</h1>
      <p class="max-w-2xl text-base text-[var(--color-muted)]">{{ $t('admin.subtitle') }}</p>
    </header>

    <BaseCard v-if="!isAuthenticated" :title="$t('admin.authRequiredTitle')" variant="outlined">
      <p class="text-sm text-[var(--color-muted)]">{{ $t('admin.authRequiredBody') }}</p>
      <template #footer>
        <button
          type="button"
          class="inline-flex items-center gap-2 rounded-full bg-emerald-500 px-5 py-2 text-sm font-semibold text-white transition hover:bg-emerald-600"
          @click="startDemoSession"
        >
          {{ $t('admin.authRequiredCta') }}
        </button>
      </template>
    </BaseCard>

    <div v-else class="space-y-6">
      <BaseCard :title="$t('admin.welcomeTitle', { name: user?.name })" :subtitle="user?.email">
        <p class="text-sm text-[var(--color-muted)]">
          {{ $t('admin.welcomeBody') }}
        </p>
      </BaseCard>

      <div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        <BaseCard
          v-for="metric in metrics"
          :key="metric.label"
          :title="$t(metric.label)"
          variant="outlined"
        >
          <p class="text-3xl font-semibold">{{ metric.value }}</p>
          <p class="text-sm text-[var(--color-muted)]">{{ $t(metric.description) }}</p>
        </BaseCard>
      </div>

      <BaseCard :title="$t('admin.activityTitle')" variant="outlined">
        <ul class="space-y-4 text-sm text-[var(--color-muted)]">
          <li v-for="activity in activityFeed" :key="activity.timestamp" class="flex items-start justify-between gap-3">
            <div>
              <p class="font-medium text-[var(--color-text)]">{{ activity.message }}</p>
              <p class="text-xs">{{ activity.category }}</p>
            </div>
            <span class="text-xs text-[var(--color-muted)]">{{ activity.timestamp }}</span>
          </li>
        </ul>
      </BaseCard>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'

import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const { user, isAuthenticated } = storeToRefs(auth)

const startDemoSession = () => {
  auth.startDemoSession()
}

const metrics = computed(() => [
  { label: 'admin.metrics.activeDishes', value: 42, description: 'admin.metrics.activeDishesDescription' },
  { label: 'admin.metrics.pendingReviews', value: 8, description: 'admin.metrics.pendingReviewsDescription' },
  { label: 'admin.metrics.tableHoldRequests', value: 14, description: 'admin.metrics.tableHoldRequestsDescription' },
])

const activityFeed = computed(() => [
  { message: 'Chef Aurora published a new seasonal tasting menu', category: 'Menus', timestamp: '3 minutes ago' },
  { message: 'Seven new reservation hold requests pending approval', category: 'Reservations', timestamp: '22 minutes ago' },
  { message: 'Guest feedback score increased by 6%', category: 'Insights', timestamp: '1 hour ago' },
])
</script>
