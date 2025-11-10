<template>
  <figure :class="wrapperClass">
    <component
      :is="componentTag"
      v-if="!hasError"
      v-bind="mediaAttributes"
      :src="src"
      loading="lazy"
      class="h-full w-full object-cover"
      @error="handleError"
    >
      <track kind="captions" v-if="type === 'video'" />
    </component>

    <div v-else class="flex h-full flex-col items-center justify-center gap-2 text-sm text-[var(--color-muted)]">
      <slot name="fallback">
        <span aria-hidden="true">📷</span>
        <span>{{ $t('media.unavailable') }}</span>
      </slot>
    </div>

    <figcaption v-if="$slots.caption" class="mt-3 text-sm text-[var(--color-muted)]">
      <slot name="caption" />
    </figcaption>
  </figure>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

const props = withDefaults(
  defineProps<{
    src: string
    alt?: string
    type?: 'image' | 'video'
    aspect?: '16:9' | '4:3' | '1:1' | '3:2'
    poster?: string
  }>(),
  {
    type: 'image',
    aspect: '16:9',
  }
)

const hasError = ref(false)

const mediaAttributes = computed(() => {
  if (props.type === 'video') {
    return {
      controls: true,
      playsinline: true,
      preload: 'metadata' as const,
      poster: props.poster,
    }
  }

  return {
    alt: props.alt ?? '',
  }
})

const aspectClass = computed(() => {
  switch (props.aspect) {
    case '4:3':
      return 'aspect-[4/3]'
    case '1:1':
      return 'aspect-square'
    case '3:2':
      return 'aspect-[3/2]'
    default:
      return 'aspect-video'
  }
})

const wrapperClass = computed(() => ['relative overflow-hidden rounded-xl border border-[var(--border-color)] bg-[var(--card-bg)]', aspectClass.value].join(' '))

const componentTag = computed(() => (props.type === 'video' ? 'video' : 'img'))

const handleError = () => {
  hasError.value = true
}
</script>
