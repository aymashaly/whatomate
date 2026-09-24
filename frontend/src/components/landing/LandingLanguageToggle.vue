<script setup lang="ts">
/**
 * Landing-page language toggle. Per the brief, the marketing page only
 * needs an Arabic/English switch (AR-default for Saudi launch, EN for
 * international visitors). The full multi-locale picker used in the app
 * shell lives at `src/components/LanguageSwitcher.vue`.
 *
 * Switching locales triggers `setLocale(...)` from the i18n module which
 * flips `<html dir="rtl|ltr">` and `<html lang>`, so the page rerenders
 * mirror-flipped for Arabic without any extra plumbing here.
 */
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { setLocale } from '@/i18n'

const { locale } = useI18n()

const current = computed(() => locale.value)

function toggle() {
  // Bounce between AR (default) and EN. We persist the choice via
  // `setLocale`'s localStorage write so the next visit remembers it.
  setLocale(current.value === 'ar' ? 'en' : 'ar')
}
</script>

<template>
  <button
    type="button"
    @click="toggle"
    class="glass-tile inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-semibold tracking-wide text-slate-700 hover:text-slate-900 transition-colors"
    :aria-label="$t('landing.changeLanguage')"
  >
    <span
      class="inline-block px-1.5 py-0.5 rounded-md text-[11px] transition-colors"
      :class="current === 'ar' ? 'brand-fill text-white' : 'text-slate-500'"
    >AR</span>
    <span class="text-slate-300">/</span>
    <span
      class="inline-block px-1.5 py-0.5 rounded-md text-[11px] transition-colors"
      :class="current === 'en' ? 'brand-fill text-white' : 'text-slate-500'"
    >EN</span>
  </button>
</template>
