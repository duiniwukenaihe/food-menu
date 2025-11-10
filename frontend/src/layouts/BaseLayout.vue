<template>
  <div class="min-h-screen bg-[var(--color-bg)] text-[var(--color-text)] transition-colors flex flex-col">
    <AppHeader />
    <main class="flex-1">
      <div class="mx-auto w-full max-w-6xl px-4 py-6 sm:px-6 lg:px-8">
        <slot />
      </div>
    </main>
    <AppFooter />
  </div>
</template>

<script setup lang="ts">
import AppHeader from '../components/layout/AppHeader.vue'
import AppFooter from '../components/layout/AppFooter.vue'
</script>
