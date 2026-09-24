<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import {
  MessageSquare,
  ChevronLeft,
  ChevronRight,
  Menu,
  X
} from 'lucide-vue-next'
import { wsService } from '@/services/websocket'
import { authService } from '@/services/api'
import OrganizationSwitcher from './OrganizationSwitcher.vue'
import UserMenu from './UserMenu.vue'
import ActiveCallPanel from '@/components/calling/ActiveCallPanel.vue'
import { ScrollToTop } from '@/components/shared'
import { useBrandingStore } from '@/stores/branding'
import { navigationSections, type NavSection } from './navigation'

useI18n()

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const brandingStore = useBrandingStore()
const isCollapsed = ref(false)
const isMobileMenuOpen = ref(false)

onMounted(() => {
  if (authStore.isAuthenticated) {
    authStore.refreshUserData()
    wsService.connect(async () => {
      try {
        const resp = await authService.getWSToken()
        return resp.data.data.token
      } catch {
        return null
      }
    })
  }
})

function filterItems(items: NavSection['items']) {
  return items
    .filter(item => {
      if (item.childPermissions) {
        return item.childPermissions.some(p => authStore.hasPermission(p, 'read'))
      }
      return !item.permission || authStore.hasPermission(item.permission, 'read')
    })
    .map(item => {
      const filteredChildren = item.children?.filter(
        child => !child.permission || authStore.hasPermission(child.permission, 'read')
      )
      let effectivePath = item.path
      if (item.childPermissions && item.permission && !authStore.hasPermission(item.permission, 'read') && filteredChildren?.length) {
        effectivePath = filteredChildren[0].path
      }
      const originalPath = item.path
      const isActive = originalPath === '/'
        ? route.name === 'dashboard'
        : originalPath === '/chat'
          ? route.name === 'chat' || route.name === 'chat-conversation'
          : route.path.startsWith(originalPath)
      return {
        ...item,
        path: effectivePath,
        active: isActive,
        children: filteredChildren
      }
    })
}

const navSections = computed(() => {
  return navigationSections
    .map(section => ({
      ...section,
      items: filterItems(section.items)
    }))
    .filter(section => section.items.length > 0)
})

const mainSections = computed(() => navSections.value.filter(s => !s.pinBottom))
const bottomSections = computed(() => navSections.value.filter(s => s.pinBottom))

const toggleSidebar = () => {
  isCollapsed.value = !isCollapsed.value
}

// Dynamic brand gradient for the logo tile. When a custom logo is uploaded
// it sits on a white card background, so the gradient is only used as the
// fallback. Both tiles share the same gradient so swapping to the brand
// color (via --brand-primary) is a one-line CSS-variable change.
const sidebarLogoStyle = computed(() => {
  if (brandingStore.logoUrl) return { background: '#ffffff' }
  return {
    background: `linear-gradient(135deg, hsl(${brandingStore.primaryColor}), hsl(${brandingStore.accentColor}))`
  }
})
const mobileLogoStyle = computed(() => {
  if (brandingStore.logoUrl) return { background: '#ffffff' }
  return {
    background: `linear-gradient(135deg, hsl(${brandingStore.primaryColor}), hsl(${brandingStore.accentColor}))`
  }
})

const handleLogout = async () => {
  await authStore.logout()
  router.push('/login')
}
</script>

