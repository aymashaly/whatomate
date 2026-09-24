<script setup lang="ts">
/**
 * OmniHub Brand Mark — the rounded glass O/H symbol. Circular outer ring,
 * integrated H sitting inside, two side connection dots hinting at
 * WhatsApp Coexistence (phone on one side, OmniHub on the other).
 *
 * Visual defaults: sits on a soft white glass tile (`.brand-mark-glass`)
 * with a brand-violet O/H glyph inside. When the admin uploads a logo
 * the BrandLogo wrapper (sibling component) takes over and this fallback
 * is never rendered.
 *
 * Pass `tinted` when you want a flat colored card with white strokes
 * (e.g. on a saturated hero — small enough that it doesn't dominate).
 *
 * The glyph is solid-color (no SVG gradient) so it stays flat across the
 * landing surface.
 */
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    size?: number
    /** Render on a flat colored card with white strokes (used on saturated backgrounds). */
    tinted?: boolean
    /** Slow floating animation — used in the final CTA watermark. */
    float?: boolean
  }>(),
  { size: 32, tinted: false, float: false }
)

const glyphColor = computed(() => (props.tinted ? '#ffffff' : '#7c3aed'))
</script>

<template>
  <span
    class="inline-flex items-center justify-center overflow-hidden align-middle brand-mark-tile"
    :class="[
      props.tinted ? 'brand-mark-tinted' : 'brand-mark-glass',
      props.float ? 'animate-brand-float' : ''
    ]"
    :style="{
      width: `${props.size}px`,
      height: `${props.size}px`
    }"
  >
    <svg
      :width="props.size * 0.7"
      :height="props.size * 0.7"
      viewBox="0 0 48 48"
      fill="none"
      aria-hidden="true"
    >
      <!-- Outer O ring -->
      <path
        d="M24 6
           a 18 18 0 1 0 0.001 0
           M 24 11
           a 13 13 0 1 1 -0.001 0"
        fill="none"
        :stroke="glyphColor"
        stroke-width="2.5"
        stroke-linecap="round"
      />

      <!-- Integrated H inside the ring -->
      <path
        d="M 16 19 V 29 M 28 19 V 29 M 16 24 H 28"
        :stroke="glyphColor"
        stroke-width="2.5"
        stroke-linecap="round"
        stroke-linejoin="round"
      />

      <!-- Two side connection dots — WhatsApp (left) + OmniHub (right) -->
      <circle cx="2"  cy="24" r="2.5" :fill="glyphColor" />
      <circle cx="46" cy="24" r="2.5" :fill="glyphColor" />
    </svg>
  </span>
</template>


