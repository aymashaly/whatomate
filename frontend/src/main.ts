import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { VueQueryPlugin } from '@tanstack/vue-query'

import App from './App.vue'
import router from './router'
import { i18n } from './i18n'
import { useBrandingStore } from './stores/branding'

import './assets/fonts.css'
import './assets/index.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(VueQueryPlugin)
app.use(i18n)

// Apply platform-wide branding (name, favicon, accent color) before the first
// paint so the user never sees a hardcoded "Whatomate" / default favicon.
// Public endpoint, so no auth needed; failures fall back to defaults silently.
useBrandingStore().load()

app.mount('#app')
