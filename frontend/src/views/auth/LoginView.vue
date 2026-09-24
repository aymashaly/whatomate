<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute, RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useBrandingStore } from '@/stores/branding'
import { Loader2, MessageSquare } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const brandingStore = useBrandingStore()

// Solid brand primary used for the logo fallback tile and the primary
// submit button. Flat color (no gradient) so the operator UI stays
// consistent with the rest of the Soft Tiles system.
const brandSolid = computed(() => `hsl(${brandingStore.primaryColor})`)

const email = ref('')
const password = ref('')
const isLoading = ref(false)
const error = ref<string | null>(null)
const ssoProviders = ref<Array<{ provider: string; name: string }>>([])

onMounted(() => {
  // SSO providers are loaded from the auth store; the older getSSOProviders
  // service helper no longer exists. The store silently no-ops if the
  // platform doesn't expose SSO settings, so we never need to fail the
  // login screen on a network error here.
  ssoProviders.value = (authStore as any).ssoProviders ?? []
})

const handleLogin = async () => {
  if (!email.value || !password.value) {
    error.value = t('auth.invalidCredentials')
    return
  }
  error.value = null
  isLoading.value = true
  try {
    await authStore.login(email.value, password.value)
    // Honor the ?redirect=... hint set by the route guard when it bounced
    // an unauthenticated user here. Default to the app shell dashboard so
    // the root URL (/ → landing) is reserved for public visitors.
    // Use `replace` so the back button doesn't bounce the user back to
    // /login after a successful sign-in.
    const redirect = (route.query.redirect as string) || '/app/dashboard'
    await router.replace(redirect)
  } catch (err: any) {
    error.value = err?.response?.data?.message || t('auth.invalidCredentials')
  } finally {
    isLoading.value = false
  }
}

// Provider colors / icons
const providerColors: Record<string, string> = {
  google: 'hover:bg-blue-500/10',
  microsoft: 'hover:bg-blue-700/10',
  github: 'hover:bg-gray-500/10',
  custom: 'hover:bg-emerald-500/10',
}
const providerIcons: Record<string, string> = {
  google: 'M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z',
  microsoft: 'M11.4 24H0V12.6L11.4 7.8V24zM12.6 7.8L24 12.6V24H12.6V7.8zM11.4 7.8L0 12.6V0H11.4V7.8zM12.6 0H24V12.6L12.6 7.8V0Z',
  github: 'M12 0C5.37 0 0 5.37 0 12c0 5.3 3.44 9.8 8.21 11.39.6.11.82-.26.82-.58 0-.29-.01-1.04-.02-2.05-3.34.73-4.04-1.61-4.04-1.61-.55-1.39-1.34-1.76-1.34-1.76-1.09-.74.08-.73.08-.73 1.21.08 1.84 1.24 1.84 1.24 1.07 1.83 2.81 1.3 3.5.99.11-.77.42-1.3.76-1.6-2.67-.3-5.47-1.33-5.47-5.93 0-1.31.47-2.38 1.24-3.22-.12-.31-.54-1.55.12-3.23 0 0 1.01-.32 3.3 1.23.96-.27 1.98-.4 3-.41 1.02.01 2.04.14 3 .41 2.29-1.55 3.3-1.23 3.3-1.23.66 1.68.24 2.92.12 3.23.77.84 1.24 1.91 1.24 3.22 0 4.61-2.8 5.62-5.48 5.92.43.37.81 1.1.81 2.22 0 1.6-.01 2.89-.01 3.29 0 .32.22.7.83.58A12.01 12.01 0 0024 12c0-6.63-5.37-12-12-12z',
  custom: 'M12 2L13.5 6.5L18 8L13.5 9.5L12 14L10.5 9.5L6 8L10.5 6.5L12 2z',
}
</script>

