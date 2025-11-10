<template>
  <section :class="wrapperClasses">
    <header v-if="hasHeader" class="flex items-start justify-between gap-3">
      <div class="flex flex-1 items-start gap-3">
        <div v-if="$slots.icon" class="mt-1 text-xl">
          <slot name="icon" />
        </div>
        <div class="flex-1">
          <h3 v-if="title" class="text-base font-semibold">{{ title }}</h3>
          <p v-if="subtitle" class="mt-1 text-sm text-[var(--color-muted)]">{{ subtitle }}</p>
        </div>
      </div>
      <div v-if="$slots.actions" class="flex items-center gap-2">
        <slot name="actions" />
      </div>
    </header>

    <div :class="bodyClasses">
      <slot />
    </div>

    <footer v-if="$slots.footer" class="mt-4 border-t border-[var(--border-color)] pt-4">
      <slot name="footer" />
    </footer>
  </section>
</template>

<script setup lang="ts">
import { computed, useSlots } from 'vue'

type CardVariant = 'elevated' | 'outlined' | 'ghost'

type Props = {
  title?: string
  subtitle?: string
  padded?: boolean
  variant?: CardVariant
}

const props = withDefaults(defineProps<Props>(), {
  padded: true,
  variant: 'elevated',
})

const slots = useSlots()

const hasHeader = computed(() => Boolean(props.title || props.subtitle || slots.icon || slots.actions))

const wrapperClasses = computed(() => {
  const base = [
    'rounded-xl',
    'border',
    'border-[var(--border-color)]',
    'bg-[var(--card-bg)]/95',
    'shadow-sm',
    'transition-shadow',
    'hover:shadow-md',
  ]

  if (props.variant === 'outlined') {
    base.push('bg-transparent')
  }

  if (props.variant === 'ghost') {
    base.push('border-transparent', 'bg-transparent', 'shadow-none', 'hover:shadow-none')
  }

  return base.join(' ')
})

const bodyClasses = computed(() => {
  if (!props.padded) {
    return hasHeader.value ? 'mt-2' : ''
  }

  const classes = ['space-y-4', 'px-4']
  if (hasHeader.value) {
    classes.push('mt-4', 'pb-4')
  } else {
    classes.push('py-4')
  }
  return classes.join(' ')
})
</script>
