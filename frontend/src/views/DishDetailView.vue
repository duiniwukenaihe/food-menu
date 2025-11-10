<template>
  <div class="grid gap-10 lg:grid-cols-[minmax(0,1fr),320px]">
    <article class="space-y-6">
      <div class="space-y-3">
        <RouterLink to="/" class="inline-flex items-center gap-2 text-sm text-[var(--color-muted)] transition hover:text-[var(--color-text)]">
          <span aria-hidden="true">←</span>
          {{ $t('dish.backToHome') }}
        </RouterLink>
        <p class="text-sm font-medium uppercase tracking-wider text-emerald-500">
          {{ dish?.region ?? $t('dish.unknownRegion') }}
        </p>
        <h1 class="text-3xl font-bold tracking-tight">{{ dish?.name ?? formattedId }}</h1>
        <p class="text-base text-[var(--color-muted)]">{{ dish?.description ?? $t('dish.placeholderDescription') }}</p>
      </div>

      <MediaDisplay v-if="dish" :src="dish.image" :alt="dish.name" aspect="16:9" />

      <BaseCard v-if="dish" :title="$t('dish.storyTitle')" variant="outlined">
        <p class="text-sm leading-relaxed text-[var(--color-muted)]">{{ dish.story }}</p>
      </BaseCard>

      <BaseCard :title="$t('dish.techniquesTitle')">
        <ul v-if="dish" class="grid list-disc gap-3 pl-5 text-sm text-[var(--color-muted)]">
          <li v-for="technique in dish.techniques" :key="technique">{{ technique }}</li>
        </ul>
        <p v-else class="text-sm text-[var(--color-muted)]">{{ $t('dish.techniquesFallback') }}</p>
      </BaseCard>
    </article>

    <aside class="space-y-6">
      <BaseCard :title="$t('dish.quickFactsTitle')" :subtitle="dish?.region">
        <dl v-if="dish" class="grid grid-cols-1 gap-3 text-sm">
          <div v-for="fact in dish.quickFacts" :key="fact.label" class="flex items-center justify-between gap-2">
            <dt class="text-[var(--color-muted)]">{{ fact.label }}</dt>
            <dd class="font-medium">{{ fact.value }}</dd>
          </div>
        </dl>
        <p v-else class="text-sm text-[var(--color-muted)]">{{ $t('dish.quickFactsFallback') }}</p>
      </BaseCard>

      <BaseCard :title="$t('dish.pairingsTitle')" variant="outlined">
        <ul v-if="dish" class="space-y-2 text-sm text-[var(--color-muted)]">
          <li v-for="pairing in dish.pairings" :key="pairing">{{ pairing }}</li>
        </ul>
        <p v-else class="text-sm text-[var(--color-muted)]">{{ $t('dish.pairingsFallback') }}</p>
      </BaseCard>

      <BaseCard :title="$t('dish.callToActionTitle')" :subtitle="$t('dish.callToActionSubtitle')" padded>
        <p class="text-sm text-[var(--color-muted)]">{{ $t('dish.callToActionBody') }}</p>
        <template #footer>
          <RouterLink to="/admin" class="inline-flex items-center gap-2 text-sm font-semibold text-emerald-500 transition hover:text-emerald-400">
            {{ $t('dish.callToActionCta') }}
            <span aria-hidden="true">→</span>
          </RouterLink>
        </template>
      </BaseCard>
    </aside>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'

const props = defineProps<{ id: string }>()

const catalog = {
  'artisan-pasta': {
    id: 'artisan-pasta',
    name: 'Artisan Saffron Pasta',
    region: 'Northern Italy',
    description: 'Hand-rolled pasta infused with saffron and finished with a Meyer lemon beurre blanc.',
    story:
      'Our pasta is crafted daily using heritage grains milled on-site. The saffron threads are sourced directly from a cooperative of growers in Abruzzo, steeped gently to preserve their floral aroma.',
    image: 'https://images.unsplash.com/photo-1525755662778-989d0524087e?auto=format&fit=crop&w=1200&q=80',
    techniques: ['Slow hydration', 'Bronze-die extrusion', 'Hand lamination', 'Table-side finishing'],
    quickFacts: [
      { label: 'Course', value: 'Primo' },
      { label: 'Cook time', value: '8 minutes' },
      { label: 'Spice level', value: 'Gentle warmth' },
    ],
    pairings: ['Verdicchio Superiore', 'Charred Meyer lemon', 'Fennel pollen dust'],
  },
  'fusion-sushi': {
    id: 'fusion-sushi',
    name: 'Amberjack Fusion Sushi',
    region: 'Kyoto, Japan',
    description: 'Sustainable amberjack sashimi layered with pickled cherry blossoms and foraged mountain herbs.',
    story:
      'Inspired by spring omakase traditions, this dish balances oceanic sweetness with the vibrant aromatics of sansho pepper. Each slice is brushed with barrel-aged soy for depth.',
    image: 'https://images.unsplash.com/photo-1553621042-f6e147245754?auto=format&fit=crop&w=1200&q=80',
    techniques: ['Ikejime preparation', 'Flash curing', 'Hand pressing', 'Floral smoke infusion'],
    quickFacts: [
      { label: 'Course', value: 'Chef signature' },
      { label: 'Serving temp', value: 'Chilled' },
      { label: 'Texture', value: 'Melt-in-mouth' },
    ],
    pairings: ['Junmai Daiginjo', 'Compressed cucumber', 'Granite of sudachi'],
  },
  'patagonia-lamb': {
    id: 'patagonia-lamb',
    name: 'Hearth-Roasted Patagonia Lamb',
    region: 'Patagonia, Argentina',
    description: 'Grass-fed lamb smoked over lenga wood and paired with charred chimichurri.',
    story:
      'The lamb is dry-aged for 14 days and smoked over carefully cured lenga wood. A charred chimichurri featuring indigenous merken peppers adds warmth and complexity.',
    image: 'https://images.unsplash.com/photo-1604908177522-4023ac76f015?auto=format&fit=crop&w=1200&q=80',
    techniques: ['Wood-fire roasting', 'Dry aging', 'Cold smoking', 'Table-side carving'],
    quickFacts: [
      { label: 'Course', value: 'Signature main' },
      { label: 'Cook time', value: '90 minutes' },
      { label: 'Texture', value: 'Tender & smoky' },
    ],
    pairings: ['Malbec Reserva', 'Fire-charred root vegetables', 'Smoked salt flakes'],
  },
} as const

type Dish = (typeof catalog)[keyof typeof catalog]

const dish = computed<Dish | null>(() => catalog[props.id as keyof typeof catalog] ?? null)

const formattedId = computed(() => props.id.replace(/[-_]/g, ' '))
</script>