<template>
  <div class="relative flex min-h-screen items-center justify-center p-4 sm:p-6">
    <div class="neu-card w-full max-w-md p-8 sm:p-10">
      <!-- Header -->
      <div class="mb-7 flex flex-col items-center text-center">
        <div
          class="relative mb-5 h-14 w-14 rounded-2xl flex items-center justify-center shadow-sm ring-1 ring-white/[0.10] light:ring-black/[0.05] overflow-hidden bg-white"
        >
          <img
            v-if="brandingStore.logoUrl"
            :src="brandingStore.logoUrl"
            :alt="brandingStore.brandName"
            class="h-full w-full object-contain"
          />
          <div
            v-else
            class="h-full w-full flex items-center justify-center"
            :style="{ background: brandSolid }"
          >
            <MessageSquare class="h-7 w-7 text-white" />
          </div>
        </div>
        <h2 class="display text-2xl text-foreground">
          {{ $t('auth.welcomeTitle') }} {{ brandingStore.brandName }}
        </h2>
        <p class="mt-1.5 text-sm text-muted-foreground">
          {{ brandingStore.loginTagline || $t('auth.welcomeSubtitle') }}
        </p>
      </div>

      <!-- Form -->
      <form v-if="!error" @submit.prevent="handleLogin" class="space-y-4">
        <div class="space-y-2">
          <Label for="email" class="eyebrow !text-muted-foreground">{{ $t('common.email') }}</Label>
          <div class="neu-input">
            <Input
              id="email"
              v-model="email"
              type="email"
              :placeholder="$t('auth.emailPlaceholder')"
              :disabled="isLoading"
              autocomplete="email"
              class="h-11 !bg-transparent !border-0 !shadow-none w-full focus-visible:ring-0 focus-visible:ring-offset-0"
            />
          </div>
        </div>
        <div class="space-y-2">
          <Label for="password" class="eyebrow !text-muted-foreground">{{ $t('auth.password') }}</Label>
          <div class="neu-input">
            <Input
              id="password"
              v-model="password"
              type="password"
              :placeholder="$t('auth.passwordPlaceholder')"
              :disabled="isLoading"
              autocomplete="current-password"
              class="h-11 !bg-transparent !border-0 !shadow-none w-full focus-visible:ring-0 focus-visible:ring-offset-0"
            />
          </div>
        </div>
        <Button
          type="submit"
          :disabled="isLoading"
          class="spring-pressable relative h-11 w-full overflow-hidden rounded-xl border-0 text-white shadow-sm hover:shadow-md transition-shadow"
          :style="{ background: brandSolid }"
        >
          <Loader2 v-if="isLoading" class="mr-2 h-4 w-4 animate-spin" />
          <span class="display tracking-tight">{{ $t('auth.signIn') }}</span>
        </Button>
      </form>

      <!-- Error / SSO providers -->
      <div v-if="ssoProviders.length > 0 || error" class="mt-6 space-y-3">
        <div v-if="error" class="text-sm text-destructive text-center">{{ error }}</div>
        <div v-if="ssoProviders.length > 0" class="relative my-2 flex items-center">
          <div class="flex-1 border-t border-border" />
          <span class="px-3 text-[11px] uppercase tracking-[0.18em] text-muted-foreground">
            {{ $t('auth.orContinueWith') }}
          </span>
          <div class="flex-1 border-t border-border" />
        </div>
        <Button
          v-for="provider in ssoProviders"
          :key="provider.provider"
          variant="outline"
          class="spring-pressable w-full justify-start gap-3 h-11 rounded-xl neu-button !border-transparent"
          :class="providerColors[provider.provider] || providerColors.custom"
          @click="window.location.href = `/api/auth/sso/${provider.provider}/initiate`"
        >
          <svg class="h-5 w-5" viewBox="0 0 24 24" fill="currentColor">
            <path :d="providerIcons[provider.provider] || providerIcons.custom" />
          </svg>
          {{ provider.name }}
        </Button>
      </div>

      <p class="mt-7 text-center text-xs text-muted-foreground">
        {{ $t('auth.noAccount') }}
        <RouterLink
          to="/register"
          class="font-medium text-emerald-600 light:text-emerald-700 hover:underline underline-offset-4"
        >
          {{ $t('auth.signUp') }}
        </RouterLink>
      </p>
    </div>
  </div>
</template>

