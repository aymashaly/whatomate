<script setup lang="ts">
/**
 * OmniHub marketing landing page.
 *
 * Restructure notes (vs v1):
 *   - Launch offer (100 SAR / 5 months) is no longer splashed in the hero,
 *     pricing cards, sticky CTA, or footer. It now lives only in a
 *     dismissible floating popup that appears after a few seconds.
 *   - Meta Business Partner status gets its own utility bar at the very top
 *     and a dedicated section right after the trust bar.
 *   - Pricing is now three tiers (Starter / Standard / Pro) with an
 *     annual-monthly toggle. The "100 SAR / 5 months" promo is the
 *     Starter plan's launch offer, surfaced in the popup only.
 *   - Sticky mobile CTA becomes a neutral "Get started" that links to
 *     /app/register (the offer lives in the popup, not the sticky bar).
 *
 * Sections:
 *   T0.  Utility bar (Meta Business Partner)
 *   01.  Sticky navbar (logo + nav + language + sign in + Get started)
 *   02.  Hero (no offer glass card — just headline + CTAs + product visual)
 *   03.  Trust bar (industry categories)
 *   04.  Meta partnership section ("Built on the Official WhatsApp Business Platform")
 *   05.  Problem
 *   06.  WhatsApp Coexistence
 *   07.  AI replies
 *   08.  Campaigns
 *   09.  Shared inbox
 *   10.  Bento grid
 *   11.  How it works
 *   12.  Pricing (3 tiers w/ annual toggle)
 *   13.  Social proof placeholder
 *   14.  FAQ
 *   15.  Final CTA
 *   16.  Footer
 *   17.  Demo modal
 *   18.  Launch offer popup (floating button + modal)
 *   19.  Mobile sticky CTA (just "Get started")
 */
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useBrandingStore } from '@/stores/branding'
import BrandMark from '@/components/landing/BrandMark.vue'
import BrandLogo from '@/components/landing/BrandLogo.vue'
import LandingLanguageToggle from '@/components/landing/LandingLanguageToggle.vue'
import {
  Menu, X, Plus, Minus, MessageCircle, Sparkles, Send, Inbox, Users,
  BarChart3, Megaphone, Plug, ShieldCheck, Smartphone, Phone,
  ArrowRight, BadgeCheck, Check
} from 'lucide-vue-next'

// `tm` (translation message) returns the raw value from the locale tree —
// required for array-valued keys (we use this for pricing tier feature
// lists and offer popup features, which are JSON arrays, not strings).
const { t, tm } = useI18n()
const brandingStore = useBrandingStore()

// Brand-driven values — solid brand primary so the marketing surface reads
// flat instead of gradient. Falls back to the violet brand default when the
// store hasn't loaded yet.
const brandName = computed(() => brandingStore.brandName || t('landing.brandName', 'OmniHub') as string)
const heroSolid = computed(() => `hsl(${brandingStore.primaryColor})`)
const heroSolidMuted = computed(() => `hsl(${brandingStore.primaryColor} / 0.1)`)

function applyDocumentTitle() {
  if (typeof document === 'undefined') return
  document.title = `${brandName.value} — ${t('landing.metaTitle')}`
  document.documentElement.setAttribute('description', t('landing.metaDescription'))
}
onMounted(applyDocumentTitle)
watch([brandName, () => t('landing.metaTitle')], applyDocumentTitle)

// FAQ accordion
const openFaq = ref<number | null>(0)
function toggleFaq(i: number) {
  openFaq.value = openFaq.value === i ? null : i
}

// Mobile nav drawer + demo modal
const mobileNavOpen = ref(false)
const demoOpen = ref(false)
function openDemo() { demoOpen.value = true }
function closeDemo() { demoOpen.value = false }

// ---------------------------------------------------------------------------
// Pricing toggle (monthly vs annual). Default to monthly — annual is the
// "save 17%" upsell. Spec-driven, no extra math at render time; both
// values come pre-formatted from i18n.
// ---------------------------------------------------------------------------
const billingAnnual = ref(false)

// ---------------------------------------------------------------------------
// Launch-offer popup. Appears as a floating button in the bottom-right
// after a 3.5s delay, only the first time per session. Clicking either
// the FAB or the secondary CTA inside the offer card opens the popup
// modal. Dismissing persists in sessionStorage so it doesn't reappear
// during the same session.
// ---------------------------------------------------------------------------
const offerFabVisible = ref(false)
const offerModalOpen = ref(false)
const OFFER_DISMISS_KEY = 'landing:offer:dismissed'

function maybeShowOfferFab() {
  if (typeof window === 'undefined') return
  try {
    if (sessionStorage.getItem(OFFER_DISMISS_KEY)) return
  } catch (_) { /* private mode → default to visible */ }
  setTimeout(() => { offerFabVisible.value = true }, 3500)
}
function openOfferModal() {
  offerModalOpen.value = true
}
function closeOfferModal() {
  offerModalOpen.value = false
  try { sessionStorage.setItem(OFFER_DISMISS_KEY, '1') } catch (_) {}
}
function dismissOfferFab() {
  offerFabVisible.value = false
  try { sessionStorage.setItem(OFFER_DISMISS_KEY, '1') } catch (_) {}
}

