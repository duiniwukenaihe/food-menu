<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="modelValue" class="modal-open fixed inset-0 z-[100] flex items-center justify-center">
        <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="handleBackdrop"></div>
        <section
          class="relative z-10 w-full max-w-lg overflow-hidden rounded-2xl border border-[var(--border-color)] bg-[var(--card-bg)] shadow-xl"
          role="dialog"
          aria-modal="true"
        >
          <header class="flex items-center justify-between border-b border-[var(--border-color)] px-6 py-4">
            <h2 v-if="title" class="text-lg font-semibold">{{ title }}</h2>
            <button type="button" class="rounded-full p-2 text-sm text-[var(--color-muted)] transition hover:bg-slate-100/50 dark:hover:bg-slate-700/50" @click="close">
              <span aria-hidden="true">×</span>
            </button>
          </header>
          <div class="max-h-[70vh] overflow-y-auto px-6 py-4">
            <slot />
          </div>
          <footer v-if="$slots.footer" class="border-t border-[var(--border-color)] bg-[var(--card-bg)]/80 px-6 py-4">
            <slot name="footer" />
          </footer>
        </section>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, watch } from 'vue'

type Props = {
  modelValue: boolean
  title?: string
  closeOnBackdrop?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  closeOnBackdrop: true,
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  close: []
}>

const lockScroll = (shouldLock: boolean) => {
  if (typeof document === 'undefined') return
  document.body.classList.toggle('modal-open', shouldLock)
}

const close = () => {
  emit('update:modelValue', false)
  emit('close')
}

const handleBackdrop = () => {
  if (props.closeOnBackdrop) {
    close()
  }
}

watch(
  () => props.modelValue,
  value => {
    lockScroll(value)
  },
  { immediate: true }
)

onMounted(() => {
  if (props.modelValue) {
    lockScroll(true)
  }
})

onBeforeUnmount(() => {
  lockScroll(false)
})
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
</style>
