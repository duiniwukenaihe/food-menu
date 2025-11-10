<template>
  <section class="flex flex-col items-center justify-center gap-6 py-24 text-center">
    <BaseCard :title="$t('notFound.title')" variant="outlined">
      <p class="text-sm text-[var(--color-muted)]">{{ $t('notFound.body') }}</p>
      <template #footer>
        <RouterLink
          to="/"
          class="inline-flex items-center gap-2 rounded-full bg-emerald-500 px-6 py-2 text-sm font-semibold text-white transition hover:bg-emerald-600"
        >
          {{ $t('notFound.cta') }}
        </RouterLink>
      </template>
    </BaseCard>
  </section>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>
