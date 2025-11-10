<template>
  <section class="space-y-12">
    <div class="grid gap-10 md:grid-cols-[1.1fr,0.9fr] md:items-center">
      <div class="space-y-6">
        <p class="text-sm font-medium uppercase tracking-widest text-emerald-500">
          {{ $t('home.heroEyebrow') }}
        </p>
        <h1 class="text-3xl font-bold tracking-tight sm:text-4xl lg:text-5xl">
          {{ $t('home.heroTitle') }}
        </h1>
        <p class="max-w-xl text-lg text-[var(--color-muted)]">
          {{ $t('home.heroSubtitle') }}
        </p>
        <div class="flex flex-wrap gap-4">
          <RouterLink
            to="/dish/artisan-pasta"
            class="inline-flex items-center gap-2 rounded-full bg-emerald-500 px-6 py-3 text-sm font-semibold text-white shadow-md transition hover:bg-emerald-600"
          >
            {{ $t('home.ctaDiscover') }}
            <span aria-hidden="true">→</span>
          </RouterLink>
          <RouterLink
            to="/admin"
            class="inline-flex items-center gap-2 rounded-full border border-[var(--border-color)] px-6 py-3 text-sm font-semibold text-[var(--color-text)] transition hover:border-transparent hover:bg-slate-900 hover:text-white"
          >
            {{ $t('home.ctaDashboard') }}
          </RouterLink>
        </div>
      </div>
      <div class="relative">
        <MediaDisplay
          src="https://images.unsplash.com/photo-1504674900247-0877df9cc836?auto=format&fit=crop&w=1200&q=80"
          :alt="$t('home.heroImageAlt')"
          aspect="3:2"
        >
          <template #caption>
            {{ $t('home.heroImageCaption') }}
          </template>
        </MediaDisplay>
      </div>
    </div>

    <div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
      <BaseCard
        v-for="feature in features"
        :key="feature.title"
        :title="$t(feature.title)"
        :subtitle="$t(feature.subtitle)"
        variant="outlined"
      >
        <p class="text-sm text-[var(--color-muted)]">{{ $t(feature.description) }}</p>
      </BaseCard>
    </div>

    <section class="space-y-6">
      <div class="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h2 class="text-2xl font-semibold tracking-tight">{{ $t('home.curatedHeading') }}</h2>
          <p class="text-sm text-[var(--color-muted)]">{{ $t('home.curatedSubtitle') }}</p>
        </div>
        <RouterLink
          to="/dish/artisan-pasta"
          class="inline-flex items-center gap-2 text-sm font-semibold text-emerald-500 transition hover:text-emerald-400"
        >
          {{ $t('home.curatedCta') }}
          <span aria-hidden="true">→</span>
        </RouterLink>
      </div>

      <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <BaseCard
          v-for="dish in curatedDishes"
          :key="dish.id"
          :title="dish.name"
          :subtitle="dish.region"
        >
          <MediaDisplay :src="dish.image" :alt="dish.name" aspect="4:3" />
          <p class="text-sm text-[var(--color-muted)]">{{ dish.description }}</p>
          <template #footer>
            <RouterLink :to="`/dish/${dish.id}`" class="text-sm font-semibold text-emerald-500 transition hover:text-emerald-400">
              {{ $t('home.viewDish') }}
            </RouterLink>
          </template>
        </BaseCard>
      </div>
    </section>
  </section>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'

const features = [
  {
    title: 'home.features.seasonal.title',
    subtitle: 'home.features.seasonal.subtitle',
    description: 'home.features.seasonal.description',
  },
  {
    title: 'home.features.sourcing.title',
    subtitle: 'home.features.sourcing.subtitle',
    description: 'home.features.sourcing.description',
  },
  {
    title: 'home.features.community.title',
    subtitle: 'home.features.community.subtitle',
    description: 'home.features.community.description',
  },
]

const curatedDishes = [
  {
    id: 'artisan-pasta',
    name: 'Artisan Saffron Pasta',
    region: 'Northern Italy',
    image: 'https://images.unsplash.com/photo-1525755662778-989d0524087e?auto=format&fit=crop&w=1200&q=80',
    description: 'Hand-rolled pasta infused with saffron and served with a lemon beurre blanc.',
  },
  {
    id: 'fusion-sushi',
    name: 'Amberjack Fusion Sushi',
    region: 'Kyoto, Japan',
    image: 'https://images.unsplash.com/photo-1553621042-f6e147245754?auto=format&fit=crop&w=1200&q=80',
    description: 'Sustainable amberjack sashimi paired with foraged herbs and yuzu kosho.',
  },
  {
    id: 'patagonia-lamb',
    name: 'Hearth-Roasted Patagonia Lamb',
    region: 'Patagonia, Argentina',
    image: 'https://images.unsplash.com/photo-1525755662778-989d0524087e?auto=format&fit=crop&w=1200&q=80',
    description: 'Slow-roasted lamb with charred chimichurri and native root vegetables.',
  },
]
</script>
