import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import i18n from './plugins/i18n'
import './index.css'

import BaseCard from './components/base/BaseCard.vue'
import BaseModal from './components/base/BaseModal.vue'
import MediaDisplay from './components/base/MediaDisplay.vue'
import { configureApiClient } from './api/client'
import { useAuthStore } from './stores/auth'

async function bootstrap() {
  const app = createApp(App)

  const pinia = createPinia()
  app.use(pinia)

  const authStore = useAuthStore(pinia)
  configureApiClient({
    getAccessToken: authStore.getAccessToken,
    handleUnauthorized: authStore.handleUnauthorized,
  })
  await authStore.initialize()

  app.use(router)
  app.use(i18n)

  app.component('BaseCard', BaseCard)
  app.component('BaseModal', BaseModal)
  app.component('MediaDisplay', MediaDisplay)

  app.mount('#app')
}

bootstrap()
