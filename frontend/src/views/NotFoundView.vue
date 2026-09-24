<script setup lang="ts">
import { computed } from 'vue'
import { Button } from '@/components/ui/button'
import { Home, ArrowLeft } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
// Logged-in users get sent to the app; everyone else gets the landing page.
const homePath = computed(() => (authStore.isAuthenticated ? '/app/dashboard' : '/'))
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-background p-4">
    <div class="text-center">
      <h1 class="text-9xl font-bold text-muted-foreground/20">404</h1>
      <h2 class="text-2xl font-semibold mt-4">{{ $t('notFound.title') }}</h2>
      <p class="text-muted-foreground mt-2 max-w-md">
        {{ $t('notFound.description') }}
      </p>
      <div class="flex items-center justify-center gap-4 mt-8">
        <Button variant="outline" @click="$router.back()">
          <ArrowLeft class="h-4 w-4 mr-2" />
          {{ $t('notFound.goBack') }}
        </Button>
        <RouterLink :to="homePath">
          <Button>
            <Home class="h-4 w-4 mr-2" />
            {{ $t('notFound.goHome') }}
          </Button>
        </RouterLink>
      </div>
    </div>
  </div>
</template>