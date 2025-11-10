<template>
  <footer class="border-t border-[var(--border-color)] bg-[var(--card-bg)]/80">
    <div class="mx-auto flex w-full max-w-6xl flex-col gap-4 px-4 py-6 text-sm text-[var(--color-muted)] sm:flex-row sm:items-center sm:justify-between sm:px-6 lg:px-8">
      <p>&copy; {{ new Date().getFullYear() }} Gourmet Explorer. {{ $t('footer.rightsReserved') }}</p>
      <div class="flex items-center gap-4">
        <RouterLink to="/" class="transition hover:text-[var(--color-text)]">
          {{ $t('footer.contact') }}
        </RouterLink>
        <RouterLink to="/" class="transition hover:text-[var(--color-text)]">
          {{ $t('footer.privacy') }}
        </RouterLink>
        <RouterLink to="/" class="transition hover:text-[var(--color-text)]">
          {{ $t('footer.terms') }}
        </RouterLink>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>