<template>
  <div class="relative flex h-screen overflow-hidden bg-background text-foreground">
    <a href="#main-content" class="skip-link">{{ $t('nav.skipToMain') }}</a>

    <!-- Mobile header -->
    <header class="glass mx-2 mt-2 z-40 flex h-11 items-center justify-between rounded-2xl px-3 md:hidden">
      <RouterLink to="/app/dashboard" class="flex items-center gap-2.5">
        <div
          class="relative h-7 w-7 overflow-hidden rounded-lg flex items-center justify-center shadow-sm ring-1 ring-white/[0.10]"
          :style="mobileLogoStyle"
        >
          <img
            v-if="brandingStore.logoUrl"
            :src="brandingStore.logoUrl"
            :alt="brandingStore.brandName"
            class="h-full w-full object-contain bg-white"
          />
          <MessageSquare v-else class="h-4 w-4 text-white" />
        </div>
        <span class="display text-sm">{{ brandingStore.brandName }}</span>
      </RouterLink>
      <Button
        variant="ghost"
        size="icon"
        class="spring-pressable h-8 w-8 text-muted-foreground hover:text-foreground hover:bg-black/[0.04]"
        aria-label="Toggle menu"
        :aria-expanded="isMobileMenuOpen"
        @click="isMobileMenuOpen = !isMobileMenuOpen"
      >
        <X v-if="isMobileMenuOpen" class="h-5 w-5" />
        <Menu v-else class="h-5 w-5" />
      </Button>
    </header>

    <!-- Mobile scrim -->
    <div
      v-if="isMobileMenuOpen"
      class="fixed inset-0 z-30 bg-black/40 md:hidden"
      aria-hidden="true"
      @click="isMobileMenuOpen = false"
    />

    <!-- Floating sidebar -->
    <aside
      :class="[
        'glass fixed inset-y-2 left-2 z-40 my-0 flex flex-col overflow-hidden rounded-2xl',
        'transition-[width,transform] duration-[var(--duration-slow)] ease-[var(--ease-spring)]',
        'md:relative md:inset-y-2 md:left-2 md:my-0',
        isMobileMenuOpen ? 'translate-x-0' : '-translate-x-[120%] md:translate-x-0',
        isCollapsed ? 'w-[76px]' : 'w-[272px]'
      ]"
      role="navigation"
      aria-label="Main navigation"
    >
      <!-- Logo row -->
      <div class="hidden md:flex h-16 items-center justify-between px-3.5 border-b border-border/70">
        <RouterLink to="/app/dashboard" class="flex items-center gap-3 overflow-hidden">
          <div
            class="relative h-9 w-9 shrink-0 overflow-hidden rounded-[10px] flex items-center justify-center shadow-sm ring-1 ring-white/[0.12]"
            :style="sidebarLogoStyle"
          >
            <img
              v-if="brandingStore.logoUrl"
              :src="brandingStore.logoUrl"
              :alt="brandingStore.brandName"
              class="h-full w-full object-contain bg-white"
            />
            <MessageSquare v-else class="h-4 w-4 text-white" />
          </div>
          <div v-if="!isCollapsed" class="flex flex-col leading-none">
            <span class="display text-[13px] tracking-tight">{{ brandingStore.brandName }}</span>
            <span class="eyebrow mt-0.5 text-[9px] opacity-60">{{ brandingStore.brandTagline }}</span>
          </div>
        </RouterLink>
        <Button
          v-if="!isCollapsed"
          variant="ghost"
          size="icon"
          class="spring-pressable h-8 w-8 text-muted-foreground hover:text-foreground hover:bg-foreground/[0.05]"
          :aria-label="isCollapsed ? $t('nav.expandSidebar') : $t('nav.collapseSidebar')"
          :aria-expanded="!isCollapsed"
          @click="toggleSidebar"
        >
          <ChevronLeft class="h-4 w-4" />
        </Button>
        <button
          v-else
          class="mx-auto h-9 w-9 flex items-center justify-center rounded-lg text-muted-foreground hover:text-foreground hover:bg-foreground/[0.05] spring-pressable"
          :aria-label="$t('nav.expandSidebar')"
          :aria-expanded="false"
          @click="toggleSidebar"
        >
          <ChevronRight class="h-4 w-4" />
        </button>
      </div>
      <div class="h-14 md:hidden" />

      <OrganizationSwitcher :collapsed="isCollapsed" />

      <!-- Nav with generous spacing -->
      <ScrollArea class="flex-1 pt-3 pb-2">
        <nav class="px-3 space-y-5" role="menubar">
          <template v-for="(section, sIdx) in mainSections" :key="section.label">
            <div>
              <div
                v-if="section.label && !isCollapsed"
                class="px-2.5 pb-2 eyebrow text-[10px] tracking-[0.10em] uppercase text-muted-foreground/80"
              >
                {{ $t(section.label) }}
              </div>
              <!-- Collapsed: divider between groups -->
              <div
                v-else-if="isCollapsed && sIdx > 0"
                class="my-3 mx-2 border-t border-border/70"
              />

              <div class="space-y-1">
                <template v-for="item in section.items" :key="item.path">
                  <RouterLink
                    :to="item.path"
                    :class="[
                      'nav-pill spring-pressable group flex items-center gap-3 rounded-[10px] py-2.5 text-[13px] font-medium',
                      'transition-[background,color,transform] duration-[var(--duration-base)] ease-[var(--ease-spring)]',
                      item.active
                        ? 'shadow-[inset_0_1px_0_0_hsla(0,0%,100%,0.06)] light:shadow-[inset_0_1px_0_0_hsla(0,0%,100%,0.85)]'
                        : 'text-muted-foreground hover:text-foreground hover:bg-foreground/[0.05] light:hover:bg-gray-100',
                      isCollapsed ? 'md:justify-center md:px-2 md:w-12 md:mx-auto' : 'px-3'
                    ]"
                    :data-active="item.active"
                    role="menuitem"
                    :aria-current="item.active ? 'page' : undefined"
                    @click="isMobileMenuOpen = false"
                  >
                    <span
                      :class="[
                        'flex items-center justify-center h-7 w-7 rounded-lg shrink-0 transition-colors',
                        item.active
                          ? ''
                          : 'text-muted-foreground group-hover:text-foreground'
                      ]"
                      aria-hidden="true"
                    >
                      <component :is="item.icon" class="h-4 w-4" />
                    </span>
                    <span :class="isCollapsed && 'md:sr-only'" class="truncate">{{ $t(item.name) }}</span>
                  </RouterLink>

                  <template v-if="item.children && item.active && !isCollapsed">
                    <div class="ml-11 mr-2 mt-1 space-y-0.5 border-l border-border pl-3">
                      <RouterLink
                        v-for="child in item.children"
                        :key="child.path"
                        :to="child.path"
                        :class="[
                          'spring-pressable flex items-center gap-2.5 rounded-lg py-1.5 px-3 text-[13px] font-medium',
                          'transition-[background,color] duration-[var(--duration-base)] ease-[var(--ease-out-quart)]',
                          route.path === child.path
                            ? ''
                            : 'text-muted-foreground hover:text-foreground hover:bg-foreground/[0.04]'
                        ]"
                        role="menuitem"
                        :aria-current="route.path === child.path ? 'page' : undefined"
                        @click="isMobileMenuOpen = false"
                      >
                        <component :is="child.icon" class="h-3.5 w-3.5 shrink-0" aria-hidden="true" />
                        <span>{{ $t(child.name) }}</span>
                      </RouterLink>
                    </div>
                  </template>
                </template>
              </div>
            </div>
          </template>
        </nav>
      </ScrollArea>

      <!-- Bottom: Settings + User -->
      <div v-if="bottomSections.length > 0" class="border-t border-border/70 px-3 pt-3 pb-3 space-y-1">
        <template v-for="section in bottomSections" :key="section.label">
          <template v-for="item in section.items" :key="item.path">
            <RouterLink
              :to="item.path"
              :class="[
                'nav-pill spring-pressable group flex items-center gap-3 rounded-[10px] py-2.5 text-[13px] font-medium',
                'transition-[background,color] duration-[var(--duration-base)] ease-[var(--ease-spring)]',
                item.active
                  ? ''
                  : 'text-muted-foreground hover:text-foreground hover:bg-foreground/[0.05]',
                isCollapsed ? 'md:justify-center md:px-2 md:w-12 md:mx-auto' : 'px-3'
              ]"
              :data-active="item.active"
              role="menuitem"
              :aria-current="item.active ? 'page' : undefined"
              @click="isMobileMenuOpen = false"
            >
              <span
                :class="[
                  'flex items-center justify-center h-7 w-7 rounded-lg shrink-0 transition-colors',
                  item.active ? '' : 'text-muted-foreground group-hover:text-foreground'
                ]"
                aria-hidden="true"
              >
                <component :is="item.icon" class="h-4 w-4" />
              </span>
              <span :class="isCollapsed && 'md:sr-only'">{{ $t(item.name) }}</span>
            </RouterLink>

            <!-- Submenu (children) — same rendering as main sections -->
            <template v-if="item.children && item.active && !isCollapsed">
              <div class="ml-11 mr-2 mt-1 space-y-0.5 border-l border-border pl-3">
                <RouterLink
                  v-for="child in item.children"
                  :key="child.path"
                  :to="child.path"
                  :class="[
                    'spring-pressable flex items-center gap-2.5 rounded-lg py-1.5 px-3 text-[13px] font-medium',
                    'transition-[background,color] duration-[var(--duration-base)] ease-[var(--ease-out-quart)]',
                    route.path === child.path
                      ? ''
                      : 'text-muted-foreground hover:text-foreground hover:bg-foreground/[0.04]'
                  ]"
                  role="menuitem"
                  :aria-current="route.path === child.path ? 'page' : undefined"
                  @click="isMobileMenuOpen = false"
                >
                  <component :is="child.icon" class="h-3.5 w-3.5 shrink-0" aria-hidden="true" />
                  <span>{{ $t(child.name) }}</span>
                </RouterLink>
              </div>
            </template>
          </template>
        </template>
      </div>

      <UserMenu :collapsed="isCollapsed" @logout="handleLogout" />
    </aside>

    <main
      id="main-content"
      class="relative z-10 flex-1 overflow-hidden pt-14 md:pt-0"
      role="main"
    >
      <div class="h-full overflow-y-auto px-2 pt-2 pb-6 md:px-4 md:pt-3">
        <RouterView v-slot="{ Component, route: viewRoute }">
          <Transition name="page" mode="out-in">
            <component
              :is="Component"
              :key="viewRoute.meta.stableKey ? String(viewRoute.name) : viewRoute.path"
            />
          </Transition>
        </RouterView>
        <ScrollToTop />
      </div>
      <ActiveCallPanel />
    </main>
  </div>
</template>

<style scoped>
.page-enter-active,
.page-leave-active {
  transition:
    opacity var(--duration-base) var(--ease-out-quart),
    transform var(--duration-base) var(--ease-spring);
}
.page-enter-from {
  opacity: 0;
  transform: translate3d(0, 4px, 0);
}
.page-leave-to {
  opacity: 0;
  transform: translate3d(0, -2px, 0);
}
@media (prefers-reduced-motion: reduce) {
  .page-enter-active,
  .page-leave-active { transition: none; transform: none; }
  .page-enter-from,
  .page-leave-to { opacity: 1; transform: none; }
}
</style>
