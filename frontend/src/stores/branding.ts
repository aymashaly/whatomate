import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import { brandingService, type PlatformBranding, type PlatformBrandingUpdate } from '@/services/api'

// Default branding used when the API hasn't loaded yet, so the UI never
// flashes a missing name / color before the singleton row comes back.
const DEFAULT_BRANDING: PlatformBranding = {
  id: 'default',
  brand_name: 'Whatomate',
  brand_tagline: 'WhatsApp Platform',
  logo_url: '',
  favicon_url: '',
  primary_color: '142 71% 45%',
  accent_color: '173 80% 40%',
  support_email: '',
  support_url: '',
  footer_text: '',
  login_tagline: '',
  updated_at: ''
}

/**
 * Applies the brand to the DOM. Side effects live in one place so views can
 * simply `bindTheme()` after fetching branding without each component
 * repeating the favicon/title/CSS-variable logic.
 *
 *  - <title>           uses brand_name
 *  - <link rel=icon>   swaps to favicon_url (or resets to default)
 *  - --brand-primary   CSS variable consumed by tokens that need to follow
 *                       the admin's chosen accent color
 */
function applyBrandingToDOM(b: PlatformBranding) {
  if (typeof document === 'undefined') return
  const root = document.documentElement

  // Document title — keep "<name> | <page>" pattern where possible.
  const pageSuffix = document.title?.includes(' — ')
    ? ' — ' + document.title.split(' — ').slice(1).join(' — ')
    : ''
  document.title = b.brand_name + pageSuffix

  // Favicon — swap the rel="icon" href.
  let faviconLink = document.head.querySelector('link[rel="icon"]') as HTMLLinkElement | null
  if (!faviconLink) {
    faviconLink = document.createElement('link')
    faviconLink.rel = 'icon'
    document.head.appendChild(faviconLink)
  }
  faviconLink.href = b.favicon_url || '/favicon.svg'

  // CSS variables — only apply if non-empty so a fresh install can fall back
  // to the design tokens defined in index.css.
  if (b.primary_color) root.style.setProperty('--brand-primary-hsl', b.primary_color)
  if (b.accent_color) root.style.setProperty('--brand-accent-hsl', b.accent_color)
}

export const useBrandingStore = defineStore('branding', () => {
  const branding = ref<PlatformBranding>({ ...DEFAULT_BRANDING })
  const loaded = ref(false)
  const loading = ref(false)

  const brandName = computed(() => branding.value.brand_name || DEFAULT_BRANDING.brand_name)
  const brandTagline = computed(() => branding.value.brand_tagline || DEFAULT_BRANDING.brand_tagline)
  const logoUrl = computed(() => branding.value.logo_url)
  const faviconUrl = computed(() => branding.value.favicon_url)
  const loginTagline = computed(() => branding.value.login_tagline)
  const footerText = computed(() => branding.value.footer_text)
  const supportEmail = computed(() => branding.value.support_email)
  const supportUrl = computed(() => branding.value.support_url)
  const primaryColor = computed(() => branding.value.primary_color || DEFAULT_BRANDING.primary_color)
  const accentColor = computed(() => branding.value.accent_color || DEFAULT_BRANDING.accent_color)

  /** Fetch the singleton branding. Safe to call multiple times — concurrent
   *  callers share the in-flight promise via the loading ref. */
  async function load(): Promise<void> {
    if (loading.value) return
    loading.value = true
    try {
      const resp = await brandingService.get()
      const data = (resp.data?.data ?? resp.data) as PlatformBranding
      if (data) {
        branding.value = { ...DEFAULT_BRANDING, ...data }
        loaded.value = true
        applyBrandingToDOM(branding.value)
      }
    } catch (err) {
      // Branding load failures shouldn't break the app — fall back to defaults.
      // eslint-disable-next-line no-console
      console.warn('Failed to load platform branding, using defaults', err)
    } finally {
      loading.value = false
    }
  }

  /** Update branding (super-admin). Returns the saved payload so the caller
   *  can re-bind the theme / toast success without an extra fetch. */
  async function update(payload: PlatformBrandingUpdate): Promise<void> {
    const resp = await brandingService.update(payload)
    const data = (resp.data?.data ?? resp.data) as PlatformBranding
    if (data) {
      branding.value = { ...DEFAULT_BRANDING, ...data }
      applyBrandingToDOM(branding.value)
    }
  }

  /** Upload a logo or favicon. Updates branding in place. */
  async function uploadAsset(file: File, type: 'logo' | 'favicon'): Promise<string> {
    const resp = await brandingService.uploadAsset(file, type)
    const data = (resp.data?.data ?? resp.data) as { url: string }
    if (type === 'logo') branding.value.logo_url = data.url
    else branding.value.favicon_url = data.url
    applyBrandingToDOM(branding.value)
    return data.url
  }

  /** Re-apply the current brand to the DOM (called from main.ts after load()). */
  function bindTheme() {
    applyBrandingToDOM(branding.value)
  }

  // Keep theme in sync whenever branding changes (handles HMR / save).
  watch(branding, (b) => applyBrandingToDOM(b), { deep: true })

  return {
    branding,
    loaded,
    loading,
    brandName,
    brandTagline,
    logoUrl,
    faviconUrl,
    loginTagline,
    footerText,
    supportEmail,
    supportUrl,
    primaryColor,
    accentColor,
    load,
    update,
    uploadAsset,
    bindTheme
  }
})