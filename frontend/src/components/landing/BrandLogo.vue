<script setup lang="ts">
/**
 * Smart logo component for the landing surface.
 *
 * - When the platform admin has uploaded a custom logo at
 *   /app/settings/branding, this renders that logo image with a white
 *   background (so transparent logos stay legible).
 * - Otherwise it falls back to the inline BrandMark so the page never
 *   shows a broken/missing image.
 *
 * Sized in pixels via the `size` prop. The wrapping span keeps the
 * box exactly the requested size so layout doesn't shift when the
 * branding store loads after first paint.
 */
import { computed } from 'vue'
import { useBrandingStore } from '@/stores/branding'
import BrandMark from './BrandMark.vue'

const props = withDefaults(
  defineProps<{
    /** Edge length in pixels. */
    size?: number
    /** When the BrandMark fallback is used, force the tinted variant. */
    tinted?: boolean
    /** Slow floating animation — used in the final CTA watermark. */
    float?: boolean
  }>(),
  { size: 36, tinted: false, float: false }
)

const brandingStore = useBrandingStore()
const hasLogo = computed(() => !!brandingStore.logoUrl)
const logoUrl = computed(() => brandingStore.logoUrl)
const altText = computed(() => brandingStore.brandName || 'Brand logo')
</script>

<template>
  <span
    class="inline-flex items-center justify-center overflow-hidden align-middle flex-shrink-0"
    :class="float ? 'animate-brand-float' : ''"
    :style="{
      width: `${size}px`,
      height: `${size}px`,
      borderRadius: '0.5rem'
    }"
  >
    <img
      v-if="hasLogo"
      :src="logoUrl"
      :alt="altText"
      class="w-full h-full object-contain"
      loading="lazy"
    />
    <BrandMark v-else :size="size" :tinted="tinted" :float="float" />
  </span>
</template>