// Smooth-scroll helper for in-page anchors
function scrollToId(id: string) {
  mobileNavOpen.value = false
  const el = document.getElementById(id)
  if (!el) return
  el.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

const currentYear = new Date().getFullYear()

// Trust-bar industry chips
const trustCategories = [
  { key: 'trustCatEcommerce' },
  { key: 'trustCatRestaurants' },
  { key: 'trustCatClinics' },
  { key: 'trustCatRealEstate' },
  { key: 'trustCatServices' },
  { key: 'trustCatRetail' }
]

// Problem list
const problemKeys = [
  'problemP1', 'problemP2', 'problemP3', 'problemP4', 'problemP5'
] as const

// Campaigns
const campaignSteps = [
  { key: 'campaignsStep1Title', descKey: 'campaignsStep1Desc' },
  { key: 'campaignsStep2Title', descKey: 'campaignsStep2Desc' },
  { key: 'campaignsStep3Title', descKey: 'campaignsStep3Desc' }
] as const

const campaignMetrics = [
  { key: 'campaignsMetricSent',     value: '1,248' },
  { key: 'campaignsMetricDelivered', value: '1,201' },
  { key: 'campaignsMetricRead',      value: '942' },
  { key: 'campaignsMetricReplies',   value: '318' }
]

// Inbox features
const inboxFeatures = [
  { key: 'inboxF1Title', icon: Users },
  { key: 'inboxF2Title', icon: Inbox },
  { key: 'inboxF3Title', icon: MessageCircle },
  { key: 'inboxF4Title', icon: ShieldCheck },
  { key: 'inboxF5Title', icon: Smartphone },
  { key: 'inboxF6Title', icon: Sparkles }
] as const

// Bento
const bentoTiles = [
  { key: 'bentoAiTitle', descKey: 'bentoAiDesc', icon: Sparkles, accent: 'from-violet-500/20 to-indigo-500/10',  span: 'md:col-span-2 md:row-span-2' },
  { key: 'bentoInboxTitle', descKey: 'bentoInboxDesc', icon: Inbox, accent: 'from-sky-500/15 to-cyan-500/10',    span: '' },
  { key: 'bentoCampaignsTitle', descKey: 'bentoCampaignsDesc', icon: Megaphone, accent: 'from-fuchsia-500/15 to-pink-500/10', span: '' },
  { key: 'bentoCoexTitle', descKey: 'bentoCoexDesc', icon: Phone, accent: 'from-emerald-500/15 to-teal-500/10', span: '' },
  { key: 'bentoTeamTitle', descKey: 'bentoTeamDesc', icon: Users, accent: 'from-amber-500/15 to-orange-500/10',  span: '' },
  { key: 'bentoAnalyticsTitle', descKey: 'bentoAnalyticsDesc', icon: BarChart3, accent: 'from-rose-500/15 to-red-500/10', span: 'md:col-span-2' }
] as const

// How-it-works
const howSteps = [
  { key: 'howStep1Title', descKey: 'howStep1Desc', num: '01', icon: Plug },
  { key: 'howStep2Title', descKey: 'howStep2Desc', num: '02', icon: Sparkles },
  { key: 'howStep3Title', descKey: 'howStep3Desc', num: '03', icon: Send }
] as const

// FAQ
const faqs = [
  { qKey: 'faqQ1', aKey: 'faqA1' },
  { qKey: 'faqQ2', aKey: 'faqA2' },
  { qKey: 'faqQ3', aKey: 'faqA3' },
  { qKey: 'faqQ4', aKey: 'faqA4' },
  { qKey: 'faqQ5', aKey: 'faqA5' },
  { qKey: 'faqQ6', aKey: 'faqA6' },
  { qKey: 'faqQ7', aKey: 'faqA7' }
] as const

// Meta partner section — 4 differentiators
const metaPoints = [
  { key: 'metaPartnerP1Title', descKey: 'metaPartnerP1Desc', icon: BadgeCheck },
  { key: 'metaPartnerP2Title', descKey: 'metaPartnerP2Desc', icon: ShieldCheck },
  { key: 'metaPartnerP3Title', descKey: 'metaPartnerP3Desc', icon: MessageCircle },
  { key: 'metaPartnerP4Title', descKey: 'metaPartnerP4Desc', icon: Plug }
] as const

// Pricing tiers — id maps to i18n keys.
const pricingTiers = [
  {
    id: 'starter',
    featured: false,
    ctaKey: 'pricingStarterCta',
    nameKey: 'pricingStarterName',
    tagKey:  'pricingStarterTag',
    descKey: 'pricingStarterDesc',
    monthlyKey: 'pricingStarterPriceMonthly',
    annualKey:  'pricingStarterPriceAnnual',
    annualSaveKey: 'pricingStarterAnnualSave',
    featuresPath: 'landing.pricingStarterFeatures',
    noIncludesKey: 'pricingStarterNoIncludes'
  },
  {
    id: 'standard',
    featured: true,
    ctaKey: 'pricingStandardCta',
    nameKey: 'pricingStandardName',
    tagKey:  'pricingStandardTag',
    descKey: 'pricingStandardDesc',
    monthlyKey: 'pricingStandardPriceMonthly',
    annualKey:  'pricingStandardPriceAnnual',
    annualSaveKey: 'pricingStandardAnnualSave',
    featuresPath: 'landing.pricingStandardFeatures',
    noIncludesKey: null
  },
  {
    id: 'pro',
    featured: false,
    ctaKey: 'pricingProCta',
    nameKey: 'pricingProName',
    tagKey:  'pricingProTag',
    descKey: 'pricingProDesc',
    monthlyKey: 'pricingProPriceMonthly',
    annualKey:  'pricingProPriceAnnual',
    annualSaveKey: 'pricingProAnnualSave',
    featuresPath: 'landing.pricingProFeatures',
    noIncludesKey: null
  }
] as const

// vue-i18n `tm()` resolves an array-typed i18n key into a string[] —
// needed because `t()` doesn't natively return arrays.
function tierFeatures(path: string): string[] {
  const raw = (tm as unknown as (k: string) => unknown)(path)
  return Array.isArray(raw) ? (raw as string[]) : []
}

// Demo modal steps
const demoSteps = [
  { key: 'demoStep1' },
  { key: 'demoStep2' },
  { key: 'demoStep3' },
  { key: 'demoStep4' }
] as const

// Offer popup feature list comes through i18n as an array
function offerFeatureList(): string[] {
  const raw = (tm as unknown as (k: string) => unknown)('landing.popupOfferFeatures')
  return Array.isArray(raw) ? (raw as string[]) : []
}

// Close on Esc — covers demo, offer modal, and mobile nav
function onKeydown(e: KeyboardEvent) {
  if (e.key !== 'Escape') return
  if (demoOpen.value)      demoOpen.value = false
  if (offerModalOpen.value) offerModalOpen.value = false
  if (mobileNavOpen.value) mobileNavOpen.value = false
}

onMounted(() => {
  if (typeof window !== 'undefined') {
    window.addEventListener('keydown', onKeydown)
  }
  maybeShowOfferFab()
})

// Pricing primary CTA — points to WhatsApp chat. Brief allows this until
// Stripe/Mada integration lands.
const whatsappLink = 'https://wa.me/966000000000?text=' + encodeURIComponent(t('landing.heroPrimaryCta') as string)
</script>

<template>
  <div class="landing-shell min-h-screen relative overflow-x-hidden bg-slate-50 text-slate-900">
    <!-- Ambient color layer -->
    <div aria-hidden="true" class="pointer-events-none fixed inset-0 -z-10 overflow-hidden">
      <div class="absolute -top-40 -start-40 h-[28rem] w-[28rem] rounded-full bg-violet-400/30 mix-blend-multiply blur-3xl opacity-70 animate-shimmer"></div>
      <div class="absolute top-1/3 -end-32 h-[26rem] w-[26rem] rounded-full bg-sky-400/30 mix-blend-multiply blur-3xl opacity-60 animate-shimmer"></div>
      <div class="absolute bottom-0 start-1/3 h-[24rem] w-[24rem] rounded-full bg-cyan-300/30 mix-blend-multiply blur-3xl opacity-60 animate-shimmer"></div>
      <div class="absolute inset-0 bg-[radial-gradient(circle_at_center,rgba(255,255,255,0.6)_0%,rgba(255,255,255,0)_60%)]"></div>
    </div>

    <!-- ============================================================== -->
    <!-- T0. META UTILITY BAR — sits at the very top of the page, before -->
    <!-- the navbar. Establishes Meta Business Partner credibility.       -->
    <!-- ============================================================== -->
    <div class="border-b border-slate-200/60 bg-gradient-to-r from-slate-900 via-slate-800 to-slate-900 text-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 py-2 flex items-center justify-center sm:justify-between gap-3">
        <div class="flex items-center gap-2 text-[12px] sm:text-xs font-medium">
          <span class="inline-flex h-6 w-6 items-center justify-center rounded-md bg-white/10 text-white text-[11px] font-bold border border-white/20">M</span>
          <span class="font-semibold tracking-wide">{{ t('landing.metaPartnerTopBadge') }}</span>
          <span class="hidden sm:inline text-white/60">·</span>
          <span class="hidden sm:inline text-white/80">{{ t('landing.metaPartnerTopSub') }}</span>
        </div>
        <div class="hidden sm:flex items-center gap-1.5 text-[11px] text-white/70">
          <BadgeCheck class="h-3.5 w-3.5 text-cyan-300" />
          {{ t('landing.heroEyebrow') }}
        </div>
      </div>
    </div>

    <!-- ============================================================== -->
    <!-- 01. NAVBAR                                                     -->
    <!-- ============================================================== -->
    <header class="sticky top-0 z-40">
      <div class="max-w-7xl mx-auto px-4 sm:px-6">
        <div class="mt-3 sm:mt-4 glass-panel flex items-center justify-between px-3 sm:px-4 py-2.5 rounded-2xl">
          <RouterLink to="/" class="flex items-center gap-2.5 group">
            <BrandLogo :size="36" />
            <div class="flex flex-col leading-none">
              <span class="font-semibold text-[15px] tracking-tight text-slate-900">{{ brandName }}</span>
              <span class="text-[10px] mt-1 text-slate-500 tracking-wide">{{ t('landing.brandTagline') }}</span>
            </div>
          </RouterLink>

          <nav class="hidden lg:flex items-center gap-1 text-sm">
            <button class="nav-link-anim px-3 py-2 rounded-lg text-slate-600 hover:text-slate-900" @click="scrollToId('features')">{{ t('landing.navFeatures') }}</button>
            <button class="nav-link-anim px-3 py-2 rounded-lg text-slate-600 hover:text-slate-900" @click="scrollToId('ai')">{{ t('landing.navAi') }}</button>
            <button class="nav-link-anim px-3 py-2 rounded-lg text-slate-600 hover:text-slate-900" @click="scrollToId('campaigns')">{{ t('landing.navCampaigns') }}</button>
            <button class="nav-link-anim px-3 py-2 rounded-lg text-slate-600 hover:text-slate-900" @click="scrollToId('pricing')">{{ t('landing.navPricing') }}</button>
            <button class="nav-link-anim px-3 py-2 rounded-lg text-slate-600 hover:text-slate-900" @click="scrollToId('faq')">{{ t('landing.navFaq') }}</button>
          </nav>

          <div class="flex items-center gap-2">
            <LandingLanguageToggle class="hidden sm:inline-flex" />
            <RouterLink to="/login" class="hidden md:inline-flex text-sm font-medium text-slate-700 hover:text-slate-900 px-3 py-2 rounded-lg hover:bg-white/50 transition-colors">
              {{ t('landing.signIn') }}
            </RouterLink>
            <RouterLink to="/register">
              <button class="inline-flex items-center gap-1.5 rounded-xl px-3.5 py-2 text-sm font-semibold text-white shadow-md hover:shadow-lg transition-shadow" :style="{ background: heroSolid }">
                {{ t('landing.getStarted') }}
                <ArrowRight class="h-3.5 w-3.5 flip-arrow-end" />
              </button>
            </RouterLink>
            <button class="lg:hidden p-2 rounded-lg text-slate-700 hover:bg-white/60" :aria-label="mobileNavOpen ? 'Close menu' : 'Open menu'" @click="mobileNavOpen = !mobileNavOpen">
              <component :is="mobileNavOpen ? X : Menu" class="h-5 w-5" />
            </button>
          </div>
        </div>
      </div>

      <!-- Mobile drawer -->
      <div v-if="mobileNavOpen" class="lg:hidden mt-2 max-w-7xl mx-auto px-4 sm:px-6">
        <div class="glass-panel rounded-2xl p-3 flex flex-col gap-1 text-base">
          <button class="text-start px-3 py-3 rounded-lg hover:bg-white/60" @click="scrollToId('features')">{{ t('landing.navFeatures') }}</button>
          <button class="text-start px-3 py-3 rounded-lg hover:bg-white/60" @click="scrollToId('ai')">{{ t('landing.navAi') }}</button>
          <button class="text-start px-3 py-3 rounded-lg hover:bg-white/60" @click="scrollToId('campaigns')">{{ t('landing.navCampaigns') }}</button>
          <button class="text-start px-3 py-3 rounded-lg hover:bg-white/60" @click="scrollToId('pricing')">{{ t('landing.navPricing') }}</button>
          <button class="text-start px-3 py-3 rounded-lg hover:bg-white/60" @click="scrollToId('faq')">{{ t('landing.navFaq') }}</button>
          <div class="flex items-center justify-between pt-2 mt-2 border-t border-slate-200/60">
            <RouterLink to="/login" class="text-slate-700 px-3 py-2 rounded-lg hover:bg-white/60">{{ t('landing.signIn') }}</RouterLink>
            <LandingLanguageToggle />
          </div>
        </div>
      </div>
    </header>

    <!-- ============================================================== -->
    <!-- 02. HERO                                                        -->
    <!-- ============================================================== -->
    <section class="relative pt-12 pb-20 md:pt-20 md:pb-28 px-4 sm:px-6">
      <div class="max-w-7xl mx-auto grid lg:grid-cols-12 gap-10 lg:gap-16 items-center">
        <div class="lg:col-span-6 animate-fade-up">
          <div class="inline-flex items-center gap-2 rounded-full glass-tile px-3 py-1.5 text-xs font-semibold text-slate-700">
            <span class="h-1.5 w-1.5 rounded-full bg-violet-500 animate-pulse"></span>
            {{ t('landing.heroEyebrow') }}
          </div>

          <h1 class="mt-6 text-4xl sm:text-5xl md:text-6xl font-bold tracking-tight leading-[1.05] text-slate-900">
            <template v-if="(t('landing.heroTitle') as string).includes(t('landing.heroHighlightedWord') as string)">
              <span>{{ (t('landing.heroTitle') as string).split(t('landing.heroHighlightedWord') as string)[0] }}</span><span class="brand-fill-text">{{ t('landing.heroHighlightedWord') }}</span><span>{{ (t('landing.heroTitle') as string).split(t('landing.heroHighlightedWord') as string)[1] }}</span>
            </template>
            <template v-else>
              {{ t('landing.heroTitle') }}
            </template>
          </h1>

          <p class="mt-5 text-base sm:text-lg text-slate-600 leading-relaxed max-w-xl">
            {{ t('landing.heroSubtitle') }}
          </p>

          <div class="mt-7 flex flex-wrap items-center gap-3">
            <RouterLink to="/register">
              <button class="inline-flex items-center gap-2 rounded-2xl px-5 py-3 text-sm sm:text-base font-semibold text-white shadow-xl hover:shadow-2xl transition-shadow" :style="{ background: heroSolid }">
                {{ t('landing.heroPrimaryCta') }}
                <ArrowRight class="h-4 w-4 flip-arrow-end" />
              </button>
            </RouterLink>
            <button class="inline-flex items-center gap-2 rounded-2xl px-5 py-3 text-sm sm:text-base font-semibold text-slate-800 glass-card hover:translate-y-[-1px] transition-transform" @click="openDemo">
              <span class="inline-flex h-7 w-7 items-center justify-center rounded-full brand-fill text-white">
                <svg class="h-3 w-3" viewBox="0 0 24 24" fill="currentColor"><path d="M8 5v14l11-7z"/></svg>
              </span>
              {{ t('landing.heroSecondaryCta') }}
            </button>
          </div>

          <!-- Trust micro-row: no offer. Just credibility. -->
          <div class="mt-8 flex flex-wrap items-center gap-x-5 gap-y-2 text-xs text-slate-500">
            <span class="inline-flex items-center gap-1.5">
              <BadgeCheck class="h-4 w-4 text-emerald-600" />
              {{ t('landing.metaPartnerTopBadge') }}
            </span>
            <span class="inline-flex items-center gap-1.5">
              <ShieldCheck class="h-4 w-4 text-emerald-600" />
              {{ t('landing.metaPartnerTopSub') }}
            </span>
          </div>
        </div>

        <!-- Right: 3D product visual -->
        <div class="lg:col-span-6 animate-fade-up-2">
          <div class="relative">
            <div aria-hidden="true" class="absolute inset-0 -z-10">
              <div class="absolute inset-x-10 top-10 bottom-10 rounded-3xl bg-gradient-to-br from-violet-300/40 via-sky-300/30 to-cyan-300/40 blur-3xl"></div>
            </div>

            <div class="glass-panel relative p-4 sm:p-5 rounded-3xl overflow-hidden">
              <div class="flex items-center justify-between mb-3 px-1">
                <div class="flex items-center gap-2">
                  <BrandLogo :size="28" />
                  <div class="text-sm font-semibold text-slate-800">{{ brandName }}</div>
                </div>
                <div class="flex items-center gap-1.5 text-[11px] text-emerald-700">
                  <span class="h-1.5 w-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
                  live
                </div>
              </div>

              <div class="glass-tile rounded-2xl p-3 sm:p-4 space-y-3">
                <div class="flex justify-end">
                  <div class="max-w-[80%] rounded-2xl rounded-tr-md bg-violet-500/10 px-3 py-2 text-sm text-slate-800">
                    {{ t('landing.aiExample1Customer') }}
                  </div>
                </div>
                <div class="flex items-center gap-2">
                  <span class="inline-flex h-5 items-center gap-1 rounded-full px-2 text-[10px] font-bold uppercase tracking-wider text-white brand-fill">
                    <Sparkles class="h-3 w-3" /> {{ t('landing.aiBadge') }}
                  </span>
                </div>
                <div class="flex justify-start">
                  <div class="max-w-[80%] rounded-2xl rounded-tl-md bg-white/80 px-3 py-2 text-sm text-slate-800 shadow-sm">
                    {{ t('landing.aiExample1Reply') }}
                  </div>
                </div>
                <div class="flex justify-end">
                  <div class="inline-flex items-center gap-2 rounded-full bg-amber-100/80 px-3 py-1.5 text-[11px] font-medium text-amber-900">
                    <span class="h-1.5 w-1.5 rounded-full bg-amber-500"></span>
                    {{ t('landing.aiHumanHandoff') }}
                  </div>
                </div>
              </div>

              <div class="mt-3 grid grid-cols-4 gap-2">
                <div class="glass-tile p-2 sm:p-3 text-center">
                  <div class="text-[10px] text-slate-500 uppercase tracking-wide">{{ t('landing.campaignsMetricSent') }}</div>
                  <div class="text-sm font-bold text-slate-900 mt-0.5">1,248</div>
                </div>
                <div class="glass-tile p-2 sm:p-3 text-center">
                  <div class="text-[10px] text-slate-500 uppercase tracking-wide">{{ t('landing.campaignsMetricDelivered') }}</div>
                  <div class="text-sm font-bold text-slate-900 mt-0.5">1,201</div>
                </div>
                <div class="glass-tile p-2 sm:p-3 text-center">
                  <div class="text-[10px] text-slate-500 uppercase tracking-wide">{{ t('landing.campaignsMetricRead') }}</div>
                  <div class="text-sm font-bold text-slate-900 mt-0.5">942</div>
                </div>
                <div class="glass-tile p-2 sm:p-3 text-center">
                  <div class="text-[10px] text-slate-500 uppercase tracking-wide">{{ t('landing.campaignsMetricReplies') }}</div>
                  <div class="text-sm font-bold text-slate-900 mt-0.5">318</div>
                </div>
              </div>
            </div>

            <div class="absolute -bottom-8 -end-2 sm:-end-6 w-44 sm:w-52 animate-brand-float">
              <div class="glass-panel rounded-[2rem] p-2 shadow-2xl">
                <div class="rounded-[1.5rem] bg-slate-900 p-3 space-y-2 overflow-hidden">
                  <div class="flex items-center justify-between text-[10px] text-slate-300 px-1">
                    <span>9:41</span>
                    <span class="h-1 w-6 rounded-full bg-slate-700"></span>
                  </div>
                  <div class="space-y-1.5 px-1">
                    <div class="rounded-xl bg-slate-800 px-2 py-1.5 text-[11px] text-slate-200 max-w-[85%] mr-auto">
                      {{ t('landing.aiExample1Customer') }}
                    </div>
                    <div class="rounded-xl bg-emerald-600/90 px-2 py-1.5 text-[11px] text-white max-w-[85%] ml-auto">
                      {{ t('landing.aiExample1Reply') }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ============================================================== -->
    <!-- 03. TRUST BAR                                                  -->
    <!-- ============================================================== -->
    <section class="px-4 sm:px-6 py-10">
      <div class="max-w-6xl mx-auto glass-card rounded-3xl px-4 sm:px-8 py-6 sm:py-8 text-center animate-fade-up">
        <div class="text-sm font-medium text-slate-600">
          {{ t('landing.trustEyebrow') }}
        </div>
        <div class="mt-4 flex flex-wrap items-center justify-center gap-2 sm:gap-3 text-xs sm:text-sm">
          <span v-for="(cat, i) in trustCategories" :key="cat.key" class="inline-flex items-center gap-2 rounded-full bg-white/70 px-3 py-1.5 text-slate-700 border border-slate-200/60">
            <span class="h-1.5 w-1.5 rounded-full brand-fill" :style="`opacity: ${0.6 + (i / trustCategories.length) * 0.4}`"></span>
            {{ t(`landing.${cat.key}`) }}
          </span>
        </div>
        <div class="mt-4 text-xs text-slate-500">
          {{ t('landing.trustSubtext') }}
        </div>
      </div>
    </section>

    <!-- ============================================================== -->
    <!-- 04. META BUSINESS PARTNER SECTION                             -->
    <!-- ============================================================== -->
    <section class="px-4 sm:px-6 py-14 sm:py-20 relative">
      <div aria-hidden="true" class="absolute inset-0 -z-10 overflow-hidden">
        <div class="absolute top-1/2 start-1/2 -translate-x-1/2 -translate-y-1/2 h-96 w-96 rounded-full bg-slate-300/40 blur-3xl"></div>
      </div>
      <div class="max-w-6xl mx-auto">
        <div class="glass-card relative overflow-hidden p-8 sm:p-12 rounded-3xl">
          <div aria-hidden="true" class="absolute -top-20 -end-20 h-60 w-60 rounded-full bg-violet-300/30 blur-3xl"></div>
          <div aria-hidden="true" class="absolute -bottom-24 -start-24 h-60 w-60 rounded-full bg-cyan-300/30 blur-3xl"></div>

          <div class="relative grid md:grid-cols-12 gap-10 items-start">
            <div class="md:col-span-5">
              <!-- Partner badge -->
              <div class="inline-flex items-center gap-2 rounded-full bg-slate-900 text-white px-3 py-1.5 text-xs font-semibold">
                <span class="inline-flex h-5 w-5 items-center justify-center rounded bg-white text-slate-900 text-[11px] font-bold">M</span>
                <span>{{ t('landing.metaPartnerTopBadge') }}</span>
                <span class="text-white/60">·</span>
                <span class="text-white/80 font-normal">{{ t('landing.metaPartnerTopSub') }}</span>
              </div>
              <div class="mt-4 text-xs uppercase tracking-wider font-bold brand-fill-text inline-block">{{ t('landing.metaPartnerEyebrow') }}</div>
              <h2 class="mt-2 text-3xl sm:text-4xl md:text-5xl font-bold tracking-tight text-slate-900 leading-tight">
                {{ t('landing.metaPartnerTitle') }}
              </h2>
              <p class="mt-4 text-base sm:text-lg text-slate-600 leading-relaxed">
                {{ t('landing.metaPartnerSubtitle') }}
              </p>
            </div>

            <div class="md:col-span-7 grid sm:grid-cols-2 gap-4">
              <div v-for="(p, i) in metaPoints" :key="p.key" class="glass-tile p-5 rounded-2xl">
                <span class="inline-flex h-10 w-10 items-center justify-center rounded-xl brand-fill text-white">
                  <component :is="p.icon" class="h-5 w-5" />
                </span>
                <div class="mt-3 text-base font-bold text-slate-900">{{ t(`landing.${p.key}`) }}</div>
                <div class="mt-1 text-sm text-slate-600 leading-relaxed">{{ t(`landing.${p.descKey}`) }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ============================================================== -->
    <!-- 05. PROBLEM                                                    -->
    <!-- ============================================================== -->
    <section class="px-4 sm:px-6 py-16 sm:py-20">
      <div class="max-w-5xl mx-auto text-center animate-fade-up">
        <div class="text-xs uppercase tracking-wider font-bold brand-fill-text inline-block">{{ t('landing.problemEyebrow') }}</div>
        <h2 class="mt-3 text-3xl sm:text-4xl md:text-5xl font-bold tracking-tight text-slate-900 max-w-3xl mx-auto leading-tight">
          {{ t('landing.problemTitle') }}
        </h2>
        <div class="mt-10 grid sm:grid-cols-2 lg:grid-cols-5 gap-3 text-start">
          <div v-for="(k, i) in problemKeys" :key="k" class="glass-card p-4 sm:p-5 rounded-2xl">
            <span class="inline-flex h-7 w-7 items-center justify-center rounded-lg brand-fill text-white text-xs font-bold">{{ i + 1 }}</span>
            <div class="mt-3 text-sm font-semibold text-slate-800">{{ t(`landing.${k}`) }}</div>
          </div>
        </div>
        <p class="mt-10 text-lg sm:text-xl text-slate-700 font-medium">
          {{ t('landing.problemPivot') }}
        </p>
      </div>
    </section>

    <!-- ============================================================== -->
    <!-- 06. WHATSAPP COEXISTENCE                                       -->
    <!-- ============================================================== -->
    <section id="features" class="px-4 sm:px-6 py-16 sm:py-24 relative">
      <div aria-hidden="true" class="absolute inset-0 -z-10 overflow-hidden">
        <div class="absolute top-1/2 start-1/4 h-72 w-72 -translate-y-1/2 rounded-full bg-emerald-200/30 blur-3xl"></div>
        <div class="absolute top-1/2 end-1/4 h-72 w-72 -translate-y-1/2 rounded-full bg-violet-300/40 blur-3xl"></div>
      </div>

      <div class="max-w-7xl mx-auto animate-fade-up">
        <div class="text-center max-w-3xl mx-auto">
          <div class="text-xs uppercase tracking-wider font-bold text-emerald-600 inline-block">{{ t('landing.coexEyebrow') }}</div>
          <h2 class="mt-3 text-3xl sm:text-4xl md:text-5xl font-bold tracking-tight text-slate-900 leading-tight">
            {{ t('landing.coexTitle') }}
          </h2>
          <p class="mt-4 text-base sm:text-lg text-slate-600 max-w-2xl mx-auto leading-relaxed">
            {{ t('landing.coexSubtitle') }}
          </p>
        </div>

        <div class="mt-14 grid md:grid-cols-2 gap-6 items-stretch relative">
          <div aria-hidden="true" class="hidden md:block absolute inset-y-12 start-1/2 -translate-x-1/2 w-1 brand-fill rounded-full opacity-50"></div>
          <div aria-hidden="true" class="hidden md:block absolute start-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 h-6 w-6 rounded-full brand-fill shadow-lg animate-pulse"></div>

          <div class="glass-panel p-6 sm:p-8 rounded-3xl relative overflow-hidden">
            <div aria-hidden="true" class="absolute -top-10 -end-10 h-32 w-32 rounded-full bg-emerald-300/50 blur-3xl"></div>
            <div class="flex items-center gap-3">
              <span class="inline-flex h-10 w-10 items-center justify-center rounded-xl bg-emerald-500/15 text-emerald-600">
                <Smartphone class="h-5 w-5" />
              </span>
              <div>
                <div class="text-base font-bold text-slate-900">{{ t('landing.coexLeftTitle') }}</div>
                <div class="text-xs text-slate-500 mt-0.5">{{ t('landing.coexLeftSubtitle') }}</div>
              </div>
            </div>
            <div class="mt-6 flex justify-center">
              <div class="w-56 sm:w-64 rounded-[2.2rem] bg-slate-900 p-2 shadow-2xl">
                <div class="rounded-[1.7rem] bg-gradient-to-b from-slate-800 to-slate-900 p-3 space-y-2 overflow-hidden">
                  <div class="flex items-center justify-between text-[10px] text-slate-400 px-2">
                    <span>9:41</span>
                    <span>WhatsApp</span>
                  </div>
                  <div class="space-y-1.5">
                    <div class="rounded-2xl rounded-bl-sm bg-slate-700/80 px-3 py-2 text-[12px] text-slate-100 max-w-[80%]">
                      {{ t('landing.aiExample1Customer') }}
                    </div>
                    <div class="rounded-2xl rounded-br-sm bg-emerald-600 px-3 py-2 text-[12px] text-white max-w-[80%] ms-auto">
                      {{ t('landing.aiExample1Reply') }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="glass-panel p-6 sm:p-8 rounded-3xl relative overflow-hidden">
            <div aria-hidden="true" class="absolute -top-10 -start-10 h-32 w-32 rounded-full bg-violet-400/40 blur-3xl"></div>
            <div class="flex items-center gap-3">
              <span class="inline-flex h-10 w-10 items-center justify-center rounded-xl brand-fill text-white">
                <BrandLogo :size="22" />
              </span>
              <div>
                <div class="text-base font-bold text-slate-900">{{ t('landing.coexRightTitle') }}</div>
                <div class="text-xs text-slate-500 mt-0.5">{{ t('landing.coexRightSubtitle') }}</div>
              </div>
            </div>
            <div class="mt-6 glass-tile p-3 rounded-2xl space-y-2.5">
              <div class="flex items-center justify-between text-[10px] text-slate-500 px-1">
                <span>{{ t('landing.bentoInboxTitle') }}</span>
                <span class="text-emerald-600">●</span>
              </div>
              <div class="space-y-1.5 px-1">
                <div class="rounded-xl bg-white/70 px-2.5 py-1.5 text-[12px] text-slate-700 max-w-[80%] ml-auto">{{ t('landing.aiExample1Customer') }}</div>
                <div class="rounded-xl bg-violet-100 px-2.5 py-1.5 text-[12px] text-slate-700 max-w-[80%] flex items-start gap-1">
                  <span class="inline-flex h-4 items-center rounded-full px-1.5 text-[8px] font-bold uppercase tracking-wider text-white brand-fill">{{ t('landing.aiBadge') }}</span>
                  <span>{{ t('landing.aiExample2Reply') }}</span>
                </div>
              </div>
              <div class="grid grid-cols-3 gap-1.5 pt-1">
                <div class="rounded-lg bg-white/60 p-1.5 text-center"><div class="text-[8px] text-slate-500">{{ t('landing.campaignsMetricSent') }}</div><div class="text-[11px] font-bold">1.2k</div></div>
                <div class="rounded-lg bg-white/60 p-1.5 text-center"><div class="text-[8px] text-slate-500">{{ t('landing.campaignsMetricRead') }}</div><div class="text-[11px] font-bold">942</div></div>
                <div class="rounded-lg bg-white/60 p-1.5 text-center"><div class="text-[8px] text-slate-500">{{ t('landing.campaignsMetricReplies') }}</div><div class="text-[11px] font-bold">318</div></div>
              </div>
            </div>
          </div>
        </div>

        <p class="mt-12 text-center text-xl sm:text-2xl font-bold tracking-tight text-slate-900 max-w-3xl mx-auto">
          <span class="brand-fill-text">{{ t('landing.coexHeadline') }}</span>
        </p>
      </div>
    </section>

    <!-- ============================================================== -->
    <!-- 07. AI REPLIES                                                 -->
    <!-- ============================================================== -->
    <section id="ai" class="px-4 sm:px-6 py-16 sm:py-20">
      <div class="max-w-6xl mx-auto grid lg:grid-cols-12 gap-10 items-center animate-fade-up">
        <div class="lg:col-span-5">
          <div class="text-xs uppercase tracking-wider font-bold brand-fill-text inline-block">{{ t('landing.aiEyebrow') }}</div>
          <h2 class="mt-3 text-3xl sm:text-4xl md:text-5xl font-bold tracking-tight text-slate-900 leading-tight">
            {{ t('landing.aiTitle') }}
          </h2>
          <p class="mt-5 text-base sm:text-lg text-slate-600 leading-relaxed">
            {{ t('landing.aiSubtitle') }}
          </p>
          <div class="mt-6 inline-flex items-start gap-3 rounded-2xl bg-amber-50 border border-amber-200/60 p-3.5">
            <span class="inline-flex h-6 w-6 items-center justify-center rounded-full bg-amber-500 text-white text-[11px] font-bold">!</span>
            <div class="text-sm text-amber-900">{{ t('landing.aiHumanHandoff') }}</div>
          </div>
        </div>

        <div class="lg:col-span-7">
          <div class="glass-panel p-4 sm:p-6 rounded-3xl space-y-3">
            <div class="flex justify-end">
              <div class="max-w-[85%] rounded-2xl rounded-tr-md bg-violet-500/10 px-4 py-2.5 text-sm text-slate-800">{{ t('landing.aiExample1Customer') }}</div>
            </div>
            <div class="flex items-center gap-2">
              <span class="inline-flex h-5 items-center gap-1 rounded-full px-2 text-[10px] font-bold uppercase tracking-wider text-white brand-fill">
                <Sparkles class="h-3 w-3" /> {{ t('landing.aiBadge') }}
              </span>
            </div>
            <div class="flex justify-start">
              <div class="max-w-[85%] rounded-2xl rounded-tl-md bg-white/85 px-4 py-2.5 text-sm text-slate-800 shadow-sm">{{ t('landing.aiExample1Reply') }}</div>
            </div>
            <div class="flex justify-end pt-2">
              <div class="max-w-[85%] rounded-2xl rounded-tr-md bg-violet-500/10 px-4 py-2.5 text-sm text-slate-800">{{ t('landing.aiExample2Customer') }}</div>
            </div>
            <div class="flex items-center gap-2">
              <span class="inline-flex h-5 items-center gap-1 rounded-full px-2 text-[10px] font-bold uppercase tracking-wider text-white brand-fill">
                <Sparkles class="h-3 w-3" /> {{ t('landing.aiBadge') }}
              </span>
            </div>
            <div class="flex justify-start">
              <div class="max-w-[85%] rounded-2xl rounded-tl-md bg-white/85 px-4 py-2.5 text-sm text-slate-800 shadow-sm">{{ t('landing.aiExample2Reply') }}</div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ============================================================== -->
    <!-- 08. CAMPAIGNS                                                   -->
    <!-- ============================================================== -->
    <section id="campaigns" class="px-4 sm:px-6 py-16 sm:py-20">
      <div class="max-w-6xl mx-auto animate-fade-up">
        <div class="text-center max-w-3xl mx-auto">
          <div class="text-xs uppercase tracking-wider font-bold brand-fill-text inline-block">{{ t('landing.campaignsEyebrow') }}</div>
          <h2 class="mt-3 text-3xl sm:text-4xl md:text-5xl font-bold tracking-tight text-slate-900 leading-tight">
            {{ t('landing.campaignsTitle') }}
          </h2>
          <p class="mt-4 text-base sm:text-lg text-slate-600 leading-relaxed">
            {{ t('landing.campaignsSubtitle') }}
          </p>
        </div>

        <div class="mt-12 grid md:grid-cols-3 gap-5">
          <div v-for="(step, i) in campaignSteps" :key="step.key" class="glass-card p-6 rounded-2xl">
            <div class="flex items-center justify-between">
              <span class="inline-flex h-10 w-10 items-center justify-center rounded-xl brand-fill text-white text-sm font-bold">{{ i + 1 }}</span>
              <span class="text-[11px] text-slate-400 font-mono">0{{ i + 1 }}</span>
            </div>
            <div class="mt-4 text-base font-bold text-slate-900">{{ t(`landing.${step.key}`) }}</div>
            <div class="mt-2 text-sm text-slate-600 leading-relaxed">{{ t(`landing.${step.descKey}`) }}</div>
          </div>
        </div>

        <div class="mt-10 grid grid-cols-2 sm:grid-cols-4 gap-3">
          <div v-for="m in campaignMetrics" :key="m.key" class="glass-tile p-4 rounded-2xl text-center">
            <div class="text-xs text-slate-500 uppercase tracking-wide">{{ t(`landing.${m.key}`) }}</div>
            <div class="mt-1 text-2xl font-bold brand-fill-text">{{ m.value }}</div>
          </div>
        </div>

        <p class="mt-8 max-w-3xl mx-auto text-center text-xs text-slate-500 leading-relaxed">
          {{ t('landing.campaignsCompliantNote') }}
        </p>
      </div>
    </section>

    <!-- ============================================================== -->
    <!-- 09. SHARED INBOX                                                -->
    <!-- ============================================================== -->
    <section class="px-4 sm:px-6 py-16 sm:py-20">
      <div class="max-w-6xl mx-auto animate-fade-up">
        <div class="text-center max-w-3xl mx-auto">
          <div class="text-xs uppercase tracking-wider font-bold brand-fill-text inline-block">{{ t('landing.inboxEyebrow') }}</div>
          <h2 class="mt-3 text-3xl sm:text-4xl md:text-5xl font-bold tracking-tight text-slate-900 leading-tight">
            {{ t('landing.inboxTitle') }}
          </h2>
          <p class="mt-4 text-base sm:text-lg text-slate-600 leading-relaxed">
            {{ t('landing.inboxSubtitle') }}
          </p>
        </div>

        <div class="mt-12 grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
          <div v-for="f in inboxFeatures" :key="f.key" class="glass-card p-5 rounded-2xl">
            <span class="inline-flex h-10 w-10 items-center justify-center rounded-xl bg-white/80 text-slate-700 border border-slate-200/60">
              <component :is="f.icon" class="h-5 w-5" />
            </span>
            <div class="mt-4 text-sm font-bold text-slate-900">{{ t(`landing.${f.key}`) }}</div>
          </div>
        </div>
      </div>
    </section>

    <!-- ============================================================== -->
    <!-- 10. BENTO GRID                                                  -->
    <!-- ============================================================== -->
    <section class="px-4 sm:px-6 py-16 sm:py-20">
      <div class="max-w-6xl mx-auto animate-fade-up">
        <div class="text-center max-w-3xl mx-auto">
          <div class="text-xs uppercase tracking-wider font-bold brand-fill-text inline-block">{{ t('landing.bentoEyebrow') }}</div>
          <h2 class="mt-3 text-3xl sm:text-4xl md:text-5xl font-bold tracking-tight text-slate-900 leading-tight">
            {{ t('landing.bentoTitle') }}
          </h2>
        </div>

        <div class="mt-12 grid md:grid-cols-3 gap-4 sm:gap-5">
          <div v-for="tile in bentoTiles" :key="tile.key" class="glass-card relative overflow-hidden p-5 sm:p-6 rounded-3xl bg-gradient-to-br" :class="[tile.accent, tile.span]">
            <span class="inline-flex h-10 w-10 items-center justify-center rounded-xl bg-white/70 text-slate-800">
              <component :is="tile.icon" class="h-5 w-5" />
            </span>
            <div class="mt-4 text-base sm:text-lg font-bold text-slate-900">{{ t(`landing.${tile.key}`) }}</div>
            <div class="mt-1 text-sm text-slate-600 leading-relaxed">{{ t(`landing.${tile.descKey}`) }}</div>
          </div>
        </div>
      </div>
    </section>

    <!-- ============================================================== -->
    <!-- 11. HOW IT WORKS                                                -->
    <!-- ============================================================== -->
    <section class="px-4 sm:px-6 py-16 sm:py-20">
      <div class="max-w-6xl mx-auto animate-fade-up">
        <div class="text-center max-w-3xl mx-auto">
          <div class="text-xs uppercase tracking-wider font-bold brand-fill-text inline-block">{{ t('landing.howEyebrow') }}</div>
          <h2 class="mt-3 text-3xl sm:text-4xl md:text-5xl font-bold tracking-tight text-slate-900 leading-tight">
            {{ t('landing.howTitle') }}
          </h2>
        </div>

        <div class="mt-12 grid md:grid-cols-3 gap-5">
          <div v-for="(step, i) in howSteps" :key="step.key" class="glass-panel relative p-6 rounded-2xl overflow-hidden">
            <div aria-hidden="true" class="absolute -top-12 -end-12 h-32 w-32 rounded-full bg-violet-300/40 blur-3xl"></div>
            <div class="relative">
              <div class="flex items-center justify-between">
                <span class="font-mono text-xs font-bold text-slate-400 tracking-wider">{{ step.num }}</span>
                <span class="inline-flex h-9 w-9 items-center justify-center rounded-xl brand-fill text-white">
                  <component :is="step.icon" class="h-4 w-4" />
                </span>
              </div>
              <div class="mt-6 text-base font-bold text-slate-900">{{ t(`landing.${step.key}`) }}</div>
              <div class="mt-2 text-sm text-slate-600 leading-relaxed">{{ t(`landing.${step.descKey}`) }}</div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ============================================================== -->
    <!-- 12. PRICING — 3 tiers w/ monthly-annual toggle                  -->
    <!-- ============================================================== -->
    <section id="pricing" class="px-4 sm:px-6 py-16 sm:py-24 relative">
      <div aria-hidden="true" class="absolute inset-0 -z-10 overflow-hidden">
        <div class="absolute top-1/2 start-1/2 -translate-x-1/2 -translate-y-1/2 h-96 w-96 rounded-full bg-violet-300/40 blur-3xl"></div>
      </div>

      <div class="max-w-6xl mx-auto animate-fade-up">
        <div class="text-center max-w-3xl mx-auto">
          <div class="text-xs uppercase tracking-wider font-bold brand-fill-text inline-block">{{ t('landing.pricingHeaderEyebrow') }}</div>
          <h2 class="mt-3 text-3xl sm:text-4xl md:text-5xl font-bold tracking-tight text-slate-900 leading-tight">
            {{ t('landing.pricingHeaderTitle') }}
          </h2>
          <p class="mt-4 text-base sm:text-lg text-slate-600 leading-relaxed">
            {{ t('landing.pricingHeaderSubtitle') }}
          </p>
        </div>

        <!-- Monthly / Annual toggle -->
        <div class="mt-10 flex justify-center">
          <div class="inline-flex items-center gap-1 rounded-full glass-tile p-1">
            <button
              class="px-4 sm:px-5 py-2 rounded-full text-sm font-semibold transition-colors"
              :class="!billingAnnual ? 'bg-slate-900 text-white shadow-md' : 'text-slate-600 hover:text-slate-900'"
              @click="billingAnnual = false"
            >
              {{ t('landing.pricingToggleMonthly') }}
            </button>
            <button
              class="px-4 sm:px-5 py-2 rounded-full text-sm font-semibold transition-colors flex items-center gap-2"
              :class="billingAnnual ? 'bg-slate-900 text-white shadow-md' : 'text-slate-600 hover:text-slate-900'"
              @click="billingAnnual = true"
            >
              {{ t('landing.pricingToggleAnnual') }}
              <span
                class="inline-flex items-center px-1.5 py-0.5 rounded-full text-[10px] font-bold"
                :class="billingAnnual ? 'bg-white/20 text-white' : 'brand-fill text-white'"
              >
                {{ t('landing.pricingAnnualBadge') }}
              </span>
            </button>
          </div>
        </div>

        <!-- Tier cards -->
        <div class="mt-10 grid lg:grid-cols-3 gap-5 sm:gap-6">
          <div
            v-for="tier in pricingTiers"
            :key="tier.id"
            class="glass-panel relative overflow-hidden p-7 rounded-3xl flex flex-col"
            :class="tier.featured ? 'ring-2 ring-violet-400/40 lg:scale-[1.02]' : ''"
          >
            <div v-if="tier.featured" aria-hidden="true" class="absolute -top-20 -end-20 h-60 w-60 rounded-full bg-violet-400/40 blur-3xl"></div>
            <div v-if="tier.featured" aria-hidden="true" class="absolute -bottom-20 -start-20 h-60 w-60 rounded-full bg-cyan-300/40 blur-3xl"></div>

            <div class="relative">
              <div class="flex items-center justify-between gap-2">
                <div class="text-base font-bold text-slate-900">{{ t(`landing.${tier.nameKey}`) }}</div>
                <span
                  v-if="tier.tagKey"
                  class="inline-flex items-center px-2.5 py-1 rounded-full text-[11px] font-bold uppercase tracking-wider text-white brand-fill"
                >
                  {{ t(`landing.${tier.tagKey}`) }}
                </span>
              </div>
              <p class="mt-2 text-sm text-slate-600 leading-relaxed">{{ t(`landing.${tier.descKey}`) }}</p>

              <!-- Price -->
              <div class="mt-6 flex items-baseline gap-2">
                <span class="text-4xl sm:text-5xl font-bold tracking-tight brand-fill-text">
                  {{ t(`landing.${billingAnnual ? tier.annualKey : tier.monthlyKey}`) }}
                </span>
              </div>
              <div class="mt-1 text-xs text-slate-500">
                <template v-if="!billingAnnual">{{ t('landing.pricingUnitMonth') }}</template>
                <template v-else>{{ t('landing.pricingUnitAnnualPrefix') }}{{ t('landing.pricingAnnualBadge') }}</template>
              </div>
              <div v-if="billingAnnual" class="mt-2 text-xs font-semibold text-emerald-700">
                ✓ {{ t(`landing.${tier.annualSaveKey}`) }}
              </div>

              <!-- Features -->
              <ul class="mt-6 space-y-2.5 text-sm text-slate-800">
                <li
                  v-for="(line, i) in tierFeatures(tier.featuresPath)"
                  :key="`${tier.id}-f-${i}`"
                  class="flex items-start gap-2.5"
                >
                  <span class="inline-flex h-5 w-5 flex-shrink-0 items-center justify-center rounded-full brand-fill text-white mt-0.5">
                    <Check class="h-3 w-3" />
                  </span>
                  <span :class="i === 0 ? 'font-semibold' : ''">{{ line }}</span>
                </li>
              </ul>

              <p v-if="tier.noIncludesKey" class="mt-4 text-xs text-slate-500 italic">
                {{ t(`landing.${tier.noIncludesKey}`) }}
              </p>
            </div>

            <a :href="whatsappLink" target="_blank" rel="noopener" class="relative mt-7 block">
              <button
                class="w-full inline-flex items-center justify-center gap-2 rounded-2xl px-5 py-3 text-sm font-bold transition-shadow"
                :class="tier.featured ? 'text-white shadow-xl hover:shadow-2xl' : 'text-slate-800 glass-card hover:translate-y-[-1px]'"
                :style="tier.featured ? { background: heroSolid } : {}"
              >
                {{ t(`landing.${tier.ctaKey}`) }}
                <ArrowRight class="h-4 w-4 flip-arrow-end" />
              </button>
            </a>
          </div>
        </div>

        <p class="mt-8 text-[11px] leading-relaxed text-slate-500 text-center max-w-2xl mx-auto">
          {{ t('landing.pricingFinePrint') }}
        </p>
      </div>
    </section>

    <!-- ============================================================== -->
    <!-- 13. SOCIAL PROOF PLACEHOLDER                                   -->
    <!-- ============================================================== -->
    <section class="px-4 sm:px-6 py-12">
      <div class="max-w-5xl mx-auto text-center animate-fade-up">
        <div class="text-xs uppercase tracking-wider font-bold brand-fill-text inline-block">{{ t('landing.socialProofEyebrow') }}</div>
        <h3 class="mt-2 text-xl sm:text-2xl font-bold text-slate-900">
          {{ t('landing.socialProofTitle') }}
        </h3>
        <p class="mt-2 text-sm text-slate-500">{{ t('landing.socialProofNote') }}</p>
        <div class="mt-6 glass-card rounded-3xl px-6 py-8 flex flex-wrap items-center justify-center gap-3">
          <span v-for="placeholder in 4" :key="placeholder" class="inline-flex items-center gap-2 rounded-full bg-white/70 border border-slate-200/60 px-3 py-1.5 text-xs text-slate-400">
            <span class="h-2 w-2 rounded-full bg-slate-300"></span>
            —
          </span>
        </div>
      </div>
    </section>

    <!-- ============================================================== -->
    <!-- 14. FAQ                                                         -->
    <!-- ============================================================== -->
    <section id="faq" class="px-4 sm:px-6 py-16 sm:py-20">
      <div class="max-w-3xl mx-auto animate-fade-up">
        <div class="text-center">
          <div class="text-xs uppercase tracking-wider font-bold brand-fill-text inline-block">{{ t('landing.faqEyebrow') }}</div>
          <h2 class="mt-3 text-3xl sm:text-4xl md:text-5xl font-bold tracking-tight text-slate-900 leading-tight">
            {{ t('landing.faqTitle') }}
          </h2>
        </div>

        <div class="mt-10 space-y-3">
          <div v-for="(faq, i) in faqs" :key="faq.qKey" class="glass-card rounded-2xl overflow-hidden">
            <button
              class="w-full px-5 py-4 sm:py-5 text-start flex items-center justify-between gap-4"
              @click="toggleFaq(i)"
              :aria-expanded="openFaq === i"
            >
              <span class="text-sm sm:text-base font-semibold text-slate-900 leading-relaxed">{{ t(`landing.${faq.qKey}`) }}</span>
              <span class="flex-shrink-0 inline-flex h-7 w-7 items-center justify-center rounded-full brand-fill text-white">
                <component :is="openFaq === i ? Minus : Plus" class="h-4 w-4" />
              </span>
            </button>
            <div v-if="openFaq === i" class="px-5 pb-5 -mt-1 text-sm sm:text-[15px] leading-relaxed text-slate-600">
              {{ t(`landing.${faq.aKey}`) }}
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ============================================================== -->
    <!-- 15. FINAL CTA                                                   -->
    <!-- ============================================================== -->
    <section class="px-4 sm:px-6 py-20 sm:py-28">
      <div class="max-w-5xl mx-auto relative animate-fade-up">
        <div class="glass-panel relative overflow-hidden p-10 sm:p-16 rounded-3xl text-center">
          <div aria-hidden="true" class="absolute inset-0 -z-0">
            <div class="absolute inset-0 bg-gradient-to-br from-violet-500/10 via-transparent to-cyan-400/10"></div>
            <div class="absolute -top-24 -end-24 h-72 w-72 rounded-full bg-violet-400/40 blur-3xl"></div>
            <div class="absolute -bottom-24 -start-24 h-72 w-72 rounded-full bg-cyan-300/40 blur-3xl"></div>
            <div class="absolute start-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 opacity-20">
              <BrandLogo :size="280" float />
            </div>
          </div>
          <div class="relative">
            <h2 class="text-3xl sm:text-4xl md:text-5xl font-bold tracking-tight text-slate-900 leading-tight max-w-3xl mx-auto">
              {{ t('landing.ctaTitle') }}
            </h2>
            <p class="mt-4 text-base sm:text-lg text-slate-600">{{ t('landing.ctaSubtitle') }}</p>
            <div class="mt-8 flex flex-wrap items-center justify-center gap-3">
              <RouterLink to="/register">
                <button class="inline-flex items-center gap-2 rounded-2xl px-6 py-3.5 text-base font-bold text-white shadow-xl hover:shadow-2xl transition-shadow" :style="{ background: heroSolid }">
                  {{ t('landing.ctaPrimary') }}
                  <ArrowRight class="h-4 w-4 flip-arrow-end" />
                </button>
              </RouterLink>
              <a :href="whatsappLink" target="_blank" rel="noopener" class="inline-flex items-center gap-2 rounded-2xl px-6 py-3.5 text-base font-semibold text-slate-800 glass-card hover:translate-y-[-1px] transition-transform">
                <MessageCircle class="h-4 w-4" />
                {{ t('landing.ctaSecondary') }}
              </a>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ============================================================== -->
    <!-- 16. FOOTER                                                      -->
    <!-- ============================================================== -->
    <footer class="px-4 sm:px-6 pt-16 pb-12 border-t border-slate-200/60">
      <div class="max-w-7xl mx-auto grid sm:grid-cols-2 lg:grid-cols-5 gap-8">
        <div class="lg:col-span-2">
          <div class="flex items-center gap-2.5">
            <BrandLogo :size="36" />
            <div class="flex flex-col leading-none">
              <span class="font-semibold text-[15px] tracking-tight text-slate-900">{{ brandName }}</span>
              <span class="text-[10px] mt-1 text-slate-500 tracking-wide">{{ t('landing.footerTagline') }}</span>
            </div>
          </div>
          <p class="mt-4 text-sm text-slate-600 max-w-xs leading-relaxed">{{ t('landing.footerTagline') }}</p>

          <!-- Trust chips in footer -->
          <div class="mt-5 flex flex-wrap gap-2 text-xs text-slate-600">
            <span class="inline-flex items-center gap-1.5 rounded-full bg-white/70 border border-slate-200/60 px-2.5 py-1">
              <BadgeCheck class="h-3.5 w-3.5 text-emerald-600" />
              {{ t('landing.metaPartnerTopBadge') }}
            </span>
            <span class="inline-flex items-center gap-1.5 rounded-full bg-white/70 border border-slate-200/60 px-2.5 py-1">
              <ShieldCheck class="h-3.5 w-3.5 text-emerald-600" />
              {{ t('landing.metaPartnerTopSub') }}
            </span>
          </div>
        </div>

        <div>
          <div class="text-xs font-bold uppercase tracking-wider text-slate-500 mb-3">{{ t('landing.footerProduct') }}</div>
          <ul class="space-y-2 text-sm text-slate-700">
            <li><button class="hover:text-slate-900" @click="scrollToId('features')">{{ t('landing.footerFeatures') }}</button></li>
            <li><button class="hover:text-slate-900" @click="scrollToId('pricing')">{{ t('landing.footerPricing') }}</button></li>
            <li><button class="hover:text-slate-900" @click="scrollToId('faq')">{{ t('landing.footerFaq') }}</button></li>
          </ul>
        </div>

        <div>
          <div class="text-xs font-bold uppercase tracking-wider text-slate-500 mb-3">{{ t('landing.footerContact') }}</div>
          <ul class="space-y-2 text-sm text-slate-700">
            <li><a :href="whatsappLink" target="_blank" rel="noopener" class="hover:text-slate-900">WhatsApp</a></li>
            <li><a href="mailto:hi@omnihub.mavacore.com" class="hover:text-slate-900">hi@omnihub.mavacore.com</a></li>
          </ul>
        </div>

        <div>
          <div class="text-xs font-bold uppercase tracking-wider text-slate-500 mb-3">{{ t('landing.footerLegal') }}</div>
          <ul class="space-y-2 text-sm text-slate-700">
            <li><a href="#" class="hover:text-slate-900">{{ t('landing.footerTerms') }}</a></li>
            <li><a href="#" class="hover:text-slate-900">{{ t('landing.footerPrivacy') }}</a></li>
          </ul>
        </div>
      </div>

      <div class="max-w-7xl mx-auto mt-10 pt-6 border-t border-slate-200/60 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 text-xs text-slate-500">
        <div>{{ t('landing.footerCopyright', { year: currentYear }) }}</div>
        <LandingLanguageToggle class="sm:hidden" />
      </div>
    </footer>

    <!-- ============================================================== -->
    <!-- 17. DEMO MODAL                                                  -->
    <!-- ============================================================== -->
    <Teleport to="body">
      <div v-if="demoOpen" class="fixed inset-0 z-[100] flex items-center justify-center p-4 animate-fade-up" @click.self="closeDemo">
        <div aria-hidden="true" class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm"></div>
        <div class="glass-panel relative max-w-md w-full rounded-3xl overflow-hidden">
          <div aria-hidden="true" class="absolute -top-20 -end-20 h-60 w-60 rounded-full bg-violet-400/40 blur-3xl"></div>
          <div class="relative p-6 sm:p-8">
            <div class="flex items-center justify-between gap-4">
              <div class="text-base font-bold text-slate-900">{{ t('landing.demoTitle') }}</div>
              <button class="inline-flex h-8 w-8 items-center justify-center rounded-full bg-white/70 text-slate-700 hover:bg-white" @click="closeDemo" :aria-label="t('landing.demoClose') as string">
                <X class="h-4 w-4" />
              </button>
            </div>
            <p class="mt-2 text-sm text-slate-600 leading-relaxed">{{ t('landing.demoSubtitle') }}</p>

            <ol class="mt-6 space-y-3">
              <li v-for="(step, i) in demoSteps" :key="step.key" class="flex items-start gap-3">
                <span class="inline-flex h-7 w-7 flex-shrink-0 items-center justify-center rounded-full brand-fill text-white text-xs font-bold">{{ i + 1 }}</span>
                <span class="text-sm text-slate-800 leading-relaxed">{{ t(`landing.${step.key}`) }}</span>
              </li>
            </ol>

            <div class="mt-7 flex items-center justify-end gap-2">
              <button class="inline-flex items-center gap-2 rounded-xl px-4 py-2 text-sm font-semibold text-slate-700 bg-white/70 hover:bg-white" @click="closeDemo">
                {{ t('landing.demoClose') }}
              </button>
              <RouterLink to="/register">
                <button class="inline-flex items-center gap-2 rounded-xl px-4 py-2 text-sm font-semibold text-white shadow-md" :style="{ background: heroSolid }">
                  {{ t('landing.ctaPrimary') }}
                  <ArrowRight class="h-3.5 w-3.5 flip-arrow-end" />
                </button>
              </RouterLink>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- ============================================================== -->
    <!-- 18. LAUNCH OFFER POPUP — floating FAB + modal                  -->
    <!-- ============================================================== -->
    <Teleport to="body">
      <button
        v-if="offerFabVisible && !offerModalOpen"
        type="button"
        @click="openOfferModal"
        :aria-label="t('landing.popupOfferAriaLabel') as string"
        class="fixed z-40 bottom-24 md:bottom-8 end-4 sm:end-6 group animate-fade-up"
      >
        <span aria-hidden="true" class="absolute inset-0 rounded-full brand-fill blur-md opacity-70 group-hover:opacity-90 animate-pulse"></span>
        <span class="relative inline-flex items-center gap-2 rounded-full brand-fill text-white pl-3 pr-4 py-3 text-xs sm:text-sm font-bold shadow-2xl border border-white/30">
          <Sparkles class="h-4 w-4" />
          <span class="hidden sm:inline">{{ t('landing.popupOfferFloatingLabel') }}</span>
          <span class="sm:hidden">🎉</span>
        </span>
      </button>
    </Teleport>

    <Teleport to="body">
      <div v-if="offerModalOpen" class="fixed inset-0 z-[110] flex items-end sm:items-center justify-center p-0 sm:p-4 animate-fade-up" @click.self="closeOfferModal">
        <div aria-hidden="true" class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm"></div>
        <div class="glass-panel relative w-full sm:max-w-md rounded-t-3xl sm:rounded-3xl overflow-hidden">
          <div aria-hidden="true" class="absolute -top-20 -end-20 h-60 w-60 rounded-full bg-violet-400/40 blur-3xl"></div>
          <div aria-hidden="true" class="absolute -bottom-20 -start-20 h-60 w-60 rounded-full bg-cyan-300/40 blur-3xl"></div>
          <div class="relative p-6 sm:p-8 text-start">
            <div class="flex items-center justify-between gap-4">
              <span class="inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-[11px] font-bold uppercase tracking-wider text-white brand-fill">
                🔥 {{ t('landing.popupOfferBadge') }}
              </span>
              <button class="inline-flex h-8 w-8 items-center justify-center rounded-full bg-white/70 text-slate-700 hover:bg-white" @click="closeOfferModal" :aria-label="t('landing.popupOfferClose') as string">
                <X class="h-4 w-4" />
              </button>
            </div>

            <h3 class="mt-4 text-2xl sm:text-3xl font-bold tracking-tight text-slate-900">
              {{ t('landing.popupOfferTitle') }}
            </h3>
            <p class="mt-2 text-sm text-slate-600 leading-relaxed">
              {{ t('landing.popupOfferSubtitle') }}
            </p>

            <ul class="mt-5 space-y-2.5 text-sm text-slate-800">
              <li v-for="(line, i) in offerFeatureList()" :key="i" class="flex items-start gap-2.5">
                <span class="inline-flex h-5 w-5 flex-shrink-0 items-center justify-center rounded-full brand-fill text-white mt-0.5">
                  <Check class="h-3 w-3" />
                </span>
                <span>{{ line }}</span>
              </li>
            </ul>

            <p class="mt-5 text-xs text-slate-600 border-t border-slate-200/60 pt-3">
              {{ t('landing.popupOfferAfter') }}
            </p>

            <div class="mt-6 flex flex-col sm:flex-row gap-2">
              <a :href="whatsappLink" target="_blank" rel="noopener" class="flex-1">
                <button class="w-full inline-flex items-center justify-center gap-2 rounded-2xl px-5 py-3 text-sm font-bold text-white shadow-xl" :style="{ background: heroSolid }">
                  {{ t('landing.popupOfferCta') }}
                  <ArrowRight class="h-4 w-4 flip-arrow-end" />
                </button>
              </a>
              <button class="inline-flex items-center justify-center rounded-2xl px-5 py-3 text-sm font-semibold text-slate-700 glass-card hover:translate-y-[-1px] transition-transform" @click="closeOfferModal">
                {{ t('landing.popupOfferDismiss') }}
              </button>
            </div>

            <p class="mt-4 text-[11px] leading-relaxed text-slate-500">
              {{ t('landing.popupOfferFinePrint') }}
            </p>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- ============================================================== -->
    <!-- 19. MOBILE STICKY CTA — neutral "Get started"                  -->
    <!-- ============================================================== -->
    <div class="mobile-sticky-cta md:hidden">
      <RouterLink to="/app/register" class="block">
        <button class="w-full inline-flex items-center justify-center gap-2 rounded-2xl px-4 py-3 text-sm font-bold text-white shadow-xl" :style="{ background: heroSolid }">
          {{ t('landing.mobileStickyCta') }}
          <ArrowRight class="h-4 w-4 flip-arrow-end" />
        </button>
      </RouterLink>
    </div>
  </div>
</template>

<style scoped>
.landing-shell { font-feature-settings: 'cv02','cv11'; }

/* Force-light for the landing surface. The marketing page is light-only —
   we don't follow the user's OS dark setting here so the glassmorphism
   looks identical in every visitor's preview (and matches the Meta
   partner + offer tone). Overrides any `.dark` ancestor. */
.landing-shell,
.landing-shell * {
  color-scheme: light !important;
}
.landing-shell { background-color: #f8fafc !important; color: #0f172a !important; }
.landing-shell .text-slate-900 { color: #0f172a !important; }
.landing-shell .bg-slate-50 { background-color: #f8fafc !important; }
.landing-shell .border-slate-200\/60 { border-color: rgba(226, 232, 240, 0.6) !important; }

/* Mobile safe-area padding so the sticky CTA never sits behind the home indicator. */
.mobile-sticky-cta { padding-bottom: calc(0.75rem + env(safe-area-inset-bottom)); z-index: 40; }
</style>
