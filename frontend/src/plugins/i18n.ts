import { createI18n } from 'vue-i18n'

const messages = {
  en: {
    navigation: {
      menu: 'Menu',
      close: 'Close',
      home: 'Home',
      featuredDish: 'Featured dish',
      admin: 'Admin',
      switchToLight: 'Switch to light mode',
      switchToDark: 'Switch to dark mode',
    },
    auth: {
      startDemo: 'Start demo session',
      signOut: 'Sign out',
    },
    footer: {
      rightsReserved: 'All rights reserved.',
      contact: 'Contact',
      privacy: 'Privacy',
      terms: 'Terms',
    },
    home: {
      heroEyebrow: 'Seasonal tasting collective',
      heroTitle: 'Celebrate storied dishes crafted by chefs around the globe',
      heroSubtitle:
        'Discover a curated collection of modern cuisine, connect with chefs, and orchestrate elevated dining experiences for your guests.',
      heroImageAlt: 'Chefs preparing a contemporary tasting dish',
      heroImageCaption: 'Seasonal tasting flight prepared by Chef Aurora',
      ctaDiscover: 'Discover dishes',
      ctaDashboard: 'Open dashboard',
      curatedHeading: 'Curated for this season',
      curatedSubtitle: 'Handpicked signatures that champion sustainable sourcing and thoughtful storytelling.',
      curatedCta: 'View the feature story',
      viewDish: 'Explore this dish',
      features: {
        seasonal: {
          title: 'Seasonality perfected',
          subtitle: 'Micro-seasonal menus',
          description: 'Curate weekly tastings that respond to micro-seasons and the biodiversity of each region.',
        },
        sourcing: {
          title: 'Transparent sourcing',
          subtitle: 'Farmer partnerships',
          description: 'Track provenance from coast to canopy with supply transparency for every ingredient used.',
        },
        community: {
          title: 'Chef collaboration',
          subtitle: 'Global brigade',
          description: 'Coordinate across your culinary collective with live planning tools and shared mise en place notes.',
        },
      },
    },
    dish: {
      backToHome: 'Back to home',
      unknownRegion: 'Origin unknown',
      placeholderDescription: 'We are gathering the full story behind this signature dish. Check back soon for poetic tasting notes.',
      storyTitle: 'The story behind the dish',
      techniquesTitle: 'Key techniques',
      techniquesFallback: 'Techniques will appear once this dish is curated.',
      quickFactsTitle: 'Quick facts',
      quickFactsFallback: 'We are compiling quick facts for this dish.',
      pairingsTitle: 'Pairings & accents',
      pairingsFallback: 'Pairing suggestions will be added soon.',
      callToActionTitle: 'Ready to feature this experience?',
      callToActionSubtitle: 'Unlock chef collaboration tools',
      callToActionBody: 'Access mise en place plans, sourcing partners, and plating videos by joining the admin workspace.',
      callToActionCta: 'Launch the admin workspace',
    },
    admin: {
      eyebrow: 'Operations cockpit',
      title: 'Admin control center',
      subtitle: 'Monitor dining experiences, coordinate chef contributions, and keep your guests delighted across every seating.',
      authRequiredTitle: 'Sign in to access the control center',
      authRequiredBody: 'Start a demo session to explore how the dashboard tracks experiences, reservations, and culinary collaborations.',
      authRequiredCta: 'Start demo session',
      welcomeTitle: 'Welcome back, {name}',
      welcomeBody: 'You have dashboard access with elevated privileges. Review the highlights below to stay ahead of tonight’s service.',
      metrics: {
        activeDishes: 'Live dishes',
        activeDishesDescription: 'Signature experiences currently available for booking.',
        pendingReviews: 'Reviews awaiting feedback',
        pendingReviewsDescription: 'Guest reflections needing a thoughtful reply.',
        tableHoldRequests: 'Table hold requests',
        tableHoldRequestsDescription: 'Pre-service requests awaiting confirmation.',
      },
      activityTitle: 'Latest activity',
    },
    media: {
      unavailable: 'Media unavailable',
    },
    notFound: {
      title: 'We could not find that page',
      body: 'The route you visited is not part of our tasting journey yet. Try returning to the home view or explore the curated dishes.',
      cta: 'Return home',
    },
  },
}

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages,
})

export default i18n
